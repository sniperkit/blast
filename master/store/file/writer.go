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
