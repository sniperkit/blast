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
	"github.com/mosuka/blast/proto"
	"golang.org/x/net/context"
)

type ClusterService struct {
	ClusterName string
}

func NewClusterService(clusterName string) *ClusterService {
	return &ClusterService{
		ClusterName: clusterName,
	}
}

func (c *ClusterService) PutNode(ctx context.Context, req *proto.PutNodeRequest) (*proto.PutNodeResponse, error) {
	return &proto.PutNodeResponse{}, nil
}

func (c *ClusterService) GetNode(ctx context.Context, req *proto.GetNodeRequest) (*proto.GetNodeResponse, error) {
	return &proto.GetNodeResponse{}, nil
}

func (c *ClusterService) DeleteNode(ctx context.Context, req *proto.DeleteNodeRequest) (*proto.DeleteNodeResponse, error) {
	return &proto.DeleteNodeResponse{}, nil
}
