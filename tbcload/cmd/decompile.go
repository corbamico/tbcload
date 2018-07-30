// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/corbamico/tbcload"
	"github.com/spf13/cobra"
)

// decompileCmd represents the decompile command
var decompileCmd = &cobra.Command{
	Use:   "decompile [file|url]",
	Short: "decompile a .tbc file, which can be on disk/url",
	Long: `decompile a .tbc file, which can be on disk/url.

Example:
    tbcload decompile  test.tbc  #decompile a file named test.tbc
    tbcload decompile  https://github.com/ActiveState/teapot/raw/master/lib/tbcload/tests/tbc10/proc.tbc
			#decompile from a url`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			uri := args[0]

			if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
				parseURL(uri)
			} else {
				parseFile(uri)
			}
		} else {
			cmd.Usage()
		}
	},
}

var detail bool

func init() {
	rootCmd.AddCommand(decompileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disassembleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	decompileCmd.Flags().BoolVarP(&detail, "detail", "d", false, "decompile bytecode instruction too")
}

func parseFile(uri string) {
	r, err := os.Open(uri)
	if err != nil {
		fmt.Printf("failed read from file (%s), error as (%s)\n", uri, err)
		return
	}
	p := tbcload.NewParser(r, os.Stdout)
	p.Detail = detail
	if err = p.Parse(); err != nil {
		fmt.Printf("failed parse file (%s), error as (%s)\n", uri, err)
		return
	}
}
func parseURL(uri string) {
	r, err := http.Get(uri)
	if err != nil {
		fmt.Printf("failed read from uri (%s), error as (%s)\n", uri, err)
		return
	}
	p := tbcload.NewParser(r.Body, os.Stdout)
	p.Detail = detail
	
	if err = p.Parse(); err != nil {
		fmt.Printf("failed parse uri (%s), error as (%s)\n", uri, err)
		return
	}
	r.Body.Close()
}
