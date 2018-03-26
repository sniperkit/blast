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
	"github.com/mosuka/blast/master/store/file"
	"io"
	"io/ioutil"
)

type SupervisorConfig struct {
	Storage string                 `json:"storage"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func NewSupervisorConfig() *SupervisorConfig {
	return &SupervisorConfig{
		Storage: file.Name,
		Config:  make(map[string]interface{}),
	}
}

func LoadSupervisorConfig(reader io.Reader) (*SupervisorConfig, error) {
	supervisorConfig := NewSupervisorConfig()

	resourceBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resourceBytes, supervisorConfig)
	if err != nil {
		return nil, err
	}

	return supervisorConfig, nil
}
