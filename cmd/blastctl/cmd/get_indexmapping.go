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
	"github.com/mosuka/blast/node/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexMappingCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
}

var getIndexMappingCmdOpts = GetIndexMappingCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
}

var getIndexMappingCmd = &cobra.Command{
	Use:   "indexmapping",
	Short: "gets the index mapping",
	Long:  `The get index command gets the index mapping.`,
	RunE:  runEGetIndexMappingCmd,
}

func runEGetIndexMappingCmd(cmd *cobra.Command, args []string) error {
	// create client
	c, err := client.NewGRPCClient(context.Background(), getIndexMappingCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexMappingCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	indexMapping, err := c.GetIndexMapping(ctx)
	resp := struct {
		IndexMapping *mapping.IndexMappingImpl `json:"index_mapping,omitempty"`
		Error        error                     `json:"error,omitempty"`
	}{
		IndexMapping: indexMapping,
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
	getIndexMappingCmd.Flags().SortFlags = false

	getIndexMappingCmd.Flags().StringVar(&getIndexMappingCmdOpts.grpcServerAddress, "grpc-server-address", getIndexMappingCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getIndexMappingCmd.Flags().IntVar(&getIndexMappingCmdOpts.dialTimeout, "dial-timeout", getIndexMappingCmdOpts.dialTimeout, "dial timeout")
	getIndexMappingCmd.Flags().IntVar(&getIndexMappingCmdOpts.requestTimeout, "request-timeout", getIndexMappingCmdOpts.requestTimeout, "request timeout")

	getCmd.AddCommand(getIndexMappingCmd)
}
