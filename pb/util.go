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

package pb

import (
	"encoding/json"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"reflect"
)

var (
	typeRegistry = make(map[string]reflect.Type)
)

func init() {
	typeRegistry["mapping.IndexMappingImpl"] = reflect.TypeOf(mapping.IndexMappingImpl{})
	typeRegistry["bleve.IndexStat"] = reflect.TypeOf(bleve.IndexStat{})
	typeRegistry["interface {}"] = reflect.TypeOf((map[string]interface{})(nil))
	typeRegistry["bleve.SearchRequest"] = reflect.TypeOf(bleve.SearchRequest{})
	typeRegistry["bleve.SearchResult"] = reflect.TypeOf(bleve.SearchResult{})
}

func MarshalAny(instance interface{}) (google_protobuf.Any, error) {
	var message google_protobuf.Any

	if instance == nil {
		return message, nil
	}

	value, err := json.Marshal(instance)
	if err != nil {
		return message, err
	}

	message.TypeUrl = reflect.TypeOf(instance).Elem().String()
	message.Value = value

	return message, nil
}

func UnmarshalAny(message *google_protobuf.Any) (interface{}, error) {
	if message == nil {
		return nil, nil
	}

	typeUrl := message.TypeUrl
	value := message.Value

	instance := reflect.New(typeRegistry[typeUrl]).Interface()

	err := json.Unmarshal(value, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
