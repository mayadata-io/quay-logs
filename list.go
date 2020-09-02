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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// ListableConfig is used to initialise Listable instance
type ListableConfig struct {
	Namespace          string
	AuthToken          string
	BaseOutputFilePath string
	IsWriteToFile      bool
	Debug              bool
	Windows            bool
}

// Listable is used to list all images of the given namespace
type Listable struct {
	*Popularity
}

// NewLister returns a new instance of Listable
// It creates a new folder by the mkdir command using the arguments
// passed to it.
func NewLister(config ListableConfig) (*Listable, error) {
	if config.IsWriteToFile {
		folder := path.Join(
			config.BaseOutputFilePath,
			config.Namespace,
		)
		if config.Windows == true {
			// since windows doesn't support '\'
			p := filepath.FromSlash(folder)
			if config.Debug == true {
				fmt.Println("Creating folders: " + p)
			}
			err := os.MkdirAll(p, 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// this block is for linux. Supports '/' filepath systems
			cmd := exec.Command(
				"mkdir",
				"-p",
				folder,
				//It will be something like "mkdir -p ./popularity/namespace/
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
// This actually calls `ListReposByPopularityAndWriteToFileOptionally( )`function.
func (l *Listable) ListReposAndWriteToFileOptionally() (PopularList, error) {
	return l.Popularity.ListReposByPopularityAndWriteToFileOptionally()
}

// PopularityOption is typed function to mutate Popularity instance
type PopularityOption func(*Popularity) error

// Popularity helps fetch images ranked with their popularity
type Popularity struct {
	Namespace          string
	AuthToken          string
	BaseOutputFilePath string
	IsWriteToFile      bool
	Debug              bool
	Windows            bool
	// the below will be assigned in the next function
	currentFileName string
	fileNamePath    string
}

// ListReposByPopularityAndWriteToFileOptionally requests for repos by
// invoking API and subsequently writes them to files.
//
// It calls the `RequestReposForPageToken( )` which returns all repos
// name in order of popularity.
// -- Right now we don't have 100 repos that's why all the data are in
// one page. Thus some codes are commented below.
func (p *Popularity) ListReposByPopularityAndWriteToFileOptionally() (PopularList, error) {
	var out = &PopularList{}

	var isNextpage = true
	var pagetoken = ""
	var index int

	// File names for all downloads need to have same prefix
	// Variable 'now' defines this prefix
	var now = time.Now().Format("Jan-02-2006-15:04:05")
	if p.Windows == true {
		// since windows doesn't support ':'
		now = time.Now().Format("Jan-02-2006-15-04-05")
	}
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
		//folderPath example ./popularity/namespace/
		folderPath := path.Join(p.BaseOutputFilePath, p.Namespace)
		if p.Windows == true {
			p.fileNamePath = filepath.FromSlash(folderPath)
		} else {
			p.fileNamePath = ""
		}

		//jsonPath example ./popularity/namespace/filename
		jsonPath := path.Join(folderPath, filename) // relative filename
		if p.Windows == true {
			p.currentFileName = filepath.FromSlash(jsonPath)
		} else {
			p.currentFileName = jsonPath
		}

		// Invoke API to request for logs
		//
		// NOTE:
		//	This will run through a set of post functions if set,
		// after executing this API
		//
		// RequestReposForPageToken( ): Creates a HTTPRequest with some query
		// parameters and invokes it.
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
// Creates a HTTPRequest with some query parameters and invokes it.
func (p *Popularity) RequestReposForPageToken(pagetoken string) (PopularList, error) {
	// creating the request
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

	//resp contains the repo list in raw byte format in one page
	resp, err := req.Invoke()
	if err != nil {
		return PopularList{}, errors.Wrapf(
			err,
			"Failed to invoke popularity request: %q",
			req.URL,
		)
	}

	// Since `IsWriteToFile` is false so it **doesn't** call `WriteToFile`
	if p.IsWriteToFile {
		//writing the reponse in ./popularity/namespace/currentfileName.json
		p.WriteToFile(resp.Body(), p.currentFileName, p.fileNamePath)
		if err != nil {
			return PopularList{}, errors.Wrapf(
				err,
				"Failed to write popularity data to file: %q",
				p.currentFileName,
			)
		}
		log.Printf("Sucessfully wrote PopularList to file --------------> " + p.currentFileName)
	}

	// it is capable of holding the list of images
	// and the JSON is unmarshaled and returned.
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
	//returning the PopularList
	return out, nil
}

// WriteToFile creates a file with images having
// popularity ratings. This file is named with today's date.
// It writes the content of response body into passed filename with
// file mode 0644.
func (p *Popularity) WriteToFile(raw []byte, filename string, fpath string) error {
	if fpath != "" {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			os.MkdirAll(fpath, 0700) // Create your file
		}
	}

	errfile := ioutil.WriteFile(filename, raw, 0644)
	if errfile != nil {
		return errors.Wrapf(
			errfile,
			"Failed to write filename to %s",
			filename,
		)
	}
	return nil
}
