/*
Copyright 2020 The MayaData Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	gmetrics "github.com/mayadata.io/quay-logs"
)

// The possible list of command line arguments/flags which can be
// given to the `go run cmd/main.go` command are listed inside
// the var( ) block.
var (
	debug = flag.Bool(
		"debug",
		false,
		"when set to true will result in more verbose output",
	)
	repository = flag.String(
		"repo",
		"",
		"quay.io repository",
	)
	quayAuthToken = flag.String(
		"quay-auth-token",
		os.Getenv("QUAY-AUTH-TOKEN"),
		"authentication token to communicate with quay.io APIs",
	)
	quayNamespace = flag.String(
		"quay-namespace",
		os.Getenv("QUAY-NAMESPACE"),
		"namespace to be used while querying quay.io APIs",
	)
	logsFilePath = flag.String(
		"logs-file-path",
		"./logs",
		"(optional) absolute path to the quay repo's log files",
	)
)

// makes all the required directories/folders from the arguments
// by formatting them with -p flags.
// - Here we create 1 directory/folder i.e ./logs
//
// ./logs is used for storing the repos in order of popularity
func mkdirAll() {
	var cmds = map[string][]string{
		"logs": {"-p", *logsFilePath},
	}
	for _, commandargs := range cmds {
		cmd := exec.Command("mkdir", commandargs...)
		err := cmd.Run()
		if err != nil {
			log.Fatalf(
				"Failed to create folder: %s : %v",
				commandargs,
				err,
			)
		}
	}
}

// The main function has the following logic
// - It makes the required directories.
// - It lists all the repos in the sorted order of popularity in the
//   namespace into `repolist`.
// - It iterates through each of the repos and download its Logs and
//   stores them in different files.
func main() {
	// parses the flags. It must be called before using any of the flags.
	flag.Parse()

	if *quayAuthToken == "" {
		log.Fatal("Missing quay auth token")
	}
	if *quayNamespace == "" {
		log.Fatal("Missing quay namespace")
	}

	// create folders that will host various files downloaded from quay
	mkdirAll()

	// list repos
	log.Print("Will list all repos")

	// We create a `NewLister` (refer `list.go`) and set
	//`IsWriteToFile` false because we don't want to store the data
	// in the files.
	l, err := gmetrics.NewLister(gmetrics.ListableConfig{
		AuthToken:          *quayAuthToken,
		BaseOutputFilePath: "",
		Namespace:          *quayNamespace,
		IsWriteToFile:      false,
		Debug:              *debug,
	})
	//checking for errors
	if err != nil {
		log.Fatalf(
			"Failed to initialise lister: %v",
			err,
		)
	}

	// repolist contains repos in order of popularity
	// We call the `ListReposAndWriteToFileOptionally( )` function
	// to get all the repolist in the namespace as a JSON format.
	// It returns all the repos in sorted order of popularity.
	repolist, err := l.ListReposAndWriteToFileOptionally()
	if err != nil {
		log.Fatalf(
			"Failed to list repos: %v",
			err,
		)
	}

	// download logs of all repos
	log.Print("Will download logs of all repos")
	for _, repo := range repolist.Items {
		logger, err := gmetrics.NewLogger(gmetrics.LoggableConfig{
			AuthToken:          *quayAuthToken,
			Namespace:          *quayNamespace,
			Name:               repo.Name,
			IsWriteToFile:      true,
			BaseOutputFilePath: *logsFilePath,
			Debug:              *debug,
		})
		if err != nil {
			log.Fatalf(
				"Failed to initialise logger: %v",
				err,
			)
		}
		_, err = logger.Log()
		if err != nil {
			log.Fatalf(
				"Failed to download logs: %v",
				err,
			)
		}
	}
}
