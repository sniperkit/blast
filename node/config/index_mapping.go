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

package config

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

func LoadIndexMapping(reader io.Reader) (*mapping.IndexMappingImpl, error) {
	indexMapping := mapping.NewIndexMapping()

	resourceBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(fmt.Sprintf("failed to read index mapping."))
		return nil, err
	}

	err = json.Unmarshal(resourceBytes, indexMapping)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(fmt.Sprintf("failed to unmarshal index mapping."))
		return nil, err
	}

	return indexMapping, nil
}
