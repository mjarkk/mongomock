package mongomock

import "go.mongodb.org/mongo-driver/bson"

// documentT describes a document within the collection
type documentT struct {
	bson  bson.M
	bytes []byte
}

func TryNewDocument(value any) (documentT, error) {
	encodedValue, err := bson.Marshal(value)
	if err != nil {
		return documentT{}, err
	}
	encodedValueMap := bson.M{}
	err = bson.Unmarshal(encodedValue, &encodedValueMap)
	if err != nil {
		return documentT{}, err
	}

	return documentT{
		bson:  encodedValueMap,
		bytes: encodedValue,
	}, nil
}
