package grpc

import (
	"github.com/mosuka/blast/node/config"
	"github.com/mosuka/blast/node/server/grpc"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGRPCClient(t *testing.T) {
	dir, _ := os.Getwd()

	listenAddress := "localhost:0"
	indexPath, _ := ioutil.TempDir("/tmp", "blast")
	indexPath = indexPath + "/index/data"
	indexMappingPath := dir + "/../../../etc/index_mapping.json"
	indexMetaPath := dir + "/../../../etc/index_config.json"

	indexMappingFile, err := os.Open(indexMappingPath)
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}
	defer indexMappingFile.Close()

	indexMapping, err := config.LoadIndexMapping(indexMappingFile)
	if err != nil {
		t.Errorf("could not load IndexMapping: %v", err)
	}

	indexMetaFile, err := os.Open(indexMetaPath)
	if err != nil {
		t.Fatalf("unexpected error. %v", err)
	}
	defer indexMetaFile.Close()

	indexMeta, err := config.LoadIndexConfig(indexMetaFile)
	if err != nil {
		t.Errorf("could not load kvconfig %v", err)
	}

	blastServer, err := grpc.NewGRPCServer(listenAddress, indexPath, indexMapping, indexMeta)
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
