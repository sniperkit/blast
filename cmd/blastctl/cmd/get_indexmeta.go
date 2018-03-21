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
	"github.com/mosuka/blast/index/client"
	"github.com/mosuka/blast/index/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexMetaCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
}

var getIndexMetaCmdOpts = GetIndexMetaCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       5000,
	requestTimeout:    5000,
}

var getIndexMetaCmd = &cobra.Command{
	Use:   "indexmeta",
	Short: "gets the index meta",
	Long:  `The get index command gets the index meta.`,
	RunE:  runEGetIndexMetaCmd,
}

func runEGetIndexMetaCmd(cmd *cobra.Command, args []string) error {
	// create client
	c, err := client.NewGRPCClient(context.Background(), getIndexMetaCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexMetaCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	indexMeta, err := c.GetIndexMeta(ctx)
	resp := struct {
		IndexMeta *config.IndexConfig `json:"index_meta,omitempty"`
		Error     error               `json:"error,omitempty"`
	}{
		IndexMeta: indexMeta,
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
	getIndexMetaCmd.Flags().SortFlags = false

	getIndexMetaCmd.Flags().StringVar(&getIndexMetaCmdOpts.grpcServerAddress, "grpc-server-address", getIndexMetaCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	getIndexMetaCmd.Flags().IntVar(&getIndexMetaCmdOpts.dialTimeout, "dial-timeout", getIndexMetaCmdOpts.dialTimeout, "dial timeout")
	getIndexMetaCmd.Flags().IntVar(&getIndexMetaCmdOpts.requestTimeout, "request-timeout", getIndexMetaCmdOpts.requestTimeout, "request timeout")

	getCmd.AddCommand(getIndexMetaCmd)
}
