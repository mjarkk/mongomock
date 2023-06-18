package mongomock

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// Dump prints the full database contents in the console
// This can be used in tests to dump the contents of the database might something fail or to debug
//
// shouldPanic controls if the output is only printed or also should panic
func (c *TestConnection) Dump(shouldPanicResults bool) {
	data := map[string][]bson.M{}
	c.m.Lock()
	for _, collection := range c.collections {
		collection.m.Lock()

		documents := []bson.M{}
		for _, document := range collection.documents {
			documents = append(documents, document.bson)
		}
		data[collection.name] = documents

		collection.m.Unlock()
	}
	c.m.Unlock()

	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)
	if shouldPanicResults {
		panic(jsonString)
	} else {
		fmt.Println(jsonString)
	}
}

// DumpCollection prints a full database collection it's contents in the console
// This can be used in tests to dump the contents of the database might something fail or to debug
//
// shouldPanic controls if the output is only printed or also should panic
func (c *Collection) Dump(shouldPanicResults bool) {
	c.m.Lock()
	documents := []bson.M{}
	for _, document := range c.documents {
		documents = append(documents, document.bson)
	}
	c.m.Unlock()

	jsonBytes, err := json.MarshalIndent(documents, "", "    ")
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)
	if shouldPanicResults {
		panic(jsonString)
	} else {
		fmt.Println(jsonString)
	}
}
