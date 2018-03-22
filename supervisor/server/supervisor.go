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
	"github.com/mosuka/blast/pb"
	"github.com/mosuka/blast/supervisor/config"
	"github.com/mosuka/blast/supervisor/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Supervisor struct {
	listenAddress string
	server        *grpc.Server
	service       *service.SuperviseService
}

func NewSupervisor(listenAddress string, config *config.SupervisorConfig) (*Supervisor, error) {
	svc, err := service.NewSuperviseService(config)
	if err != nil {
		return nil, err
	}

	svr := grpc.NewServer()

	pb.RegisterClusterServer(svr, svc)

	return &Supervisor{
		listenAddress: listenAddress,
		server:        svr,
		service:       svc,
	}, nil
}

func (s *Supervisor) Start() error {
	// create listener
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"listenAddress": s.listenAddress,
			"error":         err.Error(),
		}).Error("failed to create listener.")
		return err
	}

	// start server
	go func() {
		s.server.Serve(listener)
		return
	}()

	return nil
}

func (s *Supervisor) Stop() error {
	// stop server
	s.server.GracefulStop()

	return nil
}
