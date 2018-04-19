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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mosuka/blast/cluster"
	_ "github.com/mosuka/blast/master/builtin"
	"github.com/mosuka/blast/master/registry"
	"github.com/mosuka/blast/master/store"
	"github.com/mosuka/blast/protobuf"
	"golang.org/x/net/context"
	"time"
)

type ClusterService struct {
	storeType string
	config    map[string]interface{}
	store     store.Store
}

func NewClusterService(clusterMeta *cluster.ClusterMeta) (*ClusterService, error) {
	storeConstructor := registry.StoreConstructorByName(clusterMeta.Storage)

	store, err := storeConstructor(clusterMeta.Config)
	if err != nil {
		return nil, err
	}

	return &ClusterService{
		storeType: clusterMeta.Storage,
		config:    clusterMeta.Config,
		store:     store,
	}, nil
}

func (s *ClusterService) PutNode(ctx context.Context, req *protobuf.PutNodeRequest) (*empty.Empty, error) {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", s.config["base_path"].(string), req.Cluster, req.Node)

	if req.NodeInfo == nil {
		req.NodeInfo = &protobuf.NodeInfo{}
		req.NodeInfo.Cluster = req.Cluster
		req.NodeInfo.Node = req.Node
		req.NodeInfo.Status = ""
		req.NodeInfo.Role = ""
		req.NodeInfo.Timestamp = time.Now().Format(time.RFC3339Nano)
	}

	// NodeInfo -> JSON string ([]byte)
	marshaler := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     false,
	}
	value, err := marshaler.MarshalToString(req.NodeInfo)

	// put JSON string ([]byte)
	err = s.store.Put(key, []byte(value))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ClusterService) GetNode(ctx context.Context, req *protobuf.GetNodeRequest) (*protobuf.GetNodeResponse, error) {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", s.config["base_path"].(string), req.Cluster, req.Node)

	// get JSON string ([]byte)
	value, err := s.store.Get(key)
	if err != nil {
		return nil, err
	}

	// JSON string -> NodeInfo
	nodeInfo := &protobuf.NodeInfo{}
	unmarshaler := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	err = unmarshaler.Unmarshal(bytes.NewReader(value), nodeInfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &protobuf.GetNodeResponse{
		NodeInfo: nodeInfo,
	}, nil
}

func (s *ClusterService) DeleteNode(ctx context.Context, req *protobuf.DeleteNodeRequest) (*empty.Empty, error) {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", s.config["base_path"].(string), req.Cluster, req.Node)

	err := s.store.Delete(key)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ClusterService) PutIndexMapping(ctx context.Context, req *protobuf.PutIndexMappingRequest) (*empty.Empty, error) {
	key := fmt.Sprintf("%s/clusters/%s/index_mapping.json", s.config["base_path"].(string), req.Cluster)

	// Any -> IndexMappingImpl
	v, err := protobuf.UnmarshalAny(req.IndexMapping)
	if err != nil {
		return nil, err
	}
	indexMapping := v.(*mapping.IndexMappingImpl)

	// IndexMappingImple -> JSON string ([]byte)
	value, err := json.MarshalIndent(indexMapping, "", "  ")
	if err != nil {
		return nil, err
	}

	// put JSON string ([]byte)
	err = s.store.Put(key, []byte(value))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ClusterService) GetIndexMapping(ctx context.Context, req *protobuf.GetIndexMappingRequest) (*protobuf.GetIndexMappingResponse, error) {

	return &protobuf.GetIndexMappingResponse{}, nil
}

func (s *ClusterService) DeleteIndexMapping(ctx context.Context, req *protobuf.DeleteIndexMappingRequest) (*empty.Empty, error) {

	return &empty.Empty{}, nil
}

func (s *ClusterService) PutIndexMeta(ctx context.Context, req *protobuf.PutIndexMetaRequest) (*empty.Empty, error) {

	return &empty.Empty{}, nil
}

func (s *ClusterService) GetIndexMeta(ctx context.Context, req *protobuf.GetIndexMetaRequest) (*protobuf.GetIndexMetaResponse, error) {

	return &protobuf.GetIndexMetaResponse{}, nil
}

func (s *ClusterService) DeleteIndexMeta(ctx context.Context, req *protobuf.DeleteIndexMetaRequest) (*empty.Empty, error) {

	return &empty.Empty{}, nil
}
