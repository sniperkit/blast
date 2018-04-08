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

	err = store.PutNode("test_cluster", "test_node")
	if err != nil {
		t.Fatalf("failed to put node to store. %v", err)
	}

	data, err := store.GetNode("test_cluster", "test_node")
	if err != nil {
		t.Fatalf("failed to get node from store. %v", err)
	}

	fmt.Println(data)

	err = store.DeleteNode("test_cluster", "test_node")
	if err != nil {
		t.Fatalf("failed to delete node from store. %v", err)
	}
}
