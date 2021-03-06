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
	"github.com/buger/jsonparser"
	"github.com/mosuka/blast/node/client/grpc"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type BulkHandler struct {
	client *grpc.GRPCClient
}

func NewBulkHandler(c *grpc.GRPCClient) *BulkHandler {
	return &BulkHandler{
		client: c,
	}
}

func (h *BulkHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"req": req,
	}).Info("")

	// read request
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read request body")

		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get batch_size
	batchSize, err := jsonparser.GetInt(data, "batch_size")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to get batch size")

		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get requests
	requestsBytes, _, _, err := jsonparser.Get(data, "requests")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to get update requests")

		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var requests []map[string]interface{}
	err = json.Unmarshal(requestsBytes, &requests)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create update requests")

		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// overwrite request
	if req.URL.Query().Get("batchSize") != "" {
		i, err := strconv.Atoi(req.URL.Query().Get("batchSize"))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to set batch size")

			Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		batchSize = int64(i)
	}
	if batchSize <= 0 {
		batchSize = int64(DefaultBatchSize)
	}

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
	putCount, deleteCount, errorCount, err := h.client.Bulk(ctx, requests, int32(batchSize))
	resp := struct {
		PutCount    int32 `json:"put_count,omitempty"`
		DeleteCount int32 `json:"delete_count,omitempty"`
		ErrorCount  int32 `json:"error_count,omitempty"`
		Error       error `json:"error,omitempty"`
	}{
		PutCount:    putCount,
		DeleteCount: deleteCount,
		ErrorCount:  errorCount,
		Error:       err,
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
