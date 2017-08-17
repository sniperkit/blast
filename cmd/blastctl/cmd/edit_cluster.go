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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	"github.com/coreos/etcd/clientv3"
	"github.com/mosuka/blast/util"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type EditClusterCommandOptions struct {
	etcdEndpoints      []string
	etcdDialTimeout    int
	etcdRequestTimeout int
	cluster            string
	indexMapping       string
	indexType          string
	kvstore            string
	kvconfig           string
}

var editClusterCmdOpts = EditClusterCommandOptions{
	etcdEndpoints:      []string{"localhost:2379"},
	etcdDialTimeout:    5000,
	etcdRequestTimeout: 5000,
	cluster:            "",
	indexMapping:       "",
	indexType:          "upside_down",
	kvstore:            "boltdb",
	kvconfig:           "",
}

var editClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "edits the cluster information",
	Long:  `The edit cluster command edits the cluster information.`,
	RunE:  runEEditClusterCmd,
}

func runEEditClusterCmd(cmd *cobra.Command, args []string) error {
	// check cluster name
	if editClusterCmdOpts.cluster == "" {
		return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
	}

	// IndexMapping
	indexMapping := mapping.NewIndexMapping()
	if editClusterCmdOpts.indexMapping != "" {
		file, err := os.Open(editClusterCmdOpts.indexMapping)
		if err != nil {
			return err
		}
		defer file.Close()

		indexMapping, err = util.NewIndexMapping(file)
		if err != nil {
			return err
		}
	}

	// Kvconfig
	kvconfig := make(map[string]interface{})
	if editClusterCmdOpts.kvconfig != "" {
		file, err := os.Open(editClusterCmdOpts.kvconfig)
		if err != nil {
			return err
		}
		defer file.Close()

		kvconfig, err = util.NewKvconfig(file)
		if err != nil {
			return err
		}
	}

	var err error

	var bytesIndexMapping []byte
	if indexMapping != nil {
		bytesIndexMapping, err = json.Marshal(indexMapping)
		if err != nil {
			return err
		}
	}

	var bytesKvconfig []byte
	if kvconfig != nil {
		bytesKvconfig, err = json.Marshal(kvconfig)
		if err != nil {
			return err
		}
	}

	cfg := clientv3.Config{
		Endpoints:   editClusterCmdOpts.etcdEndpoints,
		DialTimeout: time.Duration(editClusterCmdOpts.etcdDialTimeout) * time.Millisecond,
		Context:     context.Background(),
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return err
	}
	defer c.Close()

	var kv clientv3.KV
	if c != nil {
		kv = clientv3.NewKV(c)
	}

	resp := struct {
		IndexMapping *mapping.IndexMappingImpl `json:"index_mapping,omitempty"`
		IndexType    string                    `json:"index_type,omitempty"`
		Kvstore      string                    `json:"kvstore,omitempty"`
		Kvconfig     map[string]interface{}    `json:"kvconfig,omitempty"`
	}{}

	if cmd.Flag("index-mapping").Changed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(editClusterCmdOpts.etcdRequestTimeout)*time.Millisecond)
		defer cancel()

		keyIndexMapping := fmt.Sprintf("/blast/clusters/%s/indexMapping", editClusterCmdOpts.cluster)

		_, err = kv.Put(ctx, keyIndexMapping, string(bytesIndexMapping))
		if err != nil {
			return err
		}
		resp.IndexMapping = indexMapping
	}

	if cmd.Flag("index-type").Changed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(editClusterCmdOpts.etcdRequestTimeout)*time.Millisecond)
		defer cancel()

		keyIndexType := fmt.Sprintf("/blast/clusters/%s/indexType", editClusterCmdOpts.cluster)

		_, err = kv.Put(ctx, keyIndexType, editClusterCmdOpts.indexType)
		if err != nil {
			return err
		}
		resp.IndexType = editClusterCmdOpts.indexType
	}

	if cmd.Flag("kvstore").Changed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(editClusterCmdOpts.etcdRequestTimeout)*time.Millisecond)
		defer cancel()

		keyKvstore := fmt.Sprintf("/blast/clusters/%s/kvstore", editClusterCmdOpts.cluster)

		_, err = kv.Put(ctx, keyKvstore, editClusterCmdOpts.kvstore)
		if err != nil {
			return err
		}
		resp.Kvstore = editClusterCmdOpts.kvstore
	}

	if cmd.Flag("kvconfig").Changed {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(editClusterCmdOpts.etcdRequestTimeout)*time.Millisecond)
		defer cancel()

		keyKvconfig := fmt.Sprintf("/blast/clusters/%s/kvconfig", editClusterCmdOpts.cluster)

		_, err = kv.Put(ctx, keyKvconfig, string(bytesKvconfig))
		if err != nil {
			return err
		}
		resp.Kvconfig = kvconfig
	}

	// output response
	switch rootCmdOpts.outputFormat {
	case "text":
		fmt.Printf("%v\n", resp)
	case "json":
		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", output)
	default:
		fmt.Printf("%v\n", resp)
	}

	return nil
}

func init() {
	editClusterCmd.Flags().SortFlags = false

	editClusterCmd.Flags().StringSliceVar(&editClusterCmdOpts.etcdEndpoints, "etcd-endpoint", editClusterCmdOpts.etcdEndpoints, "etcd eendpoint")
	editClusterCmd.Flags().IntVar(&editClusterCmdOpts.etcdDialTimeout, "etcd-dial-timeout", editClusterCmdOpts.etcdDialTimeout, "etcd dial timeout")
	editClusterCmd.Flags().IntVar(&editClusterCmdOpts.etcdRequestTimeout, "etcd-request-timeout", editClusterCmdOpts.etcdRequestTimeout, "etcd request timeout")
	editClusterCmd.Flags().StringVar(&editClusterCmdOpts.cluster, "cluster", editClusterCmdOpts.cluster, "cluster name")
	editClusterCmd.Flags().StringVar(&editClusterCmdOpts.indexMapping, "index-mapping", editClusterCmdOpts.indexMapping, "index mapping")
	editClusterCmd.Flags().StringVar(&editClusterCmdOpts.indexType, "index-type", editClusterCmdOpts.indexType, "index type")
	editClusterCmd.Flags().StringVar(&editClusterCmdOpts.kvstore, "kvstore", editClusterCmdOpts.kvstore, "kvstore")
	editClusterCmd.Flags().StringVar(&editClusterCmdOpts.kvconfig, "kvconfig", editClusterCmdOpts.kvconfig, "kvconfig")

	editCmd.AddCommand(editClusterCmd)
}
