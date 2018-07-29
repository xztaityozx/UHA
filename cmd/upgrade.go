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
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "UHAを更新します",
	Long:  `UHAの更新をチェックし、アップグレードします`,
	Run: func(cmd *cobra.Command, args []string) {
		branch, _ := cmd.PersistentFlags().GetString("branch")
		if err := upgrade(branch); err != nil {
			log.Fatal(err)
		}
	},
}

func upgrade(branch string) error {
	//
	os.Chdir(SelfPath)

	spin := spinner.New(spinner.CharSets[35], 500*time.Millisecond)
	spin.Prefix = "UHA\t"
	spin.Suffix = " Upgrading... "
	spin.FinalMSG = ""
	spin.Start()
	defer spin.Stop()

	checkout := exec.Command("git", "checkout", "master")
	if b, err := checkout.CombinedOutput(); err != nil {
		log.Fatal(string(b))
	}

	git := exec.Command("git", "pull")
	if b, err := git.CombinedOutput(); err != nil {
		log.Fatal(string(b))
	}

	target := exec.Command("git", "checkout", branch)
	if b, err := target.CombinedOutput(); err != nil {
		log.Fatal(string(b))
	}

	get := exec.Command("go", "get")
	if b, err := get.CombinedOutput(); err != nil {
		log.Fatal(string(b))
	}

	install := exec.Command("go", "install")
	if b, err := install.CombinedOutput(); err != nil {
		log.Fatal(string(b))
	}

	log.Println("\n\033[1;32mUpgraded UHA\033[0;39m")
	return nil
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.PersistentFlags().StringP("branch", "b", "master", "Pullしてくるbranchを指定します")
}
