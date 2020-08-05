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

// import (
// 	"context"
// 	"net"
// 	"strings"
// )

const (
	// MonDDYYYYDateFormat is used for file names where these files are
	// used to store various quay api logs
	MonDDYYYYDateFormat string = "Jan-02-2006"

	// QuayLogDateFormat is the format found in quay logs
	QuayLogDateFormat string = "02 Jan 2006"
)

// Popular holds the fields that represent an image
// with its popularity rating
type Popular struct {
	Kind        string  `json:"kind"`
	Name        string  `json:"name"`
	Popularity  float64 `json:"popularity"`
	Namespace   string  `json:"namespace"`
	State       string  `json:"state"`
	IsPublic    bool    `json:"is_public"`
	IsStarred   bool    `json:"is_starred"`
	Description string  `json:"description"`
}

// PopularList holds a list of Popular items
type PopularList struct {
	Items []Popular `json:"repositories"`
}

// ResolvedIP represents quay.io repo's resolved IP details
type ResolvedIP struct {
	CountryISOCode string `json:"country_iso_code"`
	SyncToken      string `json:"sync_token"`
	Service        string `json:"service"`
	Provider       string `json:"provider"`
}

// Metadata represents quay.io repo's metadata
type Metadata struct {
	Repo       string     `json:"repo"`
	Tag        string     `json:"tag"`
	Namespace  string     `json:"namespace"`
	ResolvedIP ResolvedIP `json:"resolved_ip"`
}

// Log represents a quay.io repo's logs
type Log struct {
	IP       string   `json:"ip"`
	Kind     string   `json:"kind"`
	Datetime string   `json:"datetime"`
	Metadata Metadata `json:"metadata"`
}

// LogList holds a list of Log items
type LogList struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	NextPage  string `json:"next_page"`
	Items     []Log  `json:"logs"`
}