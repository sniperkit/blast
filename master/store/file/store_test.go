package file

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestStore(t *testing.T) {
	configPath, _ := ioutil.TempDir("/tmp", "blast")
	configPath = configPath + "/etc/blast"

	config := map[string]interface{}{
		"base_path": configPath,
	}

	store, err := NewStore(config)
	if err != nil {
		t.Fatalf("failed to create file store. %v", err)
	}

	writer, err := store.Writer()
	if err != nil {
		t.Fatalf("failed to get store writer. %v", err)
	}

	err = writer.PutNode("test_cluster", "test_node")
	if err != nil {
		t.Fatalf("failed to put node to store. %v", err)
	}

	reader, err := store.Reader()
	if err != nil {
		t.Fatalf("failed to get store reader. %v", err)
	}

	data, err := reader.GetNode("test_cluster", "test_node")
	if err != nil {
		t.Fatalf("failed to get node from store. %v", err)
	}

	fmt.Println(data)
}
