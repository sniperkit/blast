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
	"github.com/mosuka/blast/config"
	mastergrpc "github.com/mosuka/blast/master/client/grpc"
	nodegrpc "github.com/mosuka/blast/node/client/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexMappingCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
}

var getIndexMappingCmdOpts = GetIndexMappingCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
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

		var indexMapping *mapping.IndexMappingImpl
		var err error
		if getIndexMappingCmdOpts.cluster != "" {
			indexMapping, err = getIndexMappingFromMaster(getIndexMappingCmdOpts.grpcServerAddress, getIndexMappingCmdOpts.requestTimeout, getIndexMappingCmdOpts.cluster)
			if err != nil {
				return err
			}
			resp.Cluster = getIndexMappingCmdOpts.cluster
			resp.IndexMapping = indexMapping
		} else {
			indexMapping, err = getIndexMappingFromNode(getIndexMappingCmdOpts.grpcServerAddress, getIndexMappingCmdOpts.requestTimeout)
			if err != nil {
				return err
			}
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

func getIndexMappingFromNode(nodeAddress string, requestTimeout int) (*mapping.IndexMappingImpl, error) {
	// create client
	c, err := nodegrpc.NewGRPCClient(context.Background(), nodeAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(requestTimeout)*time.Millisecond)
	defer cancel()

	indexMapping, err := c.GetIndexMapping(ctx)
	if err != nil {
		return nil, err
	}

	return indexMapping, nil
}

func getIndexMappingFromMaster(masterAddress string, requestTimeout int, cluster string) (*mapping.IndexMappingImpl, error) {
	// create client
	c, err := mastergrpc.NewGRPCClient(context.Background(), masterAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(requestTimeout)*time.Millisecond)
	defer cancel()

	indexMapping, err := c.GetIndexMapping(ctx, cluster)
	if err != nil {
		return nil, err
	}

	return indexMapping, nil
}

func init() {
	getIndexMappingCmd.Flags().SortFlags = false

	getIndexMappingCmd.Flags().StringVar(&getIndexMappingCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	getIndexMappingCmd.Flags().IntVar(&getIndexMappingCmdOpts.dialTimeout, "dial-timeout", getIndexMappingCmdOpts.dialTimeout, "dial timeout")
	getIndexMappingCmd.Flags().IntVar(&getIndexMappingCmdOpts.requestTimeout, "request-timeout", getIndexMappingCmdOpts.requestTimeout, "request timeout")
	getIndexMappingCmd.Flags().StringVar(&getIndexMappingCmdOpts.cluster, "cluster", getIndexMappingCmdOpts.cluster, "cluster name. only used in connect to Blast master")

	getCmd.AddCommand(getIndexMappingCmd)
}
