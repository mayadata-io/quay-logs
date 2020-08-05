// /*
// Copyright 2020 The MayaData Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

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

// // ListableConfig is used to initialise Listable instance
type ListableConfig struct {
	Namespace          string
	AuthToken          string
	BaseOutputFilePath string
	IsWriteToFile      bool
	Debug              bool
}

// // Listable is used to list all images of the given namespace
type Listable struct {
	*Popularity
}

// // NewLister returns a new instance of Listable
func NewLister(config ListableConfig) (*Listable, error) {
	fmt.Println(config, "config")
	if(config.IsWriteToFile) {
		folder := path.Join(
			config.BaseOutputFilePath,
			config.Namespace,
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
				"Failed to create repo list folder: Name %q",
				folder,
			)
		}
	}
	
	return &Listable{
		Popularity: &Popularity{
			Namespace:          config.Namespace,
			AuthToken:          config.AuthToken,
			BaseOutputFilePath: config.BaseOutputFilePath,
			IsWriteToFile:      config.IsWriteToFile,
			Debug:              config.Debug,
		},
	}, nil
}

// ListReposAndWriteToFileOptionally invokes the API to list images
// belonging to a namespace and then write them to a file
func (l *Listable) ListReposAndWriteToFileOptionally() (PopularList, error) {
	return l.Popularity.ListReposByPopularityAndWriteToFileOptionally()
}

// // PopularityOption is typed function to mutate Popularity instance
type PopularityOption func(*Popularity) error

// // Popularity helps fetch images ranked with their popularity
type Popularity struct {
	Namespace          string
	AuthToken          string
	BaseOutputFilePath string
	IsWriteToFile      bool
	Debug              bool

	currentFilename string
}

// ListReposByPopularityAndWriteToFileOptionally requests for repos by
// invoking API and subsequently writes them to files.
func (p *Popularity) ListReposByPopularityAndWriteToFileOptionally() (PopularList, error) {
	var out = &PopularList{}

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
		p.currentFilename = path.Join(p.BaseOutputFilePath, p.Namespace, filename)

		// Invoke API to request for logs
		//
		// NOTE:
		//	This will run through a set of post functions if set,
		// after executing this API
		got, err := p.RequestReposForPageToken(pagetoken)
		if err != nil {
			return PopularList{}, err
		}
		out.Items = append(out.Items, got.Items...)

		// TODO
		// 	Uncomment following code once API response supports
		// next token field
		//
		// prepare for next iteration
		//isNextpage = got.NextPage != ""
		//pagetoken = got.NextPage

		// TODO: next_page is not defined in the response to following API:
		// - https://quay.io/api/v1/repository
		//
		// Hence next page is set to false after a single invocation
		isNextpage = false
		index++
	}
	return *out, nil
}

// RequestReposForPageToken lists the repos belonging to a namespace
func (p *Popularity) RequestReposForPageToken(pagetoken string) (PopularList, error) {
	req := &HTTPRequest{
		AuthToken: p.AuthToken,
		URL:       "https://quay.io/api/v1/repository",
		Method:    GET,
		QueryParams: map[string]string{
			"popularity": "true",
			"namespace":  p.Namespace,
			"next_page":  pagetoken,
		},
	}
	resp, err := req.Invoke()
	if err != nil {
		return PopularList{}, errors.Wrapf(
			err,
			"Failed to invoke popularity request: %q",
			req.URL,
		)
	}
	if p.IsWriteToFile {
		p.WriteToFile(resp.Body(), p.currentFilename)
	}
	var out PopularList
	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return PopularList{}, errors.Wrapf(
			err,
			"Failed to unmarshal to PopularityList",
		)
	}
	if p.Debug {
		log.Printf("Successfully invoked popularity request")
	}
	return out, nil
}

// WriteToFile creates a file with images having
// popularity ratings. This file is named with today's date.
func (p *Popularity) WriteToFile(raw []byte, filename string) error {
	err := ioutil.WriteFile(
		filename,
		raw,
		0644,
	)
	if err != nil {
		return errors.Wrapf(
			err,
			"Failed to write popularity list to %q",
			filename,
		)
	}
	log.Printf("Wrote popularity list to %q", filename)
	return nil
}