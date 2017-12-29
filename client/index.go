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

package client

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Index interface {
	GetIndexInfo(ctx context.Context, indexPath bool, indexMapping bool, indexType bool, kvstore bool, kvconfig bool, opts ...grpc.CallOption) (string, *mapping.IndexMappingImpl, string, string, map[string]interface{}, error)
	PutDocument(ctx context.Context, id string, fields map[string]interface{}, opts ...grpc.CallOption) (string, map[string]interface{}, error)
	GetDocument(ctx context.Context, id string, opts ...grpc.CallOption) (string, map[string]interface{}, error)
	DeleteDocument(ctx context.Context, id string, opts ...grpc.CallOption) (string, error)
	Bulk(ctx context.Context, requests []map[string]interface{}, batchSize int32, opts ...grpc.CallOption) (int32, int32, int32, int32, error)
	Search(ctx context.Context, searchRequest *bleve.SearchRequest, opts ...grpc.CallOption) (*bleve.SearchResult, error)
}

type index struct {
	client proto.IndexClient
}

func NewIndex(c *BlastClient) Index {
	ic := proto.NewIndexClient(c.conn)

	return &index{
		client: ic,
	}
}

func (i *index) GetIndexInfo(ctx context.Context, indexPath bool, indexMapping bool, indexType bool, kvstore bool, kvconfig bool, opts ...grpc.CallOption) (string, *mapping.IndexMappingImpl, string, string, map[string]interface{}, error) {
	protoReq := &proto.GetIndexInfoRequest{
		IndexMapping: indexMapping,
		IndexType:    indexType,
		Kvstore:      kvstore,
		Kvconfig:     kvconfig,
	}

	protoResp, err := i.client.GetIndexInfo(ctx, protoReq, opts...)
	if err != nil {
		return "", nil, "", "", nil, err
	}

	im, err := proto.UnmarshalAny(protoResp.IndexMapping)
	if err != nil {
		return "", nil, "", "", nil, err
	}

	kvc, err := proto.UnmarshalAny(protoResp.Kvconfig)
	if err != nil {
		return "", nil, "", "", nil, err
	}

	return protoResp.IndexPath, im.(*mapping.IndexMappingImpl), protoResp.IndexType, protoResp.Kvstore, *kvc.(*map[string]interface{}), nil
}

func (i *index) PutDocument(ctx context.Context, id string, fields map[string]interface{}, opts ...grpc.CallOption) (string, map[string]interface{}, error) {
	fieldAny, err := proto.MarshalAny(fields)
	if err != nil {
		return "", nil, err
	}

	protoReq := &proto.PutDocumentRequest{
		Id:     id,
		Fields: &fieldAny,
	}

	protoResp, err := i.client.PutDocument(ctx, protoReq, opts...)
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

func (i *index) GetDocument(ctx context.Context, id string, opts ...grpc.CallOption) (string, map[string]interface{}, error) {
	protoReq := &proto.GetDocumentRequest{
		Id: id,
	}

	protoResp, err := i.client.GetDocument(ctx, protoReq, opts...)
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

func (i *index) DeleteDocument(ctx context.Context, id string, opts ...grpc.CallOption) (string, error) {
	protoReq := &proto.DeleteDocumentRequest{
		Id: id,
	}

	protoResp, err := i.client.DeleteDocument(ctx, protoReq, opts...)
	if err != nil {
		return "", err
	}

	return protoResp.Id, nil
}

func (i *index) Bulk(ctx context.Context, requests []map[string]interface{}, batchSize int32, opts ...grpc.CallOption) (int32, int32, int32, int32, error) {
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

	protoResp, err := i.client.Bulk(ctx, protoReq, opts...)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return protoResp.PutCount, protoResp.PutErrorCount, protoResp.DeleteCount, protoResp.MethodErrorCount, nil
}

func (i *index) Search(ctx context.Context, searchRequest *bleve.SearchRequest, opts ...grpc.CallOption) (*bleve.SearchResult, error) {
	searchResultAny, err := proto.MarshalAny(searchRequest)
	if err != nil {
		return nil, err
	}

	protoReq := &proto.SearchRequest{
		SearchRequest: &searchResultAny,
	}

	protoResp, err := i.client.Search(ctx, protoReq, opts...)
	if err != nil {
		return nil, err
	}

	searchResult, err := proto.UnmarshalAny(protoResp.SearchResult)

	return searchResult.(*bleve.SearchResult), err
}
