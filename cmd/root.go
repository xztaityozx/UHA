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
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/c-bata/go-prompt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var config Config

var DEF_MOTES = []string{"100", "200", "500", "1000", "2000", "5000", "10000", "20000", "50000"}

var SelfPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "UHA",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", path.Join(os.Getenv("HOME"), ".config", "UHA", ".UHA.json"), "config file ")

	viper.SetDefault("Simulation", Simulation{
		Monte:  DEF_MOTES,
		Range:  Range{Start: "2.5ns", Step: "7.5ns", Stop: "17.5ns"},
		SimDir: "/tmp",
		DstDir: "/tmp",
		Signal: "N2",
		Vtn:    Node{Voltage: 0.6, Sigma: 0.0, Deviation: 1.0},
		Vtp:    Node{Voltage: 0.6, Sigma: 0.0, Deviation: 1.0},
	})
	viper.SetDefault("Repositorys", []Repository{})
	viper.SetDefault("TaskDir", path.Join(os.Getenv("HOME"), ".config", "UHA", "task"))
	viper.SetDefault("SpreadSheet", SpreadSheet{})
	viper.SetDefault("Simulation.SEED", 1)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".UHA" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".UHA")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	if !IsAccurateConfig(cfgFile) {
		log.Println("\033[1;31mConfigが変っぽいですが大丈夫ですか？(y/n)\033[0;39m")
		if y := prompt.Input(">>> ", completer); y != "y" && y != "yes" {
			log.Fatal(`https://github.com/xztaityozx/UHA/blob/master/doc/setting.md を見てみてください`)
		}
	}

	//dir
	TaskDir, _ := homedir.Expand(config.TaskDir)
	ReserveRunDir = filepath.Join(TaskDir, RUN, RESERVE)
	DoneRunDir = filepath.Join(TaskDir, RUN, DONE)
	FailedRunDir = filepath.Join(TaskDir, RUN, FAILED)
	ReserveSRunDir = filepath.Join(TaskDir, SRUN, RESERVE)
	DoneSRunDir = filepath.Join(TaskDir, SRUN, DONE)
	FailedSRunDir = filepath.Join(TaskDir, SRUN, FAILED)

	home, _ := homedir.Dir()
	ConfigDir = filepath.Join(home, ".config", "UHA")
	SelfPath = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "xztaityozx", "UHA")
	NextPath = filepath.Join(ConfigDir, "next.json")
	tryMkdir(ReserveRunDir)
	tryMkdir(DoneRunDir)
	tryMkdir(FailedRunDir)
	tryMkdir(ReserveSRunDir)
	tryMkdir(DoneSRunDir)
	tryMkdir(FailedSRunDir)

}

func WriteConfig() error {
	return viper.WriteConfig()
}
