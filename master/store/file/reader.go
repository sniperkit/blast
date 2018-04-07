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

package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Reader struct {
	store *Store
}

func NewReader(s *Store) (Reader, error) {
	return Reader{
		store: s,
	}, nil
}

func (r *Reader) Read(key string) (map[string]interface{}, error) {
	file, err := os.Open(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, data)

	return data, nil
}

func (r *Reader) GetNode(cluster string, node string) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", r.store.BasePath, cluster, node)

	value, err := r.Read(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (r *Reader) Close() error {
	return nil
}