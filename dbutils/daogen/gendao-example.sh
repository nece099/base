#!/bin/sh

echo '###### generate dao ... ...'

go run github.com/nece099/base/dbutils/daogen/daogen -p "github.com/PandUncle/TakeAwayApp/server/model"
gofmt -w ./model/model.go
