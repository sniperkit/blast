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
	"os"
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
	file, err := os.Open(key)
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
