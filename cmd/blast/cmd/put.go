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
	"github.com/spf13/cobra"
)

type PutCmdOpts struct {
	outputFormat string
}

var putCmdOpts = PutCmdOpts{
	outputFormat: "json",
}

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "puts the object",
	Long:  `The put command puts the object.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}

		_, _, err := cmd.Find(args)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	putCmd.PersistentFlags().StringVar(&putCmdOpts.outputFormat, "output-format", putCmdOpts.outputFormat, "output format")

	RootCmd.AddCommand(putCmd)
}
