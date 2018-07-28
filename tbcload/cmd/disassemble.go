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

// disassembleCmd represents the disassemble command
var disassembleCmd = &cobra.Command{
	Use:   "disassemble [file|url]",
	Short: "disassemble a .tbc file, which can be on disk/url",
	Long: `disassemble a .tbc file, which can be on disk/url.

Example:
     tbcload disassemble  test.tbc  #disassemble a file named test.tbc
     tbcload disassemble  https://github.com/ActiveState/teapot/raw/master/lib/tbcload/tests/tbc10/proc.tbc
                 	 #disassemble from a url`,
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

func init() {
	rootCmd.AddCommand(disassembleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disassembleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disassembleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func usage() {
// 	message := `
// Usage: distclbytecode [-h|--help] [file|url]
// disassemble tcl bytecode file (usally .tbc file).

// Example:
//     distclbytecode  test.tbc  #disassemble a file named test.tbc
//     distclbytecode  https://github.com/ActiveState/teapot/raw/master/lib/tbcload/tests/tbc10/proc.tbc
//                 	#disassemble from a url
// 	`
// 	fmt.Println(message)
// 	os.Exit(1)
// }

func parseFile(uri string) {
	r, err := os.Open(uri)
	if err != nil {
		fmt.Printf("failed read from file (%s), error as (%s)\n", uri, err)
		return
	}
	p := tbcload.NewParser(r, os.Stdout)
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
	if err = p.Parse(); err != nil {
		fmt.Printf("failed parse uri (%s), error as (%s)\n", uri, err)
		return
	}
	r.Body.Close()
}

// func main() {
// 	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
// 		usage()
// 	}
// 	uri := os.Args[1]

// 	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
// 		parseURL(uri)
// 	} else {
// 		parseFile(uri)
// 	}
// }
