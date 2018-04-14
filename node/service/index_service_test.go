//  Copyright (c) 2017 Minoru Osuka
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
	"github.com/mosuka/blast/index"
	"io/ioutil"
	"os"
	"testing"
)

func TestIndexService(t *testing.T) {
	dir, _ := os.Getwd()

	indexPath, _ := ioutil.TempDir("/tmp", "indigo")
	indexMappingPath := dir + "/../../example/config/index_mapping.json"
	indexMetaPath := dir + "/../../example/config/index_meta.json"

	indexMappingFile, err := os.Open(indexMappingPath)
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}
	defer indexMappingFile.Close()

	indexMapping, err := index.LoadIndexMapping(indexMappingFile)
	if err != nil {
		t.Errorf("could not load IndexMapping : %v", err)
	}

	indexMetaFile, err := os.Open(indexMetaPath)
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}
	defer indexMetaFile.Close()

	indexMeta, err := index.LoadIndexMeta(indexMetaFile)
	if err != nil {
		t.Errorf("could not load kvconfig : %v", err)
	}

	s, err := NewIndexService(indexPath, indexMapping, indexMeta)
	if err != nil {
		t.Fatalf("unexpected error.  expected not nil, actual %v", s)
	}

	if s.indexPath != indexPath {
		t.Errorf("unexpected error.  expected %v, actual %v", indexPath, s.indexPath)
	}
}
