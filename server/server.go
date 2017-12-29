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
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/cluster"
	"github.com/mosuka/blast/proto"
	"github.com/mosuka/blast/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

type BlastServer struct {
	cluster.BlastCluster

	hostname   string
	port       int
	server     *grpc.Server
	service    *service.BlastService
	collection string
	node       string
}

func NewBlastServerInClusterMode(port int, indexPath string, etcdEndpoints []string, etcdDialTimeout int, etcdRequestTimeout int, collection string) (*BlastServer, error) {

	blastCluster, err := cluster.NewBlastCluster(etcdEndpoints, etcdDialTimeout)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(etcdRequestTimeout)*time.Millisecond)
	defer cancel()

	// fetch index mapping
	indexMapping, err := blastCluster.GetIndexMapping(ctx, collection)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// fetch index type
	indexType, err := blastCluster.GetIndexType(ctx, collection)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// fetch kvstore
	kvstore, err := blastCluster.GetKvstore(ctx, collection)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// fetch kvconfig
	kvconfig, err := blastCluster.GetKvconfig(ctx, collection)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	blastServer, err := NewBlastServer(port, indexPath, indexMapping, indexType, kvstore, kvconfig)
	if err != nil {
		return nil, err
	}

	blastServer.collection = collection
	blastServer.BlastCluster = blastCluster

	return blastServer, nil
}

func NewBlastServer(port int, indexPath string, indexMapping *mapping.IndexMappingImpl, indexType string, kvstore string, kvconfig map[string]interface{}) (*BlastServer, error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	node := fmt.Sprintf("%s:%d", hostname, port)

	svr := grpc.NewServer()
	svc := service.NewBlastService(indexPath, indexMapping, indexType, kvstore, kvconfig)
	proto.RegisterIndexServer(svr, svc)

	return &BlastServer{
		hostname: hostname,
		port:     port,
		server:   svr,
		service:  svc,
		node:     node,
	}, nil
}

func (s *BlastServer) Start() error {
	// create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.WithFields(log.Fields{
			"port": s.port,
		}).Error(err.Error())
		return err
	}

	if s.BlastCluster != nil {
		s.watchCollection()
	}

	// start server
	go func() {
		s.service.OpenIndex()
		s.server.Serve(listener)
		return
	}()

	if s.BlastCluster != nil {
		err := s.joinCluster()
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

func (s *BlastServer) Stop() error {
	if s.BlastCluster != nil {
		err := s.leaveCluster()
		if err != nil {
			log.Error(err.Error())
		}

		err = s.BlastCluster.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}

	err := s.service.CloseIndex()
	if err != nil {
		log.Error(err.Error())
	}

	s.server.GracefulStop()

	return err
}

func (s *BlastServer) joinCluster() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15000)*time.Millisecond)
	defer cancel()

	err := s.BlastCluster.PutNode(ctx, s.collection, s.node)

	return err
}

func (s *BlastServer) leaveCluster() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15000)*time.Millisecond)
	defer cancel()

	err := s.BlastCluster.DeleteNode(ctx, s.collection, s.node)

	return err
}

func (s *BlastServer) watchCollection() {
	go func() {
		var ctx context.Context
		var cancel context.CancelFunc

		for {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(15000)*time.Millisecond)
			rch := s.BlastCluster.Watch(ctx, s.collection)
			for wresp := range rch {
				for _, ev := range wresp.Events {
					log.WithFields(log.Fields{
						"type":  ev.Type,
						"key":   fmt.Sprintf("%s", ev.Kv.Key),
						"value": fmt.Sprintf("%s", ev.Kv.Value),
					}).Info("the cluster information has been changed")
				}
			}
		}
		defer cancel()

		return
	}()

	return
}
