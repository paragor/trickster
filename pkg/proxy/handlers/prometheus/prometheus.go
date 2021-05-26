/*
 * Copyright 2018 The Trickster Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prometheus

import (
	"net/http"
	"net/url"

	bo "github.com/trickstercache/trickster/pkg/backends/options"
	"github.com/trickstercache/trickster/pkg/backends/prometheus"
	co "github.com/trickstercache/trickster/pkg/cache/options"
	"github.com/trickstercache/trickster/pkg/cache/registration"
	"github.com/trickstercache/trickster/pkg/routing"

	"github.com/gorilla/mux"
)

// NewAccelerator returns a new Prometheus Accelerator. only baseURL is required
func NewAccelerator(baseURL string) (http.Handler, error) {
	return NewAcceleratorWithOptions(baseURL, nil, nil)
}

// NewAcceleratorWithOptions returns a new Prometheus Accelerator. only baseURL is required
func NewAcceleratorWithOptions(baseURL string, o *bo.Options, c *co.Options) (http.Handler, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if c == nil {
		c = co.New()
		c.Name = "default"
	}
	cache := registration.NewCache(c.Name, c, nil)
	err = cache.Connect()
	if err != nil {
		return nil, err
	}
	if o == nil {
		o = bo.New()
		o.Name = "default"
	}
	o.Provider = "prometheus"
	o.CacheName = c.Name
	o.Scheme = u.Scheme
	o.Host = u.Host
	o.PathPrefix = u.Path
	r := mux.NewRouter()
	cl, err := prometheus.NewClient("default", o, mux.NewRouter(), cache, nil, nil)
	if err != nil {
		return nil, err
	}
	o.HTTPClient = cl.HTTPClient()
	routing.RegisterPathRoutes(r, cl.Handlers(), cl, o, cache, cl.DefaultPathConfigs(o), nil, "", nil)
	return o.Router, nil
}
