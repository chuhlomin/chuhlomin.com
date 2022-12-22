# search

[![main](https://github.com/chuhlomin/search/actions/workflows/main.yml/badge.svg)](https://github.com/chuhlomin/search/actions/workflows/main.yml)
[![release](https://github.com/chuhlomin/search/actions/workflows/release.yml/badge.svg)](https://github.com/chuhlomin/search/actions/workflows/release.yml)

`search` is a project to provide a simple search engine,
built on top of the [bleve](https://github.com/blevesearch/bleve) Go library.

Suitable for a small projects, like blogs.

It concist of two parts:

* `indexer` Go struct, which can be used to create index,
* `server` Go app, which reads index and serves search requests.

## Indexer

Initialize an indexer with `NewIndexer`:

```go
import "github.com/chuhlomin/search"

indexer, err := search.NewIndexer(searchIndexPath, buildPathPrefix)
```

Register types to index with `RegisterType`:

```go
err := indexer.RegisterType(someStruct{}, "en")
```

Struct may implement bleve's `mapping.Classifier` interface to specify a type name.
Otherwise, struct name will be used as type name.

```go
type Classifier interface {
	Type() string
}
```

Also, structs may implement `Lanuguage` interface to specify language
for all text fields in the struct.
Otherwise, the default language will be used (passed to `RegisterType`).

```go
type Language interface {
	Language() string
}
```

Use struct tags to define a special behaviour:

```go
type someStruct struct {
	SomeField string `indexer:"text"`
	Date      string `indexer:"date"`
	Password  string `indexer:"no_store"`
	Blob      string `indexer:"no_index"`
}
```

Index documents with `Index`:

```go
err := indexer.Index("id", someStruct{SomeField: "needle & needle"})
```

Don't forget to call `Close` when you're done:

```go
err := indexer.Close()
```

## Server

You may run the `server` container with index mounted:

```yml
version: "3.9"

services:
  search:
    image: chuhlomin/search:v0.1.0
    ports:
    - 127.0.0.1:8081:80
    environment:
    - INDEX_PATH=/index
    - BIND=0.0.0.0:80
    volumes:
    - ./index:/index
```

Then send a search query:

```bash
curl "http://127.0.0.1:8081/?q=needle"
```

Result will be a JSON array of documents:

```json
[
    {
        "id": "id",
        "score": 1.0,
        "fragments": {
            "field": "SomeField",
            "locations": {
                "needle": [
                    {
                        "start": 0,
                        "end": 6
                    },
                    {
                        "start": 10,
                        "end": 16
                    }
                ]
            }
        },
        "document": {}
    }
]
```

You may also specify which document fields to return:

```bash
## Search Local
curl -X "POST" "http://127.0.0.1:8081/?q=needle" \
     -d $'{
  "SomeField": true
}
'
```

```json
[
    {
        "id": "id",
        "score": 1.0,
        "fragments": {
            "field": "SomeField",
            "locations": {
                "needle": [
                    {
                        "start": 0,
                        "end": 6
                    },
                    {
                        "start": 10,
                        "end": 16
                    }
                ]
            }
        },
        "document": {
            "SomeField": "needle & needle"
        }
    }
]
```
