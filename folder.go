/*
Copyright 2020 The MayaData Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package growthmetrics

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
)

// FolderConfig is used to initialise Folder
type FolderConfig struct {
	Path  string
	Debug bool
}

// Folder represents the path to a directory
type Folder struct {
	Path  string
	Debug bool
}

// NewFolder returns a new instance of Folder
func NewFolder(config FolderConfig) *Folder {
	return &Folder{
		Path:  config.Path,
		Debug: config.Debug,
	}
}

// ListJSONFiles lists all json files
func (f *Folder) ListJSONFiles() ([]string, error) {
	files, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		log.Printf("No file(s) found: Path %s", f.Path)
		return nil, nil
	}
	var out []string
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() || file.Mode().IsDir() {
			// we don't want to load directory
			continue
		}
		if !strings.HasSuffix(fileName, ".json") {
			// we support json file only
			continue
		}
		fileNameWithPath := path.Join(f.Path, fileName)
		out = append(out, fileNameWithPath)
	}
	if f.Debug {
		log.Printf(
			"Found json files: file-count %d: path %s",
			len(out),
			f.Path,
		)
	}
	return out, nil
}
