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
	"github.com/buger/jsonparser"
	"github.com/mosuka/blast/config"
	blastgrpc "github.com/mosuka/blast/node/client/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type BulkCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	batchSize         int32
	request           string
	outputFormat      string
}

var bulkCmdOpts = BulkCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       60000,
	requestTimeout:    60000,
	batchSize:         1000,
	request:           "",
	outputFormat:      "json",
}

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "puts or deletes documents in bulk",
	Long:  `The bulk command puts or deletes documents in bulk.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read request
		data := []byte(bulkCmdOpts.request)

		// get batch_size from request
		batchSize, err := jsonparser.GetInt(data, "batch_size")
		if err != nil {
			fmt.Println("batch_size")
			return err
		}

		// get requests from request
		requestsBytes, _, _, err := jsonparser.Get(data, "requests")
		if err != nil {
			fmt.Println("requests")
			return err
		}
		var requests []map[string]interface{}
		err = json.Unmarshal(requestsBytes, &requests)
		if err != nil {
			return err
		}

		// overwrite batch size by command line option
		if cmd.Flag("batch-size").Changed {
			batchSize = int64(bulkCmdOpts.batchSize)
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), bulkCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(bulkCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		// update documents to index in bulk
		putCount, deleteCount, errorCount, err := c.Bulk(ctx, requests, int32(batchSize))
		resp := struct {
			PutCount    int32 `json:"put_count,omitempty"`
			DeleteCount int32 `json:"delete_count,omitempty"`
			ErrorCount  int32 `json:"error_count,omitempty"`
			Error       error `json:"error,omitempty"`
		}{
			PutCount:    putCount,
			DeleteCount: deleteCount,
			ErrorCount:  errorCount,
			Error:       err,
		}

		// output request
		switch bulkCmdOpts.outputFormat {
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
	bulkCmd.Flags().SortFlags = false

	bulkCmd.Flags().StringVar(&bulkCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	bulkCmd.Flags().IntVar(&bulkCmdOpts.dialTimeout, "dial-timeout", bulkCmdOpts.dialTimeout, "dial timeout")
	bulkCmd.Flags().IntVar(&bulkCmdOpts.requestTimeout, "request-timeout", bulkCmdOpts.requestTimeout, "request timeout")
	bulkCmd.Flags().Int32Var(&bulkCmdOpts.batchSize, "batch-size", bulkCmdOpts.batchSize, "batch size of bulk request")
	bulkCmd.Flags().StringVar(&bulkCmdOpts.request, "request", bulkCmdOpts.request, "request file")
	bulkCmd.Flags().StringVar(&bulkCmdOpts.outputFormat, "output-format", bulkCmdOpts.outputFormat, "output format")

	RootCmd.AddCommand(bulkCmd)
}
