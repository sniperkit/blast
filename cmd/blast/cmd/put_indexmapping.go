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
	"github.com/mosuka/blast/index"
	blastgrpc "github.com/mosuka/blast/master/client/grpc"
	"github.com/mosuka/blast/master/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"time"
)

type PutIndexMappingCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
	indexMappingPath  string
}

var putIndexMappingCmdOpts = PutIndexMappingCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
	indexMappingPath:  "",
}

var putIndexMappingCmd = &cobra.Command{
	Use:   "indexmapping",
	Short: "puts the index mapping",
	Long:  `The put indexmapping command puts the index mapping to the cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check cluster
		if putIndexMappingCmdOpts.cluster == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
		}

		// check indexMappingPath
		if putIndexMappingCmdOpts.indexMappingPath == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("node").Name)
		}

		// create index mapping
		file, err := os.Open(putIndexMappingCmdOpts.indexMappingPath)
		if err != nil {
			return err
		}
		defer file.Close()

		indexMapping, err := index.LoadIndexMapping(file)
		if err != nil {
			return err
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), putIndexMappingCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(putIndexMappingCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		// put node
		err = c.PutIndexMapping(ctx, putIndexMappingCmdOpts.cluster, indexMapping)
		if err != nil {
			return err
		}
		resp := struct {
			Cluster      string                    `json:"cluster,omitempty"`
			IndexMapping *mapping.IndexMappingImpl `json:"index_mapping,omitempty"`
			Error        error                     `json:"error,omitempty"`
		}{
			Cluster:      putIndexMappingCmdOpts.cluster,
			IndexMapping: indexMapping,
			Error:        err,
		}

		// output response
		switch putCmdOpts.outputFormat {
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
	putIndexMappingCmd.Flags().SortFlags = false

	putIndexMappingCmd.Flags().StringVar(&putIndexMappingCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	putIndexMappingCmd.Flags().IntVar(&putIndexMappingCmdOpts.dialTimeout, "dial-timeout", putIndexMappingCmdOpts.dialTimeout, "dial timeout")
	putIndexMappingCmd.Flags().IntVar(&putIndexMappingCmdOpts.requestTimeout, "request-timeout", putIndexMappingCmdOpts.requestTimeout, "request timeout")
	putIndexMappingCmd.Flags().StringVar(&putIndexMappingCmdOpts.cluster, "cluster", putIndexMappingCmdOpts.cluster, "cluster")
	putIndexMappingCmd.Flags().StringVar(&putIndexMappingCmdOpts.indexMappingPath, "index-mapping-path", putIndexMappingCmdOpts.indexMappingPath, "index mapping path")

	putCmd.AddCommand(putIndexMappingCmd)
}
