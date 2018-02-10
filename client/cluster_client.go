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

package client

import (
	"context"
	"github.com/mosuka/blast/proto"
	"google.golang.org/grpc"
)

type ClusterClient struct {
	server        string
	dialOptions   []grpc.DialOption
	context       context.Context
	cancel        context.CancelFunc
	conn          *grpc.ClientConn
	clusterClient proto.ClusterClient
}

func NewClusterClient(ctx context.Context, server string, dialOpts ...grpc.DialOption) (*ClusterClient, error) {
	ct, cancel := context.WithCancel(ctx)

	conn, err := grpc.DialContext(ct, server, dialOpts...)
	if err != nil {
		cancel()
		return nil, err
	}

	ic := proto.NewClusterClient(conn)

	c := &ClusterClient{
		server:        server,
		dialOptions:   dialOpts,
		context:       ct,
		cancel:        cancel,
		conn:          conn,
		clusterClient: ic,
	}

	return c, nil
}

func (c *ClusterClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.context.Err()
}

func (c *ClusterClient) Join(ctx context.Context, node string, callOpts ...grpc.CallOption) (string, error) {
	protoReq := &proto.JoinRequest{
		Node: node,
	}

	protoResp, err := c.clusterClient.Join(ctx, protoReq, callOpts...)
	if err != nil {
		return "", err
	}

	return protoResp.Node, nil
}

func (c *ClusterClient) Leave(ctx context.Context, node string, callOpts ...grpc.CallOption) (string, error) {
	protoReq := &proto.LeaveRequest{
		Node: node,
	}

	protoResp, err := c.clusterClient.Leave(ctx, protoReq, callOpts...)
	if err != nil {
		return "", err
	}

	return protoResp.Node, nil
}
