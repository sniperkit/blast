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

package cluster

import (
	"encoding/json"
	"github.com/mosuka/blast/master/store/file"
	"io"
	"io/ioutil"
)

type ClusterMeta struct {
	Storage string                 `json:"storage"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func NewClusterMeta() *ClusterMeta {
	return &ClusterMeta{
		Storage: file.Name,
		Config:  make(map[string]interface{}),
	}
}

func LoadClusterMeta(reader io.Reader) (*ClusterMeta, error) {
	clusterMeta := NewClusterMeta()

	resourceBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resourceBytes, clusterMeta)
	if err != nil {
		return nil, err
	}

	return clusterMeta, nil
}
