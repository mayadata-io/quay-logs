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

func main() {
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
	l, err := gmetrics.NewLister(gmetrics.ListableConfig{
		AuthToken:          *quayAuthToken,
		BaseOutputFilePath: "",
		Namespace:          *quayNamespace,
		IsWriteToFile:      false,
		Debug:              *debug,
	})

	if err != nil {
		log.Fatalf(
			"Failed to initialise lister: %v",
			err,
		)
	}
	//repolist contains repos in order of popularity
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
