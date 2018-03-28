//  Copyright (c) 2018 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	blastgrpc "github.com/mosuka/blast/node/client/grpc"
	"github.com/mosuka/blast/node/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type HTTPServer struct {
	listener   net.Listener
	router     *mux.Router
	grpcClient *blastgrpc.GRPCClient
}

func NewHTTPServer(httpListenAddress string, restPath string, metricsPath string, ctx context.Context, grpcListenAddress string, dialOpts ...grpc.DialOption) (*HTTPServer, error) {
	// create client
	grpcClient, err := blastgrpc.NewGRPCClient(ctx, grpcListenAddress, dialOpts...)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to create gRPC client.")
		return nil, err
	}

	router := mux.NewRouter()
	router.StrictSlash(true)

	// set handlers
	router.Handle(fmt.Sprintf("%s/{id}", restPath), handler.NewPutDocumentHandler(grpcClient)).Methods("PUT")
	router.Handle(fmt.Sprintf("%s/{id}", restPath), handler.NewGetDocumentHandler(grpcClient)).Methods("GET")
	router.Handle(fmt.Sprintf("%s/{id}", restPath), handler.NewDeleteDocumentHandler(grpcClient)).Methods("DELETE")
	router.Handle(fmt.Sprintf("%s/_bulk", restPath), handler.NewBulkHandler(grpcClient)).Methods("POST")
	router.Handle(fmt.Sprintf("%s/_search", restPath), handler.NewSearchHandler(grpcClient)).Methods("POST")
	router.Handle(fmt.Sprintf("%s/_indexpath", restPath), handler.NewGetIndexPathHandler(grpcClient)).Methods("GET")
	router.Handle(fmt.Sprintf("%s/_indexmapping", restPath), handler.NewGetIndexMappingHandler(grpcClient)).Methods("GET")
	router.Handle(fmt.Sprintf("%s/_indexmeta", restPath), handler.NewGetIndexMetaHandler(grpcClient)).Methods("GET")

	router.Handle(metricsPath, promhttp.Handler())

	listener, err := net.Listen("tcp", httpListenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to create HTTP server.")
		return nil, err
	}

	return &HTTPServer{
		listener:   listener,
		router:     router,
		grpcClient: grpcClient,
	}, nil
}

func (s *HTTPServer) Start() error {
	// start server
	go func() {
		http.Serve(s.listener, s.router)
		return
	}()

	return nil
}

func (s *HTTPServer) Stop() error {
	// stop server
	err := s.listener.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to stop HTTP server.")
		return err
	}

	// close gRPC client
	err = s.grpcClient.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to close gRPC client.")
		return err
	}

	return nil
}
