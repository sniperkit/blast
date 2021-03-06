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

package grpc

import (
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/index"
	"github.com/mosuka/blast/node/service"
	"github.com/mosuka/blast/protobuf"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	listenAddress string
	server        *grpc.Server
	service       *service.IndexService
}

func NewGRPCServer(listenAddress string, indexPath string, indexMapping *mapping.IndexMappingImpl, indexMeta *index.IndexMeta) (*GRPCServer, error) {
	svc, err := service.NewIndexService(indexPath, indexMapping, indexMeta)
	if err != nil {
		log.WithFields(log.Fields{
			"indexPath":    indexPath,
			"indexMapping": indexMapping,
			"indexMeta":    indexMeta,
			"error":        err.Error(),
		}).Error("failed to create index service.")
		return nil, err
	}

	svr := grpc.NewServer()
	protobuf.RegisterIndexServer(svr, svc)

	return &GRPCServer{
		listenAddress: listenAddress,
		server:        svr,
		service:       svc,
	}, nil
}

func (s *GRPCServer) Start() error {
	// open index
	err := s.service.OpenIndex()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to open index.")
		return err
	}
	log.Info(fmt.Sprintf("index was opened."))

	// create listener
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"listenAddress": s.listenAddress,
			"error":         err.Error(),
		}).Error("failed to create listener.")
		return err
	}
	log.Info(fmt.Sprintf("listener was created."))

	// start server
	go func() {
		s.server.Serve(listener)
		return
	}()

	return nil
}

func (s *GRPCServer) Stop() error {
	// stop server
	s.server.GracefulStop()
	log.Info(fmt.Sprintf("gRPC server was gracefully stopped."))

	// close index
	err := s.service.CloseIndex()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to close index.")
		return err
	}
	log.Info(fmt.Sprintf("index was closed."))

	return nil
}
