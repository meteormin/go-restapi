# go-restapi
go fiber with rest api

## install
```shell
go get -u github.com/miniyus/go-restapi
```

## usage

```go
package main

import "github.com/miniyus/go-restapi"

// new Handler
h = restapi.NewHandler[TestEntity, *TestReq, *TestRes](
    &TestReq{},
	// new service
	restapi.NewService[TestEntity, *TestReq, *TestRes](
        // new repository
        restapi.NewRepository[TestEntity](db, TestEntity{}),
        &TestRes{},
    ),
)

```