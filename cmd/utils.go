package cmd

import "log"

// Struct to read and write data from yaml file
type Bookmark struct {
	URL      string
	Title    string
	Category string
}

func handleErr(err error, message string) {
	if err == nil {
		return
	}

	log.Fatalf("%s: %v", message, err)
}
