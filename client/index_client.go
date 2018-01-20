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
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mosuka/blast/proto"
	"google.golang.org/grpc"
)

type IndexClient struct {
	server      string
	dialOptions []grpc.DialOption
	context     context.Context
	cancel      context.CancelFunc
	conn        *grpc.ClientConn
	indexClient proto.IndexClient
}

func NewIndexClient(ctx context.Context, server string, dialOpts ...grpc.DialOption) (*IndexClient, error) {
	ct, cancel := context.WithCancel(ctx)

	conn, err := grpc.DialContext(ct, server, dialOpts...)
	if err != nil {
		cancel()
		return nil, err
	}

	ic := proto.NewIndexClient(conn)

	c := &IndexClient{
		server:      server,
		dialOptions: dialOpts,
		context:     ct,
		cancel:      cancel,
		conn:        conn,
		indexClient: ic,
	}

	return c, nil
}

func (c *IndexClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.context.Err()
}

func (c *IndexClient) GetIndexPath(ctx context.Context, callOpts ...grpc.CallOption) (string, error) {
	protoReq := &empty.Empty{}

	protoResp, err := c.indexClient.GetIndexPath(ctx, protoReq, callOpts...)
	if err != nil {
		return "", err
	}

	return protoResp.IndexPath, nil
}

func (c *IndexClient) GetIndexMapping(ctx context.Context, callOpts ...grpc.CallOption) (*mapping.IndexMappingImpl, error) {
	protoReq := &empty.Empty{}

	protoResp, err := c.indexClient.GetIndexMapping(ctx, protoReq, callOpts...)
	if err != nil {
		return nil, err
	}

	im, err := proto.UnmarshalAny(protoResp.IndexMapping)
	if err != nil {
		return nil, err
	}

	return im.(*mapping.IndexMappingImpl), nil
}

func (c *IndexClient) GetIndexType(ctx context.Context, callOpts ...grpc.CallOption) (string, error) {
	protoReq := &empty.Empty{}

	protoResp, err := c.indexClient.GetIndexType(ctx, protoReq, callOpts...)
	if err != nil {
		return "", err
	}

	return protoResp.IndexType, nil
}

func (c *IndexClient) GetKvstore(ctx context.Context, callOpts ...grpc.CallOption) (string, error) {
	protoReq := &empty.Empty{}

	protoResp, err := c.indexClient.GetKvstore(ctx, protoReq, callOpts...)
	if err != nil {
		return "", err
	}

	return protoResp.Kvstore, nil
}

func (c *IndexClient) GetKvconfig(ctx context.Context, callOpts ...grpc.CallOption) (map[string]interface{}, error) {
	protoReq := &empty.Empty{}

	protoResp, err := c.indexClient.GetKvconfig(ctx, protoReq, callOpts...)

	kvc, err := proto.UnmarshalAny(protoResp.Kvconfig)
	if err != nil {
		return nil, err
	}

	return *kvc.(*map[string]interface{}), nil
}

func (c *IndexClient) PutDocument(ctx context.Context, id string, fields map[string]interface{}, callOpts ...grpc.CallOption) (string, map[string]interface{}, error) {
	fieldAny, err := proto.MarshalAny(fields)
	if err != nil {
		return "", nil, err
	}

	protoReq := &proto.PutDocumentRequest{
		Id:     id,
		Fields: &fieldAny,
	}

	protoResp, err := c.indexClient.PutDocument(ctx, protoReq, callOpts...)
	if err != nil {
		return "", nil, err
	}

	fieldsTmp, err := proto.UnmarshalAny(protoResp.Fields)
	if err != nil {
		return "", nil, err
	}

	var fieldsPut map[string]interface{}
	if fieldsTmp != nil {
		fieldsPut = *fieldsTmp.(*map[string]interface{})
	}

	return protoResp.Id, fieldsPut, nil
}

func (c *IndexClient) GetDocument(ctx context.Context, id string, callOpts ...grpc.CallOption) (string, map[string]interface{}, error) {
	protoReq := &proto.GetDocumentRequest{
		Id: id,
	}

	protoResp, err := c.indexClient.GetDocument(ctx, protoReq, callOpts...)
	if err != nil {
		return "", nil, err
	}

	fieldsTmp, err := proto.UnmarshalAny(protoResp.Fields)
	if err != nil {
		return "", nil, err
	}

	var fields map[string]interface{}
	if fieldsTmp != nil {
		fields = *fieldsTmp.(*map[string]interface{})
	}

	return protoResp.Id, fields, nil
}

func (c *IndexClient) DeleteDocument(ctx context.Context, id string, callOpts ...grpc.CallOption) (string, error) {
	protoReq := &proto.DeleteDocumentRequest{
		Id: id,
	}

	protoResp, err := c.indexClient.DeleteDocument(ctx, protoReq, callOpts...)
	if err != nil {
		return "", err
	}

	return protoResp.Id, nil
}

func (c *IndexClient) Bulk(ctx context.Context, requests []map[string]interface{}, batchSize int32, callOpts ...grpc.CallOption) (int32, int32, int32, int32, error) {
	updateRequests := make([]*proto.BulkRequest_UpdateRequest, 0)
	for _, updateRequest := range requests {
		r := &proto.BulkRequest_UpdateRequest{}

		// check method
		if _, ok := updateRequest["method"]; ok {
			r.Method = updateRequest["method"].(string)
		}

		// check document
		var document map[string]interface{}
		if _, ok := updateRequest["document"]; ok {
			d := &proto.BulkRequest_UpdateRequest_Document{}

			document = updateRequest["document"].(map[string]interface{})

			// check document.id
			if _, ok := document["id"]; ok {
				d.Id = document["id"].(string)
			}

			// check document.fields
			if _, ok := document["fields"]; ok {
				fields, err := proto.MarshalAny(document["fields"].(map[string]interface{}))
				if err != nil {
					continue
				}
				d.Fields = &fields
			}

			r.Document = d
		}

		updateRequests = append(updateRequests, r)
	}

	protoReq := &proto.BulkRequest{
		BatchSize:      batchSize,
		UpdateRequests: updateRequests,
	}

	protoResp, err := c.indexClient.Bulk(ctx, protoReq, callOpts...)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return protoResp.PutCount, protoResp.PutErrorCount, protoResp.DeleteCount, protoResp.MethodErrorCount, nil
}

func (c *IndexClient) Search(ctx context.Context, searchRequest *bleve.SearchRequest, opts ...grpc.CallOption) (*bleve.SearchResult, error) {
	searchResultAny, err := proto.MarshalAny(searchRequest)
	if err != nil {
		return nil, err
	}

	protoReq := &proto.SearchRequest{
		SearchRequest: &searchResultAny,
	}

	protoResp, err := c.indexClient.Search(ctx, protoReq, opts...)
	if err != nil {
		return nil, err
	}

	searchResult, err := proto.UnmarshalAny(protoResp.SearchResult)

	return searchResult.(*bleve.SearchResult), err
}
