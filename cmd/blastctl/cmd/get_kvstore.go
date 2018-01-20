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

type GetKvstoreCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
}

var getKvstoreCmdOpts = GetKvstoreCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
}

var getKvstoreCmd = &cobra.Command{
	Use:   "kvstore",
	Short: "gets the kvstore",
	Long:  `The get kvstore command gets the kvstore.`,
	RunE:  runEGetKvstoreCmd,
}

func runEGetKvstoreCmd(cmd *cobra.Command, args []string) error {
	// create client
	c, err := client.NewIndexClient(context.Background(), getKvstoreCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getKvstoreCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	kvstore, err := c.GetKvstore(ctx)
	resp := struct {
		Kvstore string `json:"kvstore,omitempty"`
		Error   error  `json:"error,omitempty"`
	}{
		Kvstore: kvstore,
		Error:   err,
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
	getKvstoreCmd.Flags().SortFlags = false

	getKvstoreCmd.Flags().StringVar(&getKvstoreCmdOpts.grpcServerAddress, "grpc-server-address", getKvstoreCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getKvstoreCmd.Flags().IntVar(&getKvstoreCmdOpts.dialTimeout, "dial-timeout", getKvstoreCmdOpts.dialTimeout, "dial timeout")
	getKvstoreCmd.Flags().IntVar(&getKvstoreCmdOpts.requestTimeout, "request-timeout", getKvstoreCmdOpts.requestTimeout, "request timeout")

	getCmd.AddCommand(getKvstoreCmd)
}
