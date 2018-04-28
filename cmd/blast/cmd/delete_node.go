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

type DeleteNodeCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
	node              string
}

var deleteNodeCmdOpts = DeleteNodeCmdOpts{
	grpcServerAddress: config.DefaultMasterGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
	node:              "",
}

var deleteNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "deletes the node",
	Long:  `The delete node command deletes the node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check cluster
		if getNodeCmdOpts.cluster == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
		}

		// check node
		if getNodeCmdOpts.node == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("node").Name)
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), getNodeCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getNodeCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		// get document from index
		err = c.DeleteNode(ctx, getNodeCmdOpts.cluster, getNodeCmdOpts.node)
		if err != nil {
			return err
		}
		resp := struct {
			Cluster string `json:"cluster,omitempty"`
			Node    string `json:"node,omitempty"`
			Error   error  `json:"error,omitempty"`
		}{
			Cluster: getNodeCmdOpts.cluster,
			Node:    getNodeCmdOpts.node,
			Error:   err,
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
	deleteNodeCmd.Flags().SortFlags = false

	deleteNodeCmd.Flags().StringVar(&getNodeCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultMasterGRPCListenAddress, "Blast server to connect to using gRPC")
	deleteNodeCmd.Flags().IntVar(&getNodeCmdOpts.dialTimeout, "dial-timeout", getNodeCmdOpts.dialTimeout, "dial timeout")
	deleteNodeCmd.Flags().IntVar(&getNodeCmdOpts.requestTimeout, "request-timeout", getNodeCmdOpts.requestTimeout, "request timeout")
	deleteNodeCmd.Flags().StringVar(&getNodeCmdOpts.cluster, "cluster", getNodeCmdOpts.cluster, "cluster")
	deleteNodeCmd.Flags().StringVar(&getNodeCmdOpts.node, "node", getNodeCmdOpts.node, "node")

	deleteCmd.AddCommand(deleteNodeCmd)
}
