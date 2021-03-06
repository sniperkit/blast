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
	"github.com/buger/jsonparser"
	"github.com/mosuka/blast/config"
	blastgrpc "github.com/mosuka/blast/node/client/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type PutDocumentCmdOpts struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	id                string
	fields            string
	request           string
}

var putDocumentCmdOpts = PutDocumentCmdOpts{
	grpcServerAddress: config.DefaultGRPCListenAddress,
	dialTimeout:       5000,
	requestTimeout:    5000,
	id:                "",
	fields:            "",
	request:           "",
}

var putDocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "puts the document",
	Long:  `The index document command puts the document.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read request
		data := []byte(putDocumentCmdOpts.request)

		// get id from request
		id, err := jsonparser.GetString(data, "id")
		if err != nil {
			return err
		}

		// get fields request
		fieldsBytes, _, _, err := jsonparser.Get(data, "fields")
		if err != nil {
			return err
		}
		var fields map[string]interface{}
		err = json.Unmarshal(fieldsBytes, &fields)
		if err != nil {
			return err
		}

		// overwrite id by command line option
		if cmd.Flag("id").Changed {
			id = putDocumentCmdOpts.id
		}

		// overwrite fields by command line option
		if cmd.Flag("fields").Changed {
			err = json.Unmarshal([]byte(putDocumentCmdOpts.fields), &fields)
			if err != nil {
				return err
			}
		}

		// create client
		c, err := blastgrpc.NewGRPCClient(context.Background(), putDocumentCmdOpts.grpcServerAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer c.Close()

		// create context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(putDocumentCmdOpts.requestTimeout)*time.Millisecond)
		defer cancel()

		// put document to index
		err = c.PutDocument(ctx, id, fields)
		resp := struct {
			Id     string                 `json:"id,omitempty"`
			Fields map[string]interface{} `json:"fields,omitempty"`
			Error  error                  `json:"error,omitempty"`
		}{
			Id:     id,
			Fields: fields,
			Error:  err,
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
	putDocumentCmd.Flags().SortFlags = false

	putDocumentCmd.Flags().StringVar(&putDocumentCmdOpts.grpcServerAddress, "grpc-server-address", config.DefaultGRPCListenAddress, "Blast server to connect to using gRPC")
	putDocumentCmd.Flags().IntVar(&putDocumentCmdOpts.dialTimeout, "dial-timeout", putDocumentCmdOpts.dialTimeout, "dial timeout")
	putDocumentCmd.Flags().IntVar(&putDocumentCmdOpts.requestTimeout, "request-timeout", putDocumentCmdOpts.requestTimeout, "request timeout")
	putDocumentCmd.Flags().StringVar(&putDocumentCmdOpts.id, "id", putDocumentCmdOpts.id, "document id")
	putDocumentCmd.Flags().StringVar(&putDocumentCmdOpts.fields, "fields", putDocumentCmdOpts.fields, "document fields")
	putDocumentCmd.Flags().StringVar(&putDocumentCmdOpts.request, "request", putDocumentCmdOpts.request, "request file")

	putCmd.AddCommand(putDocumentCmd)
}
