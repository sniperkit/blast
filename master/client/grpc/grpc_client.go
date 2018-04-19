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
	"context"
	"encoding/json"
	"github.com/blevesearch/bleve/mapping"
	"github.com/golang/protobuf/jsonpb"
	"github.com/mosuka/blast/index"
	"github.com/mosuka/blast/protobuf"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	server        string
	dialOptions   []grpc.DialOption
	context       context.Context
	cancel        context.CancelFunc
	conn          *grpc.ClientConn
	clusterClient protobuf.ClusterClient
}

func NewGRPCClient(ctx context.Context, server string, dialOpts ...grpc.DialOption) (*GRPCClient, error) {
	ct, cancel := context.WithCancel(ctx)

	conn, err := grpc.DialContext(ct, server, dialOpts...)
	if err != nil {
		cancel()
		return nil, err
	}

	clusterClient := protobuf.NewClusterClient(conn)

	grpcClient := &GRPCClient{
		server:        server,
		dialOptions:   dialOpts,
		context:       ct,
		cancel:        cancel,
		conn:          conn,
		clusterClient: clusterClient,
	}

	return grpcClient, nil
}

func (c *GRPCClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.context.Err()
}

func (c *GRPCClient) PutNode(ctx context.Context, cluster string, node string, callOpts ...grpc.CallOption) error {
	protoReq := &protobuf.PutNodeRequest{
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
	protoReq := &protobuf.GetNodeRequest{
		Cluster: cluster,
		Node:    node,
	}

	protoResp, err := c.clusterClient.GetNode(ctx, protoReq, callOpts...)
	if err != nil {
		return nil, err
	}

	marshaler := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     false,
	}
	jsonStr, err := marshaler.MarshalToString(protoResp.NodeInfo)
	if err != nil {
		return nil, err
	}

	var value map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (c *GRPCClient) DeleteNode(ctx context.Context, cluster string, node string, callOpts ...grpc.CallOption) error {
	protoReq := &protobuf.DeleteNodeRequest{
		Cluster: cluster,
		Node:    node,
	}

	_, err := c.clusterClient.DeleteNode(ctx, protoReq, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *GRPCClient) PutIndexMapping(ctx context.Context, cluster string, indexMapping *mapping.IndexMappingImpl, callOpts ...grpc.CallOption) error {
	indexMappingAny, err := protobuf.MarshalAny(indexMapping)

	protoReq := &protobuf.PutIndexMappingRequest{
		Cluster:      cluster,
		IndexMapping: &indexMappingAny,
	}

	_, err = c.clusterClient.PutIndexMapping(ctx, protoReq, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *GRPCClient) GetIndexMapping(ctx context.Context, cluster string, callOpts ...grpc.CallOption) (*mapping.IndexMappingImpl, error) {
	protoReq := &protobuf.GetIndexMappingRequest{
		Cluster: cluster,
	}

	protoResp, err := c.clusterClient.GetIndexMapping(ctx, protoReq, callOpts...)
	if err != nil {
		return nil, err
	}

	indexMapping, err := protobuf.UnmarshalAny(protoResp.IndexMapping)
	if err != nil {
		return nil, err
	}

	return indexMapping.(*mapping.IndexMappingImpl), nil
}

func (c *GRPCClient) DeleteIndexMapping(ctx context.Context, cluster string, callOpts ...grpc.CallOption) error {
	protoReq := &protobuf.DeleteIndexMappingRequest{
		Cluster: cluster,
	}

	_, err := c.clusterClient.DeleteIndexMapping(ctx, protoReq, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *GRPCClient) PutIndexMeta(ctx context.Context, cluster string, indexMeta *index.IndexMeta, callOpts ...grpc.CallOption) error {
	indexMetaAny, err := protobuf.MarshalAny(indexMeta)

	protoReq := &protobuf.PutIndexMetaRequest{
		Cluster:   cluster,
		IndexMeta: &indexMetaAny,
	}

	_, err = c.clusterClient.PutIndexMeta(ctx, protoReq, callOpts...)
	if err != nil {
		return err
	}

	return nil
}
