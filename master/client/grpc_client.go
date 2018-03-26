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
	"github.com/mosuka/blast/pb"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	server        string
	dialOptions   []grpc.DialOption
	context       context.Context
	cancel        context.CancelFunc
	conn          *grpc.ClientConn
	clusterClient pb.ClusterClient
}

func NewGRPCClient(ctx context.Context, server string, dialOpts ...grpc.DialOption) (*GRPCClient, error) {
	ct, cancel := context.WithCancel(ctx)

	conn, err := grpc.DialContext(ct, server, dialOpts...)
	if err != nil {
		cancel()
		return nil, err
	}

	ic := pb.NewClusterClient(conn)

	c := &GRPCClient{
		server:        server,
		dialOptions:   dialOpts,
		context:       ct,
		cancel:        cancel,
		conn:          conn,
		clusterClient: ic,
	}

	return c, nil
}

func (c *GRPCClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.context.Err()
}

func (c *GRPCClient) PutNode(ctx context.Context, cluster string, node string, callOpts ...grpc.CallOption) error {
	protoReq := &pb.PutNodeRequest{
		Cluster: cluster,
		Node:    node,
	}

	_, err := c.clusterClient.PutNode(ctx, protoReq, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *GRPCClient) GetNode(ctx context.Context, cluster string, node string, callOpts ...grpc.CallOption) (*map[string]interface{}, error) {
	protoReq := &pb.GetNodeRequest{
		Cluster: cluster,
		Node:    node,
	}

	protoResp, err := c.clusterClient.GetNode(ctx, protoReq, callOpts...)
	if err != nil {
		return nil, err
	}

	valueTmp, err := pb.UnmarshalAny(protoResp.Value)
	if err != nil {
		return nil, err
	}

	var value map[string]interface{}
	if valueTmp != nil {
		value = *valueTmp.(*map[string]interface{})
	}

	return &value, nil
}

func (c *GRPCClient) DeleteNode(ctx context.Context, cluster string, node string, callOpts ...grpc.CallOption) error {
	protoReq := &pb.DeleteNodeRequest{
		Cluster: cluster,
		Node:    node,
	}

	_, err := c.clusterClient.DeleteNode(ctx, protoReq, callOpts...)
	if err != nil {
		return err
	}

	return nil
}
