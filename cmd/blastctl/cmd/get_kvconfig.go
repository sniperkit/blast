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

type GetKvconfigCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
}

var getKvconfigCmdOpts = GetKvconfigCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
}

var getKvconfigCmd = &cobra.Command{
	Use:   "kvconfig",
	Short: "gets the kvconfig",
	Long:  `The get kvconfig command gets the kvconfig.`,
	RunE:  runEGetKvconfigCmd,
}

func runEGetKvconfigCmd(cmd *cobra.Command, args []string) error {
	// create client
	c, err := client.NewIndexClient(context.Background(), getKvconfigCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getKvconfigCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	kvconfig, err := c.GetKvconfig(ctx)
	resp := struct {
		Kvconfig map[string]interface{} `json:"kvconfig,omitempty"`
		Error    error                  `json:"error,omitempty"`
	}{
		Kvconfig: kvconfig,
		Error:    err,
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
	getKvconfigCmd.Flags().SortFlags = false

	getKvconfigCmd.Flags().StringVar(&getKvconfigCmdOpts.grpcServerAddress, "grpc-server-address", getKvconfigCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getKvconfigCmd.Flags().IntVar(&getKvconfigCmdOpts.dialTimeout, "dial-timeout", getKvconfigCmdOpts.dialTimeout, "dial timeout")
	getKvconfigCmd.Flags().IntVar(&getKvconfigCmdOpts.requestTimeout, "request-timeout", getKvconfigCmdOpts.requestTimeout, "request timeout")

	getCmd.AddCommand(getKvconfigCmd)
}
