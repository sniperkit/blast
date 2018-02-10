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

package service

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/mosuka/blast/proto"
	"golang.org/x/net/context"
	"time"
)

type ClusterService struct {
	etcdClient *clientv3.Client
	clusterKey string
}

func NewClusterService(etcdEndpoints []string, etcdDialTimeout int, clusterName string) *ClusterService {
	cfg := clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: time.Duration(etcdDialTimeout) * time.Millisecond,
		Context:     context.Background(),
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil
	}

	return &ClusterService{
		etcdClient: client,
		clusterKey: "/blast/clusters/" + clusterName,
	}
}

func (c *ClusterService) CloseClient() error {
	err := c.etcdClient.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *ClusterService) Join(ctx context.Context, req *proto.JoinRequest) (*proto.JoinResponse, error) {
	keyNode := c.clusterKey + "/nodes/" + req.Node

	_, err := c.etcdClient.Put(ctx, keyNode, fmt.Sprintf("%020d", time.Now().UnixNano()))
	if err != nil {
		return nil, err
	}

	return &proto.JoinResponse{}, nil
}

func (c *ClusterService) Leave(ctx context.Context, req *proto.LeaveRequest) (*proto.LeaveResponse, error) {
	keyNode := c.clusterKey + "/nodes/" + req.Node

	_, err := c.etcdClient.Delete(ctx, keyNode)
	if err != nil {
		return nil, err
	}

	return &proto.LeaveResponse{}, nil
}
