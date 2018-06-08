# Blast

The Blast is a full text search and indexing server written in [Go](https://golang.org) that built on top of the [Bleve](http://www.blevesearch.com). Blast server provides functions through [gRPC](http://www.grpc.io) ([HTTP/2](https://en.wikipedia.org/wiki/HTTP/2) + [Protocol Buffers](https://developers.google.com/protocol-buffers/)).  
This repository includes Blast CLI that is a command line interface for controlling Blast server, and Blast REST server that provides a traditional [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) API ([HTTP/1.1](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol) + [JSON](http://www.json.org)).  
Blast makes it easy for programmers to develop search applications with advanced features.


## Features

- Full-text search and indexing
- Faceting
- Result highlighting
- Text analysis

For more detailed information, refer to the [Bleve document](http://www.blevesearch.com/docs/Home/).

## Blast node

### Blast node config file (blast_node.yaml)

The Blast server parameters are described in blast.yaml.

```yaml
log_format: text
log_level: info
log_output:

grpc_listen_address: 0.0.0.0:5000
http_listen_address: 0.0.0.0:8000

index_path: ./data/index
index_mapping_path: ./example/config/index_mapping.json
index_meta_path: ./example/config/index_meta.json

rest_uri: /rest
metrics_uri: /metrics
```

### Blast node config parameters and environment variables

| Config parameter name | Environment variable | Command line option | Description | Default value |
| --- | --- | --- | --- | --- |
| -                        | BLAST_NODE_CONFIG_PATH              | --config-path              | config path. | - |
| log_format               | BLAST_NODE_LOG_FORMAT               | --log-format               | log format. `text`, `color` and `json` are available | `text` |
| log_level                | BLAST_NODE_LOG_LEVEL                | --log-level                | log level. `debug`, `info`, `warn`, `error`, `fatal` and `panic` are available. | `info` |
| log_output               | BLAST_NODE_LOG_OUTPUT               | --log-output               | log output path. | `""` |
| grpc_listen_address      | BLAST_NODE_GRPC_LISTEN_ADDRESS      | --grpc-listen-address      | address to listen for the gRPC. | `0.0.0.0:5000` |
| http_listen_address      | BLAST_NODE_HTTP_LISTEN_ADDRESS      | --http-listen-address      | address to listen for the HTTP. | `0.0.0.0:8000` |
| index_path               | BLAST_NODE_INDEX_PATH               | --index-path               | index directory path. Default is `./data/index` | `./data/index` |
| index_mapping_path       | BLAST_NODE_INDEX_MAPPING_PATH       | --index-mapping-path       | index mapping path. | `""` |
| index_meta_path          | BLAST_NODE_INDEX_META_PATH          | --index-meta-path          | index meta path. | `""` |
| rest_uri                 | BLAST_NODE_REST_URI                 | --rest-uri                 | base URI for REST API endpoint. | `/rest` |
| metrics_uri              | BLAST_NODE_METRICS_URI              | --metrics-uri              | URI for metrics exposition endpoint. | `/metrics` |


### Index Mapping

You can specify the index mapping describes how to your data model should be indexed. it contains all of the details about which fields your documents can contain, and how those fields should be dealt with when adding documents to the index, or when querying those fields. The example is following:

`./example/config/index_mapping.yaml`

See [Introduction to Index Mappings](http://www.blevesearch.com/docs/Index-Mapping/) and [type IndexMappingImpl](https://godoc.org/github.com/blevesearch/bleve/mapping#IndexMappingImpl) for more details.  


### Start Blast node

The `blast start node` command starts Blast node. You can display a help message by specifying `-h` or `--help` option.

```sh
$ ./bin/blast start node
```


## Blast CLI

### Get the index information from Blast server

The `get index` command retrieves an index information about existing opened index. You can display a help message by specifying the `- h` or` --help` option.

```sh
$ ./bin/blastcli get index --include-index-mapping --include-index-type --include-kvstore --include-kvconfig
```

The result of the above `get index` command is:

```json
{
  "path": "/var/blast/data/index",
  "index_mapping": {
    "types": {
      "document": {
        "enabled": true,
        "dynamic": true,
        "properties": {
          "category": {
            "enabled": true,
            "dynamic": true,
            "fields": [
              {
                "type": "text",
                "analyzer": "keyword",
                "store": true,
                "index": true,
                "include_term_vectors": true,
                "include_in_all": true
              }
            ],
            "default_analyzer": ""
          },
          "description": {
            "enabled": true,
            "dynamic": true,
            "fields": [
              {
                "type": "text",
                "analyzer": "en",
                "store": true,
                "index": true,
                "include_term_vectors": true,
                "include_in_all": true
              }
            ],
            "default_analyzer": ""
          },
          "name": {
            "enabled": true,
            "dynamic": true,
            "fields": [
              {
                "type": "text",
                "analyzer": "en",
                "store": true,
                "index": true,
                "include_term_vectors": true,
                "include_in_all": true
              }
            ],
            "default_analyzer": ""
          },
          "popularity": {
            "enabled": true,
            "dynamic": true,
            "fields": [
              {
                "type": "number",
                "store": true,
                "index": true,
                "include_in_all": true
              }
            ],
            "default_analyzer": ""
          },
          "release": {
            "enabled": true,
            "dynamic": true,
            "fields": [
              {
                "type": "datetime",
                "store": true,
                "index": true,
                "include_in_all": true
              }
            ],
            "default_analyzer": ""
          },
          "type": {
            "enabled": true,
            "dynamic": true,
            "fields": [
              {
                "type": "text",
                "analyzer": "keyword",
                "store": true,
                "index": true,
                "include_term_vectors": true,
                "include_in_all": true
              }
            ],
            "default_analyzer": ""
          }
        },
        "default_analyzer": ""
      }
    },
    "default_mapping": {
      "enabled": true,
      "dynamic": true,
      "default_analyzer": ""
    },
    "type_field": "type",
    "default_type": "document",
    "default_analyzer": "standard",
    "default_datetime_parser": "dateTimeOptional",
    "default_field": "_all",
    "store_dynamic": true,
    "index_dynamic": true,
    "analysis": {}
  },
  "index_type": "upside_down",
  "kvstore": "boltdb",
  "kvconfig": {}
}
```

### Put document format

```json
{
  "id": "1",
  "fields": {
    "name": "Bleve",
    "description": "Bleve is a full-text search and indexing library for Go.",
    "category": "Library",
    "popularity": 3.0,
    "release": "2014-04-18T00:00:00Z",
    "type": "document"
  }
}
```


### Put the document to Blast server

The `put document` command adds or updates a JSON formatted document in a specified index. You can display a help message by specifying the `- h` or` --help` option.  
The document example is following:

```sh
$ cat ./example/document_1.json | xargs -0 ./bin/blast put document --request
```

The result of the above `put document` command is:

```json
{
  "id": "1",
  "fields": {
    "category": "Library",
    "description": "Bleve is a full-text search and indexing library for Go.",
    "name": "Bleve",
    "popularity": 3,
    "release": "2014-04-18T00:00:00Z",
    "type": "document"
  }
}
```


### Get the document from Blast server

The `get document` command retrieves a JSON formatted document on its id from a specified index. You can display a help message by specifying the `- h` or` --help` option.

```sh
$ ./bin/blast get document --id 1
```

The result of the above `get document` command is:

```json
{
  "id": "1",
  "fields": {
    "category": "Library",
    "description": "Bleve is a full-text search and indexing library for Go.",
    "name": "Bleve",
    "popularity": 3,
    "release": "2014-04-18T00:00:00Z",
    "type": "document"
  }
}
```


### Delete the document from Blast server

The `delete document` command deletes a document on its id from a specified index. You can display a help message by specifying the `- h` or` --help` option.

```sh
$ ./bin/blast delete document --id 1
```

The result of the above `delete document` command is:

```json
{
  "id": "1"
}
```


### Bulk document format

```json
{
  "batch_size": 1000,
  "requests": [
    {
      "method": "put",
      "document": {
        "id": "1",
        "fields": {
          "name": "Bleve",
          "description": "Bleve is a full-text search and indexing library for Go.",
          "category": "Library",
          "popularity": 3,
          "release": "2014-04-18T00:00:00Z",
          "type": "document"
        }
      }
    },
    {
      "method": "delete",
      "document": {
        "id": "2"
      }
    },
    {
      "method": "delete",
      "document": {
        "id": "3"
      }
    },
    {
      "method": "delete",
      "document": {
        "id": "4"
      }
    },
    {
      "method": "delete",
      "document": {
        "id": "5"
      }
    },
    {
      "method": "delete",
      "document": {
        "id": "6"
      }
    },
    {
      "method": "put",
      "document": {
        "id": "7",
        "fields": {
          "name": "Blast",
          "description": "Blast is a full-text search and indexing server written in Go, built on top of Bleve.",
          "category": "Server",
          "popularity": 1,
          "release": "2017-01-13T00:00:00Z",
          "type": "document"
        }
      }
    }
  ]
}
```

### Index the documents in bulk to Blast server

The `bulk` command makes it possible to perform many put/delete operations in a single command execution. This can greatly increase the indexing speed. You can display a help message by specifying the `- h` or` --help` option.
The bulk example is following:

```sh
$ cat ./example/bulk_put.json | xargs -0 ./bin/blast bulk --request
```

The result of the above `bulk` command is:

```json
{
  "put_count": 7
}
```

### Search request format

```json
{
  "search_request": {
    "query": {
      "query": "name:*"
    },
    "size": 10,
    "from": 0,
    "fields": [
      "*"
    ],
    "sort": [
      "-_score"
    ],
    "facets": {
      "Category count": {
        "size": 10,
        "field": "category"
      },
      "Popularity range": {
        "size": 10,
        "field": "popularity",
        "numeric_ranges": [
          {
            "name": "less than 1",
            "max": 1
          },
          {
            "name": "more than or equal to 1 and less than 2",
            "min": 1,
            "max": 2
          },
          {
            "name": "more than or equal to 2 and less than 3",
            "min": 2,
            "max": 3
          },
          {
            "name": "more than or equal to 3 and less than 4",
            "min": 3,
            "max": 4
          },
          {
            "name": "more than or equal to 4 and less than 5",
            "min": 4,
            "max": 5
          },
          {
            "name": "more than or equal to 5",
            "min": 5
          }
        ]
      },
      "Release date range": {
        "size": 10,
        "field": "release",
        "date_ranges": [
          {
            "name": "2001 - 2010",
            "start": "2001-01-01T00:00:00Z",
            "end": "2010-12-31T23:59:59Z"
          },
          {
            "name": "2011 - 2020",
            "start": "2011-01-01T00:00:00Z",
            "end": "2020-12-31T23:59:59Z"
          }
        ]
      }
    },
    "highlight": {
      "style": "html",
      "fields": [
        "name",
        "description"
      ]
    }
  }
}
```

See [Queries](http://www.blevesearch.com/docs/Query/), [Query String Query](http://www.blevesearch.com/docs/Query-String-Query/) and [type SearchRequest](https://godoc.org/github.com/blevesearch/bleve#SearchRequest) for more details.


### Search the documents from Blast server

The `search` command can be executed with a search request, which includes the Query, within its file. Here is an example:
You can display a help message by specifying the `- h` or` --help` option.

```sh
$ cat ./example/search_request.json | xargs -0 ./bin/blast search --request
```

The result of the above `search` command is:

```json
{
  "search_result": {
    "status": {
      "total": 1,
      "failed": 0,
      "successful": 1
    },
    "request": {
      "query": {
        "query": "name:*"
      },
      "size": 10,
      "from": 0,
      "highlight": {
        "style": "html",
        "fields": [
          "name",
          "description"
        ]
      },
      "fields": [
        "*"
      ],
      "facets": {
        "Category count": {
          "size": 10,
          "field": "category"
        },
        "Popularity range": {
          "size": 10,
          "field": "popularity",
          "numeric_ranges": [
            {
              "name": "less than 1",
              "max": 1
            },
            {
              "name": "more than or equal to 1 and less than 2",
              "min": 1,
              "max": 2
            },
            {
              "name": "more than or equal to 2 and less than 3",
              "min": 2,
              "max": 3
            },
            {
              "name": "more than or equal to 3 and less than 4",
              "min": 3,
              "max": 4
            },
            {
              "name": "more than or equal to 4 and less than 5",
              "min": 4,
              "max": 5
            },
            {
              "name": "more than or equal to 5",
              "min": 5
            }
          ]
        },
        "Release date range": {
          "size": 10,
          "field": "release",
          "date_ranges": [
            {
              "end": "2010-12-31T23:59:59Z",
              "name": "2001 - 2010",
              "start": "2001-01-01T00:00:00Z"
            },
            {
              "end": "2020-12-31T23:59:59Z",
              "name": "2011 - 2020",
              "start": "2011-01-01T00:00:00Z"
            }
          ]
        }
      },
      "explain": false,
      "sort": [
        "-_score"
      ],
      "includeLocations": false
    },
    "hits": [
      {
        "index": "./data/index",
        "id": "1",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "bleve": [
              {
                "pos": 1,
                "start": 0,
                "end": 5,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Bleve is a full-text search and indexing library for Go."
          ],
          "name": [
            "\u003cmark\u003eBleve\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Library",
          "description": "Bleve is a full-text search and indexing library for Go.",
          "name": "Bleve",
          "popularity": 3,
          "release": "2014-04-18T00:00:00Z",
          "type": "document"
        }
      },
      {
        "index": "./data/index",
        "id": "2",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "lucene": [
              {
                "pos": 1,
                "start": 0,
                "end": 6,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Apache Lucene is a high-performance, full-featured text search engine library written entirely in Java."
          ],
          "name": [
            "\u003cmark\u003eLucene\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Library",
          "description": "Apache Lucene is a high-performance, full-featured text search engine library written entirely in Java.",
          "name": "Lucene",
          "popularity": 4,
          "release": "2000-03-30T00:00:00Z",
          "type": "document"
        }
      },
      {
        "index": "./data/index",
        "id": "3",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "whoosh": [
              {
                "pos": 1,
                "start": 0,
                "end": 6,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Whoosh is a fast, featureful full-text indexing and searching library implemented in pure Python."
          ],
          "name": [
            "\u003cmark\u003eWhoosh\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Library",
          "description": "Whoosh is a fast, featureful full-text indexing and searching library implemented in pure Python.",
          "name": "Whoosh",
          "popularity": 3,
          "release": "2008-02-20T00:00:00Z",
          "type": "document"
        }
      },
      {
        "index": "./data/index",
        "id": "4",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "ferret": [
              {
                "pos": 1,
                "start": 0,
                "end": 6,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Ferret is a super fast, highly configurable search library written in Ruby."
          ],
          "name": [
            "\u003cmark\u003eFerret\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Library",
          "description": "Ferret is a super fast, highly configurable search library written in Ruby.",
          "name": "Ferret",
          "popularity": 2,
          "release": "2005-10-01T00:00:00Z",
          "type": "document"
        }
      },
      {
        "index": "./data/index",
        "id": "5",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "solr": [
              {
                "pos": 1,
                "start": 0,
                "end": 4,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Solr is an open source enterprise search platform, written in Java, from the Apache Lucene project."
          ],
          "name": [
            "\u003cmark\u003eSolr\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Server",
          "description": "Solr is an open source enterprise search platform, written in Java, from the Apache Lucene project.",
          "name": "Solr",
          "popularity": 5,
          "release": "2006-12-22T00:00:00Z",
          "type": "document"
        }
      },
      {
        "index": "./data/index",
        "id": "6",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "elasticsearch": [
              {
                "pos": 1,
                "start": 0,
                "end": 13,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Elasticsearch is a search engine based on Lucene, written in Java."
          ],
          "name": [
            "\u003cmark\u003eElasticsearch\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Server",
          "description": "Elasticsearch is a search engine based on Lucene, written in Java.",
          "name": "Elasticsearch",
          "popularity": 5,
          "release": "2010-02-08T00:00:00Z",
          "type": "document"
        }
      },
      {
        "index": "./data/index",
        "id": "7",
        "score": 0.12163776688600772,
        "locations": {
          "name": {
            "blast": [
              {
                "pos": 1,
                "start": 0,
                "end": 5,
                "array_positions": null
              }
            ]
          }
        },
        "fragments": {
          "description": [
            "Blast is a full-text search and indexing server written in Go, built on top of Bleve."
          ],
          "name": [
            "\u003cmark\u003eBlast\u003c/mark\u003e"
          ]
        },
        "sort": [
          "_score"
        ],
        "fields": {
          "category": "Server",
          "description": "Blast is a full-text search and indexing server written in Go, built on top of Bleve.",
          "name": "Blast",
          "popularity": 1,
          "release": "2017-01-13T00:00:00Z",
          "type": "document"
        }
      }
    ],
    "total_hits": 7,
    "max_score": 0.12163776688600772,
    "took": 639906,
    "facets": {
      "Category count": {
        "field": "category",
        "total": 7,
        "missing": 0,
        "other": 0,
        "terms": [
          {
            "term": "library",
            "count": 4
          },
          {
            "term": "server",
            "count": 3
          }
        ]
      },
      "Popularity range": {
        "field": "popularity",
        "total": 7,
        "missing": 0,
        "other": 0,
        "numeric_ranges": [
          {
            "name": "more than or equal to 3 and less than 4",
            "min": 3,
            "max": 4,
            "count": 2
          },
          {
            "name": "more than or equal to 5",
            "min": 5,
            "count": 2
          },
          {
            "name": "more than or equal to 1 and less than 2",
            "min": 1,
            "max": 2,
            "count": 1
          },
          {
            "name": "more than or equal to 2 and less than 3",
            "min": 2,
            "max": 3,
            "count": 1
          },
          {
            "name": "more than or equal to 4 and less than 5",
            "min": 4,
            "max": 5,
            "count": 1
          }
        ]
      },
      "Release date range": {
        "field": "release",
        "total": 6,
        "missing": 0,
        "other": 0,
        "date_ranges": [
          {
            "name": "2001 - 2010",
            "start": "2001-01-01T00:00:00Z",
            "end": "2010-12-31T23:59:59Z",
            "count": 4
          },
          {
            "name": "2011 - 2020",
            "start": "2011-01-01T00:00:00Z",
            "end": "2020-12-31T23:59:59Z",
            "count": 2
          }
        ]
      }
    }
  }
}
```


## License

Apache License Version 2.0
