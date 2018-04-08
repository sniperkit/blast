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
	"errors"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/master/registry"
	"github.com/mosuka/blast/master/store"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	Name            = "file"
	DefaultBasePath = "./etc/blast"
)

func init() {
	registry.RegisterStore(Name, NewStore)
}

type Store struct {
	BasePath string
}

func NewStore(config map[string]interface{}) (store.Store, error) {
	path, ok := config["base_path"].(string)
	if !ok || path == "" {
		path = DefaultBasePath
	}

	store := Store{
		BasePath: path,
	}

	return &store, nil
}

func (s *Store) Put(key string, value map[string]interface{}) error {
	dir := filepath.Dir(key)

	// check directory
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	// check file
	_, err = os.Stat(key)
	if os.IsNotExist(err) {
		file, err := os.Create(key)
		if err != nil {
			return err
		}
		defer file.Close()

		jsonBytes, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}

		_, err = file.Write(jsonBytes)
		if err != nil {
			return err
		}
	} else {
		return errors.New("file already exists")
	}

	return nil
}

func (s *Store) Get(key string) (map[string]interface{}, error) {
	// check file
	_, err := os.Stat(key)
	if os.IsNotExist(err) {
		return nil, err
	}

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
	err = json.Unmarshal(jsonBytes, &data)

	return data, nil
}

func (s *Store) Delete(key string) error {
	// check file
	_, err := os.Stat(key)
	if os.IsNotExist(err) {
		return err
	}

	err = os.Remove(key)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) PutNode(cluster string, node string) error {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", s.BasePath, cluster, node)

	value := make(map[string]interface{})

	value["timestamp"] = time.Now().Format(time.RFC3339Nano)

	return s.Put(key, value)
}

func (s *Store) GetNode(cluster string, node string) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", s.BasePath, cluster, node)

	value, err := s.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *Store) DeleteNode(cluster string, node string) error {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", s.BasePath, cluster, node)

	err := s.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) PutIndexMapping(cluster string, indexMapping string) error {
	key := fmt.Sprintf("%s/clusters/%s/index_mapping.json", s.BasePath, cluster)

	s.Put(key, nil)

	return nil
}

func (s *Store) GetIndexMapping(cluter string) (*mapping.IndexMappingImpl, error) {

	return nil, nil
}

func (s *Store) DeleteIndexMapping(cluster string) error {

	return nil
}

func (s *Store) Close() error {
	return nil
}
