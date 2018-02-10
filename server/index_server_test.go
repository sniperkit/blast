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

package server

import (
	"github.com/mosuka/blast/index"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestBlastServer(t *testing.T) {
	dir, _ := os.Getwd()

	listenAddress := "localhost:0"
	indexPath, _ := ioutil.TempDir("/tmp", "blast")
	indexPath = indexPath + "/index/data"
	indexMappingPath := dir + "/../etc/index_mapping.json"
	indexMetaPath := dir + "/../etc/index_meta.json"

	indexMappingFile, err := os.Open(indexMappingPath)
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}
	defer indexMappingFile.Close()

	indexMapping, err := index.LoadIndexMapping(indexMappingFile)
	if err != nil {
		t.Errorf("could not load IndexMapping: %v", err)
	}

	indexMetaFile, err := os.Open(indexMetaPath)
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}
	defer indexMetaFile.Close()

	indexMeta, err := index.LoadIndexMeta(indexMetaFile)
	if err != nil {
		t.Errorf("could not load kvconfig %v", err)
	}

	blastServer, err := NewIndexServer(listenAddress, indexPath, indexMapping, indexMeta)
	if err != nil {
		t.Fatalf("unexpected error. expected not nil, actual %v", blastServer)
	}

	err = blastServer.Start()

	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}

	time.Sleep(10 * time.Second)

	err = blastServer.Stop()
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}

	os.RemoveAll(indexPath)
}
