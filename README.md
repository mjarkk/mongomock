# `mongomock` A mocking MongoDB database

mongomock is a simple to use library that mocks the mongodb library in memory.

This package is very handy for using in a testing envourment

## Quick start

```sh
go get -u github.com/mjarkk/mongomock
```

```go
package main

import (
    "log"

    "github.com/mjarkk/mongomock"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User {
    ID   primitive.ObjectID `bson:"_id" json:"id"`
    Name string             `bson:"username"`
}

func main() {
    db := mongomock.NewDB()
    collection := db.Collection("users")
    err := collection.InsertOne(User{
        ID:   primitive.NewObjectID(),
        Name: "test",
    })
    if err != nil {
        log.Fatal(err)
    }

    user := User{}
    err = collection.FindOne(&user, bson.M{"name": "test"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found user: %+v\n", user)
}
```
