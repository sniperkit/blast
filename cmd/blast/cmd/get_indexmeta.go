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
	"github.com/mosuka/blast/index"
	blastgrpc "github.com/mosuka/blast/master/client/grpc"
	"github.com/mosuka/blast/master/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type GetIndexMetaCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
}

var getIndexMetaCmdOpts = GetIndexMetaCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
}

var getIndexMetaCmd = &cobra.Command{
	Use:   "indexmeta",
	Short: "gets the index meta",
	Long:  `The get indexmeta command gets the index meta.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check cluster
		if getIndexMetaCmdOpts.cluster == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), getIndexMetaCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getIndexMetaCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		indexMeta, err := c.GetIndexMeta(ctx, getIndexMetaCmdOpts.cluster)
		if err != nil {
			return err
		}
		resp := struct {
			Cluster   string           `json:"cluster,omitempty"`
			IndexMeta *index.IndexMeta `json:"index_meta,omitempty"`
			Error     error            `json:"error,omitempty"`
		}{
			Cluster:   getIndexMetaCmdOpts.cluster,
			IndexMeta: indexMeta,
			Error:     err,
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
	getIndexMetaCmd.Flags().SortFlags = false

	getIndexMetaCmd.Flags().StringVar(&getIndexMetaCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	getIndexMetaCmd.Flags().IntVar(&getIndexMetaCmdOpts.dialTimeout, "dial-timeout", getIndexMetaCmdOpts.dialTimeout, "dial timeout")
	getIndexMetaCmd.Flags().IntVar(&getIndexMetaCmdOpts.requestTimeout, "request-timeout", getIndexMetaCmdOpts.requestTimeout, "request timeout")
	getIndexMetaCmd.Flags().StringVar(&getIndexMetaCmdOpts.cluster, "cluster", getIndexMetaCmdOpts.cluster, "cluster")

	getCmd.AddCommand(getIndexMetaCmd)
}
