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
	"github.com/blevesearch/bleve"
	"github.com/buger/jsonparser"
	"github.com/mosuka/blast/index/client"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

type SearchCommandOptions struct {
	grpcServerAddress string
	dialTimeout       int
	requestTimeout    int
	request           string
	query             string
	size              int
	from              int
	explain           bool
	fields            []string
	sorts             []string
	facets            string
	highlight         string
	highlightStyle    string
	highlightFields   []string
	includeLocations  bool
}

var searchCmdOpts = SearchCommandOptions{
	grpcServerAddress: "localhost:5000",
	dialTimeout:       15000,
	requestTimeout:    15000,
	request:           "",
	query:             "",
	size:              10,
	from:              0,
	explain:           false,
	fields:            []string{},
	sorts:             []string{},
	facets:            "",
	highlight:         "",
	highlightStyle:    "",
	highlightFields:   []string{},
	includeLocations:  false,
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "searches the documents",
	Long:  `The search command searches the documents.`,
	RunE:  runESearchCmd,
}

func runESearchCmd(cmd *cobra.Command, args []string) error {
	// read request
	data := []byte(searchCmdOpts.request)

	// get search_request from request
	var searchRequest *bleve.SearchRequest
	searchRequestBytes, _, _, err := jsonparser.Get(data, "search_request")
	if err == nil {
		err = json.Unmarshal(searchRequestBytes, &searchRequest)
		if err != nil {
			return err
		}
	} else {
		searchRequest = bleve.NewSearchRequest(nil)
	}

	// overwrite request by command line option
	if cmd.Flag("query").Changed {
		searchRequest.Query = bleve.NewQueryStringQuery(searchCmdOpts.query)
	}
	if cmd.Flag("size").Changed {
		searchRequest.Size = searchCmdOpts.size
	}
	if cmd.Flag("from").Changed {
		searchRequest.From = searchCmdOpts.from
	}
	if cmd.Flag("explain").Changed {
		searchRequest.Explain = searchCmdOpts.explain
	}
	if cmd.Flag("field").Changed {
		searchRequest.Fields = searchCmdOpts.fields
	}
	if cmd.Flag("sort").Changed {
		searchRequest.SortBy(searchCmdOpts.sorts)
	}
	if cmd.Flag("facets").Changed {
		facetRequest := bleve.FacetsRequest{}
		err := json.Unmarshal([]byte(searchCmdOpts.facets), &facetRequest)
		if err != nil {
			return err
		}
		searchRequest.Facets = facetRequest
	}
	if cmd.Flag("highlight").Changed {
		highlightRequest := bleve.NewHighlight()
		err := json.Unmarshal([]byte(searchCmdOpts.highlight), highlightRequest)
		if err != nil {
			return err
		}
		searchRequest.Highlight = highlightRequest
	}
	if cmd.Flag("highlight-style").Changed || cmd.Flag("highlight-field").Changed {
		highlightRequest := bleve.NewHighlightWithStyle(searchCmdOpts.highlightStyle)
		highlightRequest.Fields = searchCmdOpts.highlightFields
		searchRequest.Highlight = highlightRequest
	}
	if cmd.Flag("include-locations").Changed {
		searchRequest.IncludeLocations = searchCmdOpts.includeLocations
	}

	// create client
	c, err := client.NewGRPCClient(context.Background(), searchCmdOpts.grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer c.Close()

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(searchCmdOpts.requestTimeout)*time.Millisecond)
	defer cancel()

	// search documents from index
	searchResult, err := c.Search(ctx, searchRequest)
	resp := struct {
		SearchResult *bleve.SearchResult `json:"search_result,omitempty"`
		Error        error               `json:"error,omitempty"`
	}{
		SearchResult: searchResult,
		Error:        err,
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
	searchCmd.Flags().SortFlags = false

	searchCmd.Flags().StringVar(&searchCmdOpts.grpcServerAddress, "grpc-server-address", searchCmdOpts.grpcServerAddress, "Blast server to connect to using gRPC")
	searchCmd.Flags().IntVar(&searchCmdOpts.dialTimeout, "dial-timeout", searchCmdOpts.dialTimeout, "dial timeout")
	searchCmd.Flags().IntVar(&searchCmdOpts.requestTimeout, "request-timeout", searchCmdOpts.requestTimeout, "request timeout")
	searchCmd.Flags().StringVar(&searchCmdOpts.request, "request", searchCmdOpts.request, "resource file")
	searchCmd.Flags().StringVar(&searchCmdOpts.query, "query", searchCmdOpts.query, "query string")
	searchCmd.Flags().IntVar(&searchCmdOpts.size, "size", searchCmdOpts.size, "number of hits to return")
	searchCmd.Flags().IntVar(&searchCmdOpts.from, "from", searchCmdOpts.from, "starting from index of the hits to return")
	searchCmd.Flags().BoolVar(&searchCmdOpts.explain, "explain", searchCmdOpts.explain, "contain an explanation of how scoring of the hits was computed")
	searchCmd.Flags().StringSliceVar(&searchCmdOpts.fields, "field", searchCmdOpts.fields, "specify a set of fields to return")
	searchCmd.Flags().StringSliceVar(&searchCmdOpts.sorts, "sort", searchCmdOpts.sorts, "sorting to perform")
	searchCmd.Flags().StringVar(&searchCmdOpts.facets, "facets", searchCmdOpts.facets, "faceting to perform")
	searchCmd.Flags().StringVar(&searchCmdOpts.highlight, "highlight", searchCmdOpts.highlight, "highlighting to perform")
	searchCmd.Flags().StringVar(&searchCmdOpts.highlightStyle, "highlight-style", searchCmdOpts.highlightStyle, "highlighting style")
	searchCmd.Flags().StringSliceVar(&searchCmdOpts.highlightFields, "highlight-field", searchCmdOpts.highlightFields, "specify a set of fields to highlight")
	searchCmd.Flags().BoolVar(&searchCmdOpts.includeLocations, "include-locations", searchCmdOpts.includeLocations, "include terms locations")

	RootCmd.AddCommand(searchCmd)
}
