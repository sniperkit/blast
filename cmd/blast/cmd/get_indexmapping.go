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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	mastergrpc "github.com/mosuka/blast/master/client/grpc"
	masterconfig "github.com/mosuka/blast/master/config"
	nodegrpc "github.com/mosuka/blast/node/client/grpc"
	nodeconfig "github.com/mosuka/blast/node/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexMappingCmdOpts struct {
	masterAddress  string
	nodeAddress    string
	dialTimeout    int
	requestTimeout int
	cluster        string
}

var getIndexMappingCmdOpts = GetIndexMappingCmdOpts{
	masterAddress:  masterconfig.DefaultGRPCListenAddress,
	nodeAddress:    nodeconfig.DefaultGRPCListenAddress,
	dialTimeout:    5000,
	requestTimeout: 5000,
	cluster:        "",
}

var getIndexMappingCmd = &cobra.Command{
	Use:   "indexmapping",
	Short: "gets the index mapping",
	Long:  `The get indexmapping command gets the index mapping.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp := struct {
			Cluster      string                    `json:"cluster,omitempty"`
			IndexMapping *mapping.IndexMappingImpl `json:"index_mapping,omitempty"`
			Error        error                     `json:"error,omitempty"`
		}{}

		// check cluster
		if getIndexMappingCmdOpts.cluster == "" {
			// create client
			c, err := nodegrpc.NewGRPCClient(context.Background(), getIndexMappingCmdOpts.nodeAddress, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer c.Close()

			// create context
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexMappingCmdOpts.requestTimeout)*time.Millisecond)
			defer cancel()

			indexMapping, err := c.GetIndexMapping(ctx)
			if err != nil {
				return err
			}

			resp.IndexMapping = indexMapping
		} else {
			c, err := mastergrpc.NewGRPCClient(context.Background(), getIndexMappingCmdOpts.masterAddress, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer c.Close()

			// create context
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexMappingCmdOpts.requestTimeout)*time.Millisecond)
			defer cancel()

			indexMapping, err := c.GetIndexMapping(ctx, getIndexMappingCmdOpts.cluster)
			if err != nil {
				return err
			}

			resp.Cluster = getIndexMappingCmdOpts.cluster
			resp.IndexMapping = indexMapping
		}

		// output response
		switch getCmdOpts.outputFormat {
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
	},
}

func init() {
	getIndexMappingCmd.Flags().SortFlags = false

	getIndexMappingCmd.Flags().StringVar(&getIndexMappingCmdOpts.masterAddress, "master-address", masterconfig.DefaultGRPCListenAddress, "Blast master to connect to using gRPC")
	getIndexMappingCmd.Flags().StringVar(&getIndexMappingCmdOpts.nodeAddress, "node-address", nodeconfig.DefaultGRPCListenAddress, "Blast node to connect to using gRPC")
	getIndexMappingCmd.Flags().IntVar(&getIndexMappingCmdOpts.dialTimeout, "dial-timeout", getIndexMappingCmdOpts.dialTimeout, "dial timeout")
	getIndexMappingCmd.Flags().IntVar(&getIndexMappingCmdOpts.requestTimeout, "request-timeout", getIndexMappingCmdOpts.requestTimeout, "request timeout")
	getIndexMappingCmd.Flags().StringVar(&getIndexMappingCmdOpts.cluster, "cluster", getIndexMappingCmdOpts.cluster, "cluster")

	getCmd.AddCommand(getIndexMappingCmd)
}
