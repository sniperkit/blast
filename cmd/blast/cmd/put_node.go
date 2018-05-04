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
	blastgrpc "github.com/mosuka/blast/master/client/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type PutNodeCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
	node              string
}

var putNodeCmdOpts = PutNodeCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
	node:              "",
}

var putNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "puts the node",
	Long:  `The put node command puts the node to the cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check cluster
		if putNodeCmdOpts.cluster == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
		}

		// check node
		if putNodeCmdOpts.node == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("node").Name)
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), putNodeCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(putNodeCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		// put node
		err = c.PutNode(ctx, putNodeCmdOpts.cluster, putNodeCmdOpts.node)
		if err != nil {
			return err
		}
		resp := struct {
			Cluster string `json:"cluster,omitempty"`
			Node    string `json:"node,omitempty"`
			Error   error  `json:"error,omitempty"`
		}{
			Cluster: putNodeCmdOpts.cluster,
			Node:    putNodeCmdOpts.node,
			Error:   err,
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
	putNodeCmd.Flags().SortFlags = false

	putNodeCmd.Flags().StringVar(&putNodeCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	putNodeCmd.Flags().IntVar(&putNodeCmdOpts.dialTimeout, "dial-timeout", putNodeCmdOpts.dialTimeout, "dial timeout")
	putNodeCmd.Flags().IntVar(&putNodeCmdOpts.requestTimeout, "request-timeout", putNodeCmdOpts.requestTimeout, "request timeout")
	putNodeCmd.Flags().StringVar(&putNodeCmdOpts.cluster, "cluster", putNodeCmdOpts.cluster, "cluster")
	putNodeCmd.Flags().StringVar(&putNodeCmdOpts.node, "node", putNodeCmdOpts.node, "node")

	putCmd.AddCommand(putNodeCmd)
}
