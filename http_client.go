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
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

const (
	// POST based http request
	POST string = "post"

	// GET based http request
	GET string = "get"
)

// HTTPRequest defines the configuration required
// to invoke http request
type HTTPRequest struct {
	URL         string            `json:"url"`
	Method      string            `json:"method,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"queryParams,omitempty"`
	PathParams  map[string]string `json:"pathParams,omitempty"`
	Body        string            `json:"body,omitempty"`
	AuthToken   string            `json:"authToken,omitempty"`
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	OutputFile  string            `json:"outputFile"`
}

// Invoke invokes http calls
func (r *HTTPRequest) Invoke() (*resty.Response, error) {
	req := resty.New().R().
		SetBasicAuth(r.Username, r.Password).
		SetAuthToken(r.AuthToken).
		SetBody(r.Body).
		SetHeaders(r.Headers).
		SetQueryParams(r.QueryParams).
		SetPathParams(r.PathParams)

	if r.OutputFile != "" {
		req.SetOutput(r.OutputFile)
	}

	switch strings.ToLower(r.Method) {
	case POST:
		return req.Post(r.URL)
	case GET:
		return req.Get(r.URL)
	default:
		return nil, errors.Errorf(
			"Unsupported http method %q",
			r.Method,
		)
	}
}
