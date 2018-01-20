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
	"github.com/mosuka/blast/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexPathCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
}

var getIndexPathCmdOpts = GetIndexPathCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
}

var getIndexPathCmd = &cobra.Command{
	Use:   "indexpath",
	Short: "gets the index path",
	Long:  `The get index command gets the index path.`,
	RunE:  runEGetIndexPathCmd,
}

func runEGetIndexPathCmd(cmd *cobra.Command, args []string) error {
	// create client
	c, err := client.NewIndexClient(context.Background(), getIndexPathCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexPathCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	indexPath, err := c.GetIndexPath(ctx)
	resp := struct {
		IndexPath string `json:"index_path,omitempty"`
		Error     error  `json:"error,omitempty"`
	}{
		IndexPath: indexPath,
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
	getIndexPathCmd.Flags().SortFlags = false

	getIndexPathCmd.Flags().StringVar(&getIndexPathCmdOpts.grpcServerAddress, "grpc-server-address", getIndexPathCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getIndexPathCmd.Flags().IntVar(&getIndexPathCmdOpts.dialTimeout, "dial-timeout", getIndexPathCmdOpts.dialTimeout, "dial timeout")
	getIndexPathCmd.Flags().IntVar(&getIndexPathCmdOpts.requestTimeout, "request-timeout", getIndexPathCmdOpts.requestTimeout, "request timeout")

	getCmd.AddCommand(getIndexPathCmd)
}
