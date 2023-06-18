# `mongomock` A mocking MongoDB database

mongomock is a simple to use library that mocks features of a mongodb server but in memory and as a library with no external dependencies.

This package is very handy to use in a testing envourment

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
    Email string            `bson:"email"`
}

func main() {
    db := mongomock.NewDB()
    collection := db.Collection("users")
    err := collection.InsertOne(User{
        ID:   primitive.NewObjectID(),
        Name: "test",
        Email: "example@example.org",
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
    // After exit the database data is gone
}
```

## Supported methods

### `Count` - Count documents in a collection

```go
nr, err := db.Collection("users").Count(bson.M{})
```

### `DeleteByID` - Delete a document by ID

```go
err := db.Collection("users").DeleteByID(primitive.NewObjectID())
```

### `Dump` - Dump the database to std{out,err}

```go
// Dump the full database with all it's collections in json format to stdout
// Change the
panicResults false
err := db.Dump(panicResults)

// Dump a single collection to stdout
err := db.Collection("users").Dump(panicResults)
```

### `FindOne` - Find a single doucment in a collection

```go
user := User{}
err := db.Collection("users").FindOne(&user, bson.M{"email": "example@example.org"})
```

### `Find` - Find documents in a collection

```go
users := []User{}
err := db.Collection("users").Find(&users, bson.M{})
```

### `FindCursor` - Find documents in a collection using a cursor

```go
cursor, err := db.Collection("users").FindCursor(bson.M{})
if err != nil {
    log.Fatal(err)
}
for cursor.Next() {
    user := User{}
    err := cursor.Decode(&user)
    if err != nil {
        log.Fatal(err)
    }
}
```

### `Insert` - Insert a single document into a collection

```go
err := db.Collection("users").Insert(User{
    ID:    primitive.NewObjectID(),
    Name:  "test",
    Email: "example@example.org",
})
```
