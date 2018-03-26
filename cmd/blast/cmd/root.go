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
	"fmt"
	"github.com/mosuka/blast/version"
	"github.com/spf13/cobra"
	"os"
)

type RootCommandOptions struct {
	versionFlag bool
}

var rootCmdOpts = RootCommandOptions{
	versionFlag: false,
}

var RootCmd = &cobra.Command{
	Use:   "blast",
	Short: "Blast",
	Long:  `The Command Line Interface for the Blast.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if rootCmdOpts.versionFlag {
			fmt.Printf("%s\n", version.Version)
			os.Exit(0)
		}
		return nil
	},
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
	RootCmd.Flags().SortFlags = false
	RootCmd.PersistentFlags().BoolVarP(&rootCmdOpts.versionFlag, "version", "v", rootCmdOpts.versionFlag, "show version number")
}
