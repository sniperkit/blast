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
	_ "github.com/mosuka/blast/dependency"
	"github.com/mosuka/blast/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"time"
)

type BlastService struct {
	IndexPath    string
	IndexMapping *mapping.IndexMappingImpl
	IndexType    string
	Kvstore      string
	Kvconfig     map[string]interface{}
	Index        bleve.Index
}

func NewBlastService(indexPath string, indexMapping *mapping.IndexMappingImpl, indexType string, kvstore string, kvconfig map[string]interface{}) *BlastService {
	return &BlastService{
		IndexPath:    indexPath,
		IndexMapping: indexMapping,
		IndexType:    indexType,
		Kvstore:      kvstore,
		Kvconfig:     kvconfig,
		Index:        nil,
	}
}

func (s *BlastService) OpenIndex() error {
	_, err := os.Stat(s.IndexPath)
	if os.IsNotExist(err) {
		s.Index, err = bleve.NewUsing(s.IndexPath, s.IndexMapping, s.IndexType, s.Kvstore, s.Kvconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"indexPath":    s.IndexPath,
				"indexMapping": s.IndexMapping,
				"indexType":    s.IndexType,
				"kvstore":      s.Kvstore,
				"kvconfig":     s.Kvconfig,
			}).Error(err.Error())

			return err
		}

		log.WithFields(log.Fields{
			"indexPath":    s.IndexPath,
			"indexMapping": s.IndexMapping,
			"indexType":    s.IndexType,
			"kvstore":      s.Kvstore,
			"kvconfig":     s.Kvconfig,
		}).Info("index created.")
	} else {
		s.Index, err = bleve.OpenUsing(s.IndexPath, s.Kvconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"indexPath": s.IndexPath,
				"kvconfig":  s.Kvconfig,
			}).Error(err.Error())

			return err
		}

		log.WithFields(log.Fields{
			"indexPath": s.IndexPath,
			"kvconfig":  s.Kvconfig,
		}).Info("index opened.")
	}

	return nil
}

func (s *BlastService) CloseIndex() error {
	err := s.Index.Close()
	if err != nil {
		log.Error(err.Error())

		return err
	}

	log.Info("index closed.")

	return nil
}

func (s *BlastService) GetIndexInfo(ctx context.Context, req *proto.GetIndexInfoRequest) (*proto.GetIndexInfoResponse, error) {
	protoGetIndexResponse := &proto.GetIndexInfoResponse{}

	if req.IndexPath {
		protoGetIndexResponse.IndexPath = s.IndexPath
	}

	if req.IndexMapping {
		indexMappingAny, err := proto.MarshalAny(s.IndexMapping)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		protoGetIndexResponse.IndexMapping = &indexMappingAny
	}

	if req.IndexType {
		protoGetIndexResponse.IndexType = s.IndexType
	}

	if req.Kvstore {
		protoGetIndexResponse.Kvstore = s.Kvstore
	}

	if req.Kvconfig {
		kvconfigAny, err := proto.MarshalAny(s.Kvconfig)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		protoGetIndexResponse.Kvconfig = &kvconfigAny
	}

	return protoGetIndexResponse, nil
}

func (s *BlastService) PutDocument(ctx context.Context, req *proto.PutDocumentRequest) (*proto.PutDocumentResponse, error) {
	fields, err := proto.UnmarshalAny(req.Fields)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     req.Id,
			"fields": req.Fields,
		}).Error(err.Error())

		return nil, err
	}

	err = s.Index.Index(req.Id, fields)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     req.Id,
			"fields": fields,
		}).Error(err.Error())

		return nil, err
	}

	return &proto.PutDocumentResponse{
		Id:     req.Id,
		Fields: req.Fields,
	}, nil
}

func (s *BlastService) GetDocument(ctx context.Context, req *proto.GetDocumentRequest) (*proto.GetDocumentResponse, error) {
	doc, err := s.Index.Document(req.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"id": req.Id,
		}).Error(err.Error())

		return nil, err
	}

	if doc == nil {
		log.WithFields(log.Fields{
			"id": req.Id,
		}).Info("document does not exist")

		return nil, nil
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

	fieldsAny, err := proto.MarshalAny(fields)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     req.Id,
			"fields": fields,
		}).Error(err.Error())

		return nil, err
	}

	return &proto.GetDocumentResponse{
		Id:     req.Id,
		Fields: &fieldsAny,
	}, nil
}

func (s *BlastService) DeleteDocument(ctx context.Context, req *proto.DeleteDocumentRequest) (*proto.DeleteDocumentResponse, error) {
	err := s.Index.Delete(req.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"id": req.Id,
		}).Error(err.Error())

		return nil, err
	}

	return &proto.DeleteDocumentResponse{
		Id: req.Id,
	}, nil
}

func (s *BlastService) Bulk(ctx context.Context, req *proto.BulkRequest) (*proto.BulkResponse, error) {
	var (
		processedCount   int32
		putCount         int32
		putErrorCount    int32
		deleteCount      int32
		methodErrorCount int32
	)

	batch := s.Index.NewBatch()

	for _, updateRequest := range req.UpdateRequests {
		processedCount++

		switch updateRequest.Method {
		case "put":
			fields, err := proto.UnmarshalAny(updateRequest.Document.Fields)
			if err != nil {
				log.WithFields(log.Fields{
					"updateRequest": updateRequest,
				}).Warn(err.Error())

				putErrorCount++

				continue
			}

			err = batch.Index(updateRequest.Document.Id, fields)
			if err != nil {
				log.WithFields(log.Fields{
					"id":     updateRequest.Document.Id,
					"fields": fields,
				}).Warn(err.Error())

				putErrorCount++

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

			methodErrorCount++

			continue
		}

		if processedCount%req.BatchSize == 0 {
			err := s.Index.Batch(batch)
			if err != nil {
				log.Warn(err.Error())
			}

			batch = s.Index.NewBatch()
		}
	}

	if batch.Size() > 0 {
		err := s.Index.Batch(batch)
		if err != nil {
			log.Warn(err.Error())
		}
	}

	return &proto.BulkResponse{
		PutCount:         putCount,
		PutErrorCount:    putErrorCount,
		DeleteCount:      deleteCount,
		MethodErrorCount: methodErrorCount,
	}, nil
}

func (s *BlastService) Search(ctx context.Context, req *proto.SearchRequest) (*proto.SearchResponse, error) {
	searchRequest, err := proto.UnmarshalAny(req.SearchRequest)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	searchResult, err := s.Index.Search(searchRequest.(*bleve.SearchRequest))
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	searchResultAny, err := proto.MarshalAny(searchResult)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	return &proto.SearchResponse{
		SearchResult: &searchResultAny,
	}, nil
}
