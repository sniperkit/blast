//  Copyright (c) 2017 Minoru Osuka
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

package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mosuka/blast/client"
	"github.com/mosuka/blast/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type HTTPServer struct {
	router      *mux.Router
	listener    net.Listener
	grpcClient  *client.GRPCClient
	dialTimeout int
}

func NewHTTPServer(httpListenAddress string, restPath string, metricsPath string, ctx context.Context, grpcListenAddress string, dialOpts ...grpc.DialOption) (*HTTPServer, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)

	// create client
	grpcClient, err := client.NewGRPCClient(ctx, grpcListenAddress, dialOpts...)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to create gRPC client.")
		return nil, err
	}

	// set handlers
	router.Handle(fmt.Sprintf("%s/", restPath), handler.NewGetIndexInfoHandler(grpcClient)).Methods("GET")
	router.Handle(fmt.Sprintf("%s/{id}", restPath), handler.NewPutDocumentHandler(grpcClient)).Methods("PUT")
	router.Handle(fmt.Sprintf("%s/{id}", restPath), handler.NewGetDocumentHandler(grpcClient)).Methods("GET")
	router.Handle(fmt.Sprintf("%s/{id}", restPath), handler.NewDeleteDocumentHandler(grpcClient)).Methods("DELETE")
	router.Handle(fmt.Sprintf("%s/_bulk", restPath), handler.NewBulkHandler(grpcClient)).Methods("POST")
	router.Handle(fmt.Sprintf("%s/_search", restPath), handler.NewSearchHandler(grpcClient)).Methods("POST")

	router.Handle(metricsPath, promhttp.Handler())

	listener, err := net.Listen("tcp", httpListenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to create HTTP server.")
		return nil, err
	}

	return &HTTPServer{
		router:     router,
		listener:   listener,
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
	// close gRPC client
	err := s.grpcClient.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to close gRPC client.")
		return err
	}

	// stop server
	err = s.listener.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to stop HTTP server.")
		return err
	}

	return nil
}
