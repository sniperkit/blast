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
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/mosuka/blast/proto"
	"github.com/mosuka/blast/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ClusterServer struct {
	listenAddress string
	server        *grpc.Server
	service       *service.ClusterService
	endpoints     []string
	dialTimeout   int
	etcdClient    *clientv3.Client
}

func NewClusterServer(listenAddress string, endpoints []string, dialTimeout int, clusterName string) (*ClusterServer, error) {
	svr := grpc.NewServer()
	svc := service.NewClusterService(clusterName)
	proto.RegisterClusterServer(svr, svc)

	return &ClusterServer{
		listenAddress: listenAddress,
		server:        svr,
		service:       svc,
		endpoints:     endpoints,
		dialTimeout:   dialTimeout,
	}, nil
}

func (s *ClusterServer) Start() error {
	// create etcd client
	cfg := clientv3.Config{
		Endpoints:   s.endpoints,
		DialTimeout: time.Duration(s.dialTimeout) * time.Millisecond,
		Context:     context.Background(),
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to create etcd client.")
		return err
	}
	s.etcdClient = etcdClient

	// create listener
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"listenAddress": s.listenAddress,
			"error":         err.Error(),
		}).Error("failed to start Cluster server.")
		return err
	}

	// start server
	go func() {
		s.server.Serve(listener)
		return
	}()

	return nil
}

func (s *ClusterServer) Stop() error {
	// stop server
	s.server.GracefulStop()

	// close etcd client
	err := s.etcdClient.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to close etcd client.")
		return err
	}

	return nil
}
