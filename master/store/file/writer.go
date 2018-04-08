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
	"os"
	"path/filepath"
	"time"
)

type Writer struct {
	store *Store
}

func NewWriter(s *Store) (Writer, error) {
	return Writer{
		store: s,
	}, nil
}

func (w *Writer) Write(key string, value map[string]interface{}) error {
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

func (w *Writer) PutNode(cluster string, node string) error {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", w.store.BasePath, cluster, node)

	value := make(map[string]interface{})

	value["timestamp"] = time.Now().Format(time.RFC3339Nano)

	return w.Write(key, value)
}

func (w *Writer) DeleteNode(cluster string, node string) error {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", w.store.BasePath, cluster, node)

	err := os.Remove(key)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Close() error {
	return nil
}
