package update

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ApplyUpdate(document, update bson.M) (bson.M, error) {
	if update == nil {
		return document, nil
	}

	var err error

	if value, ok := update["$set"]; ok {
		document, err = applySet(document, value)
		if err != nil {
			return nil, err
		}
	}

	if value, ok := update["$unset"]; ok {
		document, err = applyUnset(document, value)
		if err != nil {
			return nil, err
		}
	}

	return document, nil
}
