// Copyright © 2018 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version UHAVersion = UHAVersion{
	Major:    1,
	Minor:    3,
	Build:    47,
	Revision: 7,
	Status:   "Beta",
	Date:     "2018/07/31",
}

type UHAVersion struct {
	Major    int
	Minor    int
	Build    int
	Revision int
	Status   string
	Date     string
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョンを出力します",
	Long:  `UHA のバージョンを出力して終わります`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getVersion(Version))
	},
}

func getVersion(v UHAVersion) string {
	return fmt.Sprintf("UHA [Ultra H_SPICE Attacker]\nVersion: %d.%d.%d.%d %s (%s)\nRepository: https://github.com/xztaityozx/UHA\nLicense: MIT", v.Major, v.Minor, v.Build, v.Revision, v.Status, v.Date)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
