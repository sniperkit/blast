package file

import (
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

	err = store.Put("test_key", []byte("test_value"))
	if err != nil {
		t.Fatalf("failed to put node to store. %v", err)
	}

	value, err := store.Get("test_key")
	if err != nil {
		t.Fatalf("failed to get node from store. %v", err)
	}
	if string(value) != "test_value" {
		t.Fatalf("got %v\nwant %v", string(value), "test_value")
	}

	err = store.Delete("test_key")
	if err != nil {
		t.Fatalf("failed to delete node from store. %v", err)
	}
}
