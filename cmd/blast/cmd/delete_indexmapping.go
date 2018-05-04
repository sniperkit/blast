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

type DeleteIndexMappingCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	cluster           string
}

var deleteIndexMappingCmdOpts = DeleteIndexMappingCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	cluster:           "",
}

var deleteIndexMappingCmd = &cobra.Command{
	Use:   "indexmapping",
	Short: "deletes the index mapping",
	Long:  `The delete indexmapping command deletes the index mapping.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check cluster
		if deleteIndexMappingCmdOpts.cluster == "" {
			return fmt.Errorf("required flag: --%s", cmd.Flag("cluster").Name)
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), deleteIndexMappingCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deleteIndexMappingCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		err = c.DeleteIndexMapping(ctx, deleteIndexMappingCmdOpts.cluster)
		if err != nil {
			return err
		}
		resp := struct {
			Cluster string `json:"cluster,omitempty"`
			Error   error  `json:"error,omitempty"`
		}{
			Cluster: deleteIndexMappingCmdOpts.cluster,
			Error:   err,
		}

		// output response
		switch deleteCmdOpts.outputFormat {
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
	deleteIndexMappingCmd.Flags().SortFlags = false

	deleteIndexMappingCmd.Flags().StringVar(&deleteIndexMappingCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	deleteIndexMappingCmd.Flags().IntVar(&deleteIndexMappingCmdOpts.dialTimeout, "dial-timeout", deleteIndexMappingCmdOpts.dialTimeout, "dial timeout")
	deleteIndexMappingCmd.Flags().IntVar(&deleteIndexMappingCmdOpts.requestTimeout, "request-timeout", deleteIndexMappingCmdOpts.requestTimeout, "request timeout")
	deleteIndexMappingCmd.Flags().StringVar(&deleteIndexMappingCmdOpts.cluster, "cluster", deleteIndexMappingCmdOpts.cluster, "cluster")

	deleteCmd.AddCommand(deleteIndexMappingCmd)
}
