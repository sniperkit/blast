//  Copyright (c) 2017 Minoru Osuka
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

package service

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mosuka/blast/index"
	_ "github.com/mosuka/blast/node/builtin"
	"github.com/mosuka/blast/protobuf"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"time"
)

type IndexService struct {
	indexPath    string
	indexMapping *mapping.IndexMappingImpl
	indexMeta    *index.IndexMeta
	index        bleve.Index
}

func NewIndexService(indexPath string, indexMapping *mapping.IndexMappingImpl, indexMeta *index.IndexMeta) (*IndexService, error) {
	return &IndexService{
		indexPath:    indexPath,
		indexMapping: indexMapping,
		indexMeta:    indexMeta,
		index:        nil,
	}, nil
}

func (s *IndexService) OpenIndex() error {
	_, err := os.Stat(s.indexPath)
	if os.IsNotExist(err) {
		s.index, err = bleve.NewUsing(s.indexPath, s.indexMapping, s.indexMeta.IndexType, s.indexMeta.Storage, s.indexMeta.Config)
		if err != nil {
			log.WithFields(log.Fields{
				"indexPath":    s.indexPath,
				"indexMapping": s.indexMapping,
				"indexMeta":    s.indexMeta,
				"error":        err.Error(),
			}).Error("failed to create new index.")

			return err
		}

		log.WithFields(log.Fields{
			"indexPath":    s.indexPath,
			"indexMapping": s.indexMapping,
			"indexMeta":    s.indexMeta,
		}).Info("new index was created.")
	} else {
		s.index, err = bleve.OpenUsing(s.indexPath, s.indexMeta.Config)
		if err != nil {
			log.WithFields(log.Fields{
				"indexPath":     s.indexPath,
				"runtimeConfig": s.indexMeta.Config,
				"error":         err.Error(),
			}).Error("failed to open existing index.")

			return err
		}

		log.WithFields(log.Fields{
			"indexPath":     s.indexPath,
			"runtimeConfig": s.indexMeta.Config,
		}).Info("existing index was opened.")
	}

	return nil
}

func (s *IndexService) CloseIndex() error {
	err := s.index.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"indexPath":     s.indexPath,
			"runtimeConfig": s.indexMeta.Config,
			"error":         err.Error(),
		}).Error("failed to close index.")

		return err
	}

	log.Info("index was closed.")

	return nil
}

func (s *IndexService) PutDocument(ctx context.Context, req *protobuf.PutDocumentRequest) (*empty.Empty, error) {
	fields, err := protobuf.UnmarshalAny(req.Fields)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    req.Id,
			"error": err.Error(),
		}).Error("failed to unmarshal fields.")

		return nil, err
	}

	err = s.index.Index(req.Id, fields)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     req.Id,
			"fields": fields,
			"error":  err.Error(),
		}).Error("failed to index document.")

		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *IndexService) GetDocument(ctx context.Context, req *protobuf.GetDocumentRequest) (*protobuf.GetDocumentResponse, error) {
	doc, err := s.index.Document(req.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    req.Id,
			"error": err.Error(),
		}).Error("failed to get document.")

		return nil, err
	}

	if doc == nil {
		log.WithFields(log.Fields{
			"id": req.Id,
		}).Info("document does not exist")

		return &protobuf.GetDocumentResponse{}, nil
	}

	fields := make(map[string]interface{})
	for _, field := range doc.Fields {
		var value interface{}

		switch field := field.(type) {
		case *document.TextField:
			value = string(field.Value())
		case *document.NumericField:
			numValue, err := field.Number()
			if err == nil {
				value = numValue
			}
		case *document.DateTimeField:
			dateValue, err := field.DateTime()
			if err == nil {
				dateValue.Format(time.RFC3339Nano)
				value = dateValue
			}
		}

		existedField, existed := fields[field.Name()]
		if existed {
			switch existedField := existedField.(type) {
			case []interface{}:
				fields[field.Name()] = append(existedField, value)
			case interface{}:
				arr := make([]interface{}, 2)
				arr[0] = existedField
				arr[1] = value
				fields[field.Name()] = arr
			}
		} else {
			fields[field.Name()] = value
		}
	}

	fieldsAny, err := protobuf.MarshalAny(fields)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     req.Id,
			"fields": fields,
			"error":  err.Error(),
		}).Error("failed to marshal fields.")

		return nil, err
	}

	return &protobuf.GetDocumentResponse{
		Id:     req.Id,
		Fields: &fieldsAny,
	}, nil
}

func (s *IndexService) DeleteDocument(ctx context.Context, req *protobuf.DeleteDocumentRequest) (*empty.Empty, error) {
	err := s.index.Delete(req.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"id": req.Id,
		}).Error(err.Error())

		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *IndexService) Bulk(ctx context.Context, req *protobuf.BulkRequest) (*protobuf.BulkResponse, error) {
	var (
		processedCount int32
		putCount       int32
		deleteCount    int32
		errorCount     int32
	)

	batch := s.index.NewBatch()

	for _, updateRequest := range req.Requests {
		processedCount++

		switch updateRequest.Method {
		case "put":
			fields, err := protobuf.UnmarshalAny(updateRequest.Document.Fields)
			if err != nil {
				log.WithFields(log.Fields{
					"updateRequest": updateRequest,
				}).Warn(err.Error())

				errorCount++

				continue
			}

			err = batch.Index(updateRequest.Document.Id, fields)
			if err != nil {
				log.WithFields(log.Fields{
					"id":     updateRequest.Document.Id,
					"fields": fields,
				}).Warn(err.Error())

				errorCount++

				continue
			}

			putCount++
		case "delete":
			batch.Delete(updateRequest.Document.Id)

			deleteCount++
		default:
			log.WithFields(log.Fields{
				"method": updateRequest.Method,
			}).Warn("unknown method")

			errorCount++

			continue
		}

		if processedCount%req.BatchSize == 0 {
			err := s.index.Batch(batch)
			if err != nil {
				log.Warn(err.Error())
			}

			batch = s.index.NewBatch()
		}
	}

	if batch.Size() > 0 {
		err := s.index.Batch(batch)
		if err != nil {
			log.Warn(err.Error())
		}
	}

	return &protobuf.BulkResponse{
		PutCount:    putCount,
		DeleteCount: deleteCount,
		ErrorCount:  errorCount,
	}, nil
}

func (s *IndexService) Search(ctx context.Context, req *protobuf.SearchRequest) (*protobuf.SearchResponse, error) {
	searchRequest, err := protobuf.UnmarshalAny(req.SearchRequest)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	searchResult, err := s.index.Search(searchRequest.(*bleve.SearchRequest))
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	searchResultAny, err := protobuf.MarshalAny(searchResult)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	return &protobuf.SearchResponse{
		SearchResult: &searchResultAny,
	}, nil
}

func (s *IndexService) GetIndexMapping(ctx context.Context, req *empty.Empty) (*protobuf.GetIndexMappingResponse, error) {
	// IndexMapping -> Any
	any, err := protobuf.MarshalAny(s.index.Mapping())
	if err != nil {
		return nil, err
	}

	return &protobuf.GetIndexMappingResponse{
		IndexMapping: &any,
	}, nil
}

func (s *IndexService) GetIndexMeta(ctx context.Context, req *empty.Empty) (*protobuf.GetIndexMetaResponse, error) {
	// IndexMeta -> Any
	any, err := protobuf.MarshalAny(s.indexMeta)
	if err != nil {
		return nil, err
	}

	return &protobuf.GetIndexMetaResponse{
		IndexMeta: &any,
	}, nil
}
