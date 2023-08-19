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
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Bookmark struct {
	URL      string
	Title    string
	Category string
}

var category string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add new bookmark",
	Long:    `Add new url to your bookmarks, saving it as a json file in your bookmarks directory and then committing it`,
	Example: `stribor add http://github.com/acanoe/stribor`,
	Run: func(cmd *cobra.Command, args []string) {
		directory := homeDirectory + "/." + directoryName

		// Parse url
		url, err := url.Parse(args[0])
		if err != nil {
			log.Fatal("Cannot parse url")
		}

		// build struct
		bookmark := Bookmark{
			URL:      url.String(),
			Title:    "",
			Category: category,
		}

		// convert to yaml
		yamlData, err := yaml.Marshal(&bookmark)
		if err != nil {
			log.Fatal("Cannot convert to yaml")
		}

		// check subfolder
		siteFolder := url.Host
		siteFolderPath := filepath.Join(directory, siteFolder)
		if _, err := os.Stat(siteFolderPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(siteFolderPath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}

		// build file name
		fileName := url.Path
		if fileName == "" {
			fileName = "root"
		}
		fileName = slug.Make(fileName)
		fileName = fileName + ".yaml"

		// write file if not already exists
		filePath := filepath.Join(siteFolderPath, fileName)
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(filePath, yamlData, 0755); err != nil {
				log.Fatal("Cannot write file")
			}
		}

		// check repo and working dir
		//   opens an already existing repository.
		r, err := git.PlainOpen(directory)
		if err != nil {
			log.Fatal("Not a git repo")
		}

		//   get worktree
		w, err := r.Worktree()
		if err != nil {
			log.Fatal("Cannot get work tree")
		}

		//   get status
		status, err := w.Status()
		if err != nil {
			log.Fatal("Cannot get status")
		}

		//   check if file is already committed
		gitFilePath := filepath.Join(siteFolder, fileName)
		if status.IsUntracked(gitFilePath) {
			//   add file
			w.Add(gitFilePath)

			//   commit the change
			commit, err := w.Commit(fmt.Sprintf("added %s", url.String()), &git.CommitOptions{
				Author: &object.Signature{
					Name:  "Stribor",
					Email: "stribor@local.store",
					When:  time.Now(),
				},
			})
			if err != nil {
				log.Fatalf("Cannot commit file: %v", err)
			}

			obj, err := r.CommitObject(commit)
			if err != nil {
				log.Fatal("Cannot read HEAD")
			}

			fmt.Println(obj)

		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("category", "c", "other", "Bookmark category")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
