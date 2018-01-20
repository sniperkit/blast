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

package server

import (
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/proto"
	"github.com/mosuka/blast/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type IndexServer struct {
	listenAddress string
	server        *grpc.Server
	service       *service.IndexService
}

func NewIndexServer(listenAddress string, indexPath string, indexMapping *mapping.IndexMappingImpl, indexType string, kvstore string, kvconfig map[string]interface{}) (*IndexServer, error) {
	svr := grpc.NewServer()
	svc := service.NewIndexService(indexPath, indexMapping, indexType, kvstore, kvconfig)
	proto.RegisterIndexServer(svr, svc)

	return &IndexServer{
		listenAddress: listenAddress,
		server:        svr,
		service:       svc,
	}, nil
}

func (s *IndexServer) Start() error {
	// open index
	err := s.service.OpenIndex()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to open index.")
		return err
	}

	// create listener
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"listenAddress": s.listenAddress,
			"error":         err.Error(),
		}).Error("failed to start gRPC server.")
		return err
	}

	// start server
	go func() {
		s.server.Serve(listener)
		return
	}()

	return nil
}

func (s *IndexServer) Stop() error {
	// stop server
	s.server.GracefulStop()

	// close index
	err := s.service.CloseIndex()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to close index.")
		return err
	}

	return nil
}
