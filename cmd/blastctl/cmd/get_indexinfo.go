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
	"github.com/mosuka/blast/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	indexPath         bool
	indexMapping      bool
	indexType         bool
	kvstore           bool
	kvconfig          bool
}

var getIndexCmdOpts = GetIndexCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
	indexPath:         false,
	indexMapping:      false,
	indexType:         false,
	kvstore:           false,
	kvconfig:          false,
}

var getIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "gets the index information",
	Long:  `The get index command gets the index information.`,
	RunE:  runEGetIndexCmd,
}

func runEGetIndexCmd(cmd *cobra.Command, args []string) error {
	if !getIndexCmdOpts.indexPath && !getIndexCmdOpts.indexMapping && !getIndexCmdOpts.indexType && !getIndexCmdOpts.kvstore && !getIndexCmdOpts.kvconfig {
		getIndexCmdOpts.indexPath = true
		getIndexCmdOpts.indexMapping = true
		getIndexCmdOpts.indexType = true
		getIndexCmdOpts.kvstore = true
		getIndexCmdOpts.kvconfig = true
	}

	// create client
	c, err := client.NewGRPCClient(context.Background(), getIndexCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getDocumentCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	// get document from index
	indexPath, indexMapping, indexType, kvstore, kvconfig, err := c.GetIndexInfo(ctx, getIndexCmdOpts.indexPath, getIndexCmdOpts.indexMapping, getIndexCmdOpts.indexType, getIndexCmdOpts.kvstore, getIndexCmdOpts.kvconfig)
	resp := struct {
		IndexPath    string                    `json:"index_path,omitempty"`
		IndexMapping *mapping.IndexMappingImpl `json:"index_mapping,omitempty"`
		IndexType    string                    `json:"index_type,omitempty"`
		Kvstore      string                    `json:"kvstore,omitempty"`
		Kvconfig     map[string]interface{}    `json:"kvconfig,omitempty"`
		Error        error                     `json:"error,omitempty"`
	}{
		IndexPath:    indexPath,
		IndexMapping: indexMapping,
		IndexType:    indexType,
		Kvstore:      kvstore,
		Kvconfig:     kvconfig,
		Error:        err,
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
	getIndexCmd.Flags().SortFlags = false

	getIndexCmd.Flags().StringVar(&getIndexCmdOpts.grpcServerAddress, "grpc-server-address", getIndexCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getIndexCmd.Flags().IntVar(&getIndexCmdOpts.dialTimeout, "dial-timeout", getIndexCmdOpts.dialTimeout, "dial timeout")
	getIndexCmd.Flags().IntVar(&getIndexCmdOpts.requestTimeout, "request-timeout", getIndexCmdOpts.requestTimeout, "request timeout")
	getIndexCmd.Flags().BoolVar(&getIndexCmdOpts.indexPath, "index-path", getIndexCmdOpts.indexPath, "include index path")
	getIndexCmd.Flags().BoolVar(&getIndexCmdOpts.indexMapping, "index-mapping", getIndexCmdOpts.indexMapping, "include index mapping")
	getIndexCmd.Flags().BoolVar(&getIndexCmdOpts.indexType, "index-type", getIndexCmdOpts.indexType, "include index type")
	getIndexCmd.Flags().BoolVar(&getIndexCmdOpts.kvstore, "kvstore", getIndexCmdOpts.kvstore, "include kvstore")
	getIndexCmd.Flags().BoolVar(&getIndexCmdOpts.kvconfig, "kvconfig", getIndexCmdOpts.kvconfig, "include kvconfig")

	getCmd.AddCommand(getIndexCmd)
}
