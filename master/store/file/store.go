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
	"github.com/mosuka/blast/master/registry"
	"github.com/mosuka/blast/master/store"
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

func (s *Store) Reader() (store.Reader, error) {
	reader, err := NewReader(s)
	if err != nil {
		return nil, err
	}

	return &reader, nil
}

func (s *Store) Writer() (store.Writer, error) {
	writer, err := NewWriter(s)
	if err != nil {
		return nil, err
	}

	return &writer, nil
}

func (s *Store) Close() error {
	return nil
}
