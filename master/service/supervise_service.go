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
	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/mosuka/blast/master/builtin"
	"github.com/mosuka/blast/master/config"
	"github.com/mosuka/blast/master/registry"
	"github.com/mosuka/blast/master/store"
	"github.com/mosuka/blast/protobuf"
	"golang.org/x/net/context"
)

type SuperviseService struct {
	storeType string
	config    map[string]interface{}
	store     store.Store
}

func NewSuperviseService(svMeta *config.SupervisorConfig) (*SuperviseService, error) {
	storeConstructor := registry.StoreConstructorByName(svMeta.Storage)

	store, err := storeConstructor(svMeta.Config)
	if err != nil {
		return nil, err
	}

	return &SuperviseService{
		storeType: svMeta.Storage,
		config:    svMeta.Config,
		store:     store,
	}, nil
}

func (s *SuperviseService) PutNode(ctx context.Context, req *protobuf.PutNodeRequest) (*empty.Empty, error) {
	err := s.store.PutNode(req.Cluster, req.Node)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *SuperviseService) GetNode(ctx context.Context, req *protobuf.GetNodeRequest) (*protobuf.GetNodeResponse, error) {
	value, err := s.store.GetNode(req.Cluster, req.Node)
	if err != nil {
		return nil, err
	}

	valueAny, err := protobuf.MarshalAny(value)
	if err != nil {
		return nil, err
	}

	return &protobuf.GetNodeResponse{
		Value: &valueAny,
	}, nil
}

func (s *SuperviseService) DeleteNode(ctx context.Context, req *protobuf.DeleteNodeRequest) (*empty.Empty, error) {
	err := s.store.DeleteNode(req.Cluster, req.Node)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
