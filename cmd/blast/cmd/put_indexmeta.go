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
	"github.com/mosuka/blast/config"
	"github.com/mosuka/blast/index"
	blastgrpc "github.com/mosuka/blast/master/client/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"time"
)

type PutIndexMetaCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
	indexMetaPath     string
}

var putIndexMetaCmdOpts = PutIndexMetaCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
	indexMetaPath:     "",
}

var putIndexMetaCmd = &cobra.Command{
	Use:   "indexmeta",
	Short: "puts the index meta",
	Long:  `The put indexmeta command puts the index meta to the cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check cluster
		if putIndexMetaCmdOpts.cluster == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
		}

		// check indexMetaPath
		if putIndexMetaCmdOpts.indexMetaPath == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("indexmeta").Name)
		}

		// create index meta
		file, err := os.Open(putIndexMetaCmdOpts.indexMetaPath)
		if err != nil {
			return err
		}
		defer file.Close()

		indexMeta, err := index.LoadIndexMeta(file)
		if err != nil {
			return err
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), putIndexMetaCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(putIndexMetaCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		// put node
		err = c.PutIndexMeta(ctx, putIndexMetaCmdOpts.cluster, indexMeta)
		if err != nil {
			return err
		}
		resp := struct {
			Cluster   string           `json:"cluster,omitempty"`
			IndexMeta *index.IndexMeta `json:"index_meta,omitempty"`
			Error     error            `json:"error,omitempty"`
		}{
			Cluster:   putIndexMetaCmdOpts.cluster,
			IndexMeta: indexMeta,
			Error:     err,
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
	putIndexMetaCmd.Flags().SortFlags = false

	putIndexMetaCmd.Flags().StringVar(&putIndexMetaCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	putIndexMetaCmd.Flags().IntVar(&putIndexMetaCmdOpts.dialTimeout, "dial-timeout", putIndexMetaCmdOpts.dialTimeout, "dial timeout")
	putIndexMetaCmd.Flags().IntVar(&putIndexMetaCmdOpts.requestTimeout, "request-timeout", putIndexMetaCmdOpts.requestTimeout, "request timeout")
	putIndexMetaCmd.Flags().StringVar(&putIndexMetaCmdOpts.cluster, "cluster", putIndexMetaCmdOpts.cluster, "cluster")
	putIndexMetaCmd.Flags().StringVar(&putIndexMetaCmdOpts.indexMetaPath, "index-meta-path", putIndexMetaCmdOpts.indexMetaPath, "index meta path")

	putCmd.AddCommand(putIndexMetaCmd)
}
