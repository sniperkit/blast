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
	"github.com/mosuka/blast/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexTypeCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
}

var getIndexTypeCmdOpts = GetIndexTypeCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
}

var getIndexTypeCmd = &cobra.Command{
	Use:   "indextype",
	Short: "gets the index type",
	Long:  `The get index command gets the index type.`,
	RunE:  runEGetIndexTypeCmd,
}

func runEGetIndexTypeCmd(cmd *cobra.Command, args []string) error {
	// create client
	c, err := client.NewGRPCClient(context.Background(), getIndexTypeCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexTypeCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	indexType, err := c.GetIndexType(ctx)
	resp := struct {
		IndexType string `json:"index_type,omitempty"`
		Error     error  `json:"error,omitempty"`
	}{
		IndexType: indexType,
		Error:     err,
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
	getIndexTypeCmd.Flags().SortFlags = false

	getIndexTypeCmd.Flags().StringVar(&getIndexTypeCmdOpts.grpcServerAddress, "grpc-server-address", getIndexTypeCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getIndexTypeCmd.Flags().IntVar(&getIndexTypeCmdOpts.dialTimeout, "dial-timeout", getIndexTypeCmdOpts.dialTimeout, "dial timeout")
	getIndexTypeCmd.Flags().IntVar(&getIndexTypeCmdOpts.requestTimeout, "request-timeout", getIndexTypeCmdOpts.requestTimeout, "request timeout")

	getCmd.AddCommand(getIndexTypeCmd)
}
