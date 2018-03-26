package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Reader struct {
	store *Store
}

func NewReader(s *Store) (Reader, error) {
	return Reader{
		store: s,
	}, nil
}

func (r *Reader) Read(key string) (map[string]interface{}, error) {
	file, err := os.Open(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, data)

	return data, nil
}

func (r *Reader) GetNode(cluster string, node string) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s/clusters/%s/nodes/%s.json", r.store.BasePath, cluster, node)

	value, err := r.Read(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (r *Reader) Close() error {
	return nil
}
