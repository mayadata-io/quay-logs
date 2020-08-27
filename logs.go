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

package growthmetrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"time"

	"github.com/pkg/errors"
)

// LoggableConfig is used to initialise a Loggable instance
type LoggableConfig struct {
	Namespace          string
	Name               string
	AuthToken          string
	BaseOutputFilePath string
	IsWriteToFile      bool
	Debug              bool
}

// LoggableOption is a typed function to mutate Loggable instance
type LoggableOption func(*Loggable) error

// Loggable fetches image logs by invoking quay.io APIs
type Loggable struct {
	Namespace          string
	Name               string
	AuthToken          string
	BaseOutputFilePath string
	IsWriteToFile      bool
	Debug              bool
	//the value for next two will be assigned in Log()
	currentLogs     []byte
	currentFilename string
}

// NewLogger returns a new instance of Loggable
// It creates a new folder by the mkdir command using the arguments
// passed to it for each of the repos.
func NewLogger(config LoggableConfig) (*Loggable, error) {
	folder := path.Join(
		config.BaseOutputFilePath,
		config.Namespace,
		config.Name,
		//example: ../logs/namespace/reponame/
	)
	cmd := exec.Command(
		"mkdir",
		"-p",
		folder,
	)
	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"Failed to create repo logs folder: Name %q",
			folder,
		)
	}
	return &Loggable{
		AuthToken:          config.AuthToken,
		Namespace:          config.Namespace,
		Name:               config.Name,
		BaseOutputFilePath: config.BaseOutputFilePath,
		IsWriteToFile:      config.IsWriteToFile,
		Debug:              config.Debug,
	}, nil
}

// Log requests for logs by invoking API and subsequently
// writes them to files.
// 
//  It calls `RequestLogsForPageToken( )` to get the logs from 
// the Quay API. It stores them in separate files by calling 
// `WriteToFile` internally. 
// --Here next page is available since the API returns 20 `logs`
// at once. So each files can contain at max 20 `logs`.
func (l *Loggable) Log() (LogList, error) {
	var out = &LogList{}

	var isNextpage = true
	var pagetoken = ""
	var index int

	// File names for all downloads need to have same prefix
	// Variable 'now' defines this prefix
	var now = time.Now().Format("Jan-02-2006-15:04:05")
	for isNextpage {
		// Set or reset filename
		//
		// NOTE:
		// 	Set the full filename by suffixing with index
		//
		// NOTE:
		//	Logs is a list API call that is paged. Each page can
		// optionally be saved to a new file.
		filename := fmt.Sprintf("%s-%d.json", now, index)
		l.currentFilename = path.Join(
			l.BaseOutputFilePath,
			l.Namespace,
			l.Name,
			filename,
		)

		// Invoke API to request for logs
		//
		// NOTE:
		//	This will run through a set of post functions if set,
		// after executing this API
		got, err := l.RequestLogsForPageToken(pagetoken)
		if err != nil {
			return LogList{}, err
		}
		out.Items = append(out.Items, got.Items...)

		// prepare for next iteration
		isNextpage = got.NextPage != ""
		pagetoken = got.NextPage
		index++
	}
	return *out, nil
}

// RequestLogsForPageToken lists the logs of the images belonging 
// to a namespace. It Creates a HTTPRequest with some query parameters
// and invokes it. 
// -- Since `IsWriteToFile` is true here so it calls `WriteToFile` 
// and the JSON is unmarshaled and returned. 
func (l *Loggable) RequestLogsForPageToken(pagetoken string) (LogList, error) {
	if l.Debug {
		log.Printf(
			"Will request logs: Namespace %q: Name %q: Page Token %q",
			l.Namespace,
			l.Name,
			pagetoken,
		)
	}
	req := &HTTPRequest{
		AuthToken: l.AuthToken,
		URL:       "https://quay.io/api/v1/repository/{namespace}/{name}/logs",
		Method:    GET,
		QueryParams: map[string]string{
			"next_page": pagetoken,
		},
		PathParams: map[string]string{
			"namespace": l.Namespace,
			"name":      l.Name,
		},
	}
	resp, err := req.Invoke()
	if err != nil {
		return LogList{}, errors.Wrapf(
			err,
			"Failed to request logs: Namespace %q: Name %q",
			l.Namespace,
			l.Name,
		)
	}
	if resp.StatusCode() != 200 {
		log.Printf(
			"Logs response: Namespace %q: Name %q: StatusCode %d: Error %q",
			l.Namespace,
			l.Name,
			resp.StatusCode(),
			resp.Header().Get("error"),
		)
		return LogList{}, nil
	}
	if l.IsWriteToFile {
		err = l.WriteToFile(resp.Body(), l.currentFilename)
		if err != nil {
			return LogList{}, errors.Wrapf(
				err,
				"Failed to write logs to %q",
				l.currentFilename,
			)
		}
		log.Printf("Wrote logs to %q", l.currentFilename)
	}
	var out LogList
	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return LogList{}, errors.Wrapf(
			err,
			"Failed to unmarshal logs to LogList",
		)
	}
	return out, nil
}

// WriteToFile creates a file with images having popularity ratings.
// This file is named with today's date.
// It writes the content of response body into passed filename with 
// file mode 0644. It stores the logs into 
// `./logs/namespace/reponame/filename.json`
func (l *Loggable) WriteToFile(raw []byte, filename string) error {
	err := ioutil.WriteFile(filename, raw, 0644)
	return errors.Wrapf(
		err,
		"Failed to write logs to %q",
		filename,
	)
}
