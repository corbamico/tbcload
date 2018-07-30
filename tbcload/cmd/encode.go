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
	"encoding/ascii85"
	"fmt"

	"github.com/corbamico/tbcload"
	"github.com/spf13/cobra"
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode [string to encode]",
	Short: "encode a string into ascii85(re-map), which tbc file used",
	Long: `tbc file use ascii85 encode and map special ascii code.
For example:
proc->,CHr@`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			dst := make([]byte, ascii85.MaxEncodedLen(len(args[0])))
			if ndst := tbcload.Encode(dst, []byte(args[0])); ndst > 0 {
				fmt.Printf("source:%s\n", args[0])
				fmt.Printf("encode:%s\n", dst[:ndst])
			} else {
				fmt.Printf("source:%s\n", args[0])
				fmt.Printf("encode error.")
			}
		} else {
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
