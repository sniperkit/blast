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

package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type IndexConfig struct {
	IndexType string                 `json:"index_type"`
	Storage   string                 `json:"storage"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

func NewIndexConfig() *IndexConfig {
	return &IndexConfig{
		IndexType: "upside_down",
		Storage:   "boltdb",
		Config:    make(map[string]interface{}),
	}
}

func LoadIndexConfig(reader io.Reader) (*IndexConfig, error) {
	indexMeta := NewIndexConfig()

	resourceBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resourceBytes, indexMeta)
	if err != nil {
		return nil, err
	}

	return indexMeta, nil
}
