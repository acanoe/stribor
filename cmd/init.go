/*
Copyright © 2020 Hendika N.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

var force bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new bookmark directory",
	Long:  `Initialize a directory to be used as a git repo for the bookmark files`,
	Run: func(cmd *cobra.Command, args []string) {
		directory := filepath.Join(homeDirectory, "."+directoryName)

		if _, err := os.Stat(homeDirectory + "/.bookmarks"); !os.IsNotExist(err) && directoryName != "bookmarks" {
			os.RemoveAll(homeDirectory + "/.bookmarks")
		}

		if force {
			fmt.Println("Using --force, I hope you know what you're doing")
			if err := os.RemoveAll(directory); err != nil {
				log.Fatal(err)
			}
		}

		err := os.Mkdir(directory, 0755)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				log.Fatal("Folder already exist, use --force to overwrite")
			}
			log.Fatal(err)
		}

		repo, err := git.PlainInit(directory, false)
		if err != nil {
			log.Fatal(err)
		}

		config := config.NewConfig()
		config.User.Name = "stribor"
		config.User.Email = "stribor@local.store"

		repo.SetConfig(config)

		fmt.Println("Bookmarks folder initialized, add your first bookmark now!")
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			viper.SetConfigType("yaml")
			err := viper.SafeWriteConfig()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Delete existing folder and start over (make sure you know what you are doing)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
