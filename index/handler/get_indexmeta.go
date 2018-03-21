//  Copyright (c) 2018 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"encoding/json"
	"github.com/mosuka/blast/index/client"
	"github.com/mosuka/blast/index/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type GetIndexMetaHandler struct {
	client *client.GRPCClient
}

func NewGetIndexMetaHandler(c *client.GRPCClient) *GetIndexMetaHandler {
	return &GetIndexMetaHandler{
		client: c,
	}
}

func (h *GetIndexMetaHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"host":           req.Host,
		"uri":            req.RequestURI,
		"method":         req.Method,
		"content_length": req.ContentLength,
		"remote_addr":    req.RemoteAddr,
		"header":         req.Header,
		"url":            req.URL,
	}).Info("")

	// request timeout
	requestTimeout := DefaultRequestTimeout
	if req.URL.Query().Get("requestTimeout") != "" {
		i, err := strconv.Atoi(req.URL.Query().Get("requestTimeout"))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to set batch size")

			Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requestTimeout = i
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(requestTimeout)*time.Millisecond)
	defer cancel()

	// request
	indexMeta, err := h.client.GetIndexMeta(ctx)
	resp := struct {
		IndexMeta *config.IndexConfig `json:"index_metag,omitempty"`
		Error     error               `json:"error,omitempty"`
	}{
		IndexMeta: indexMeta,
		Error:     err,
	}

	// output response
	output, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create response")

		Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(output)

	return
}
