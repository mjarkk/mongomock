package update

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ApplyUpdate(document, update bson.M) (bson.M, error) {
	if update == nil {
		return document, nil
	}

	document, err := applySet(document, update["$set"])
	if err != nil {
		return nil, err
	}

	document, err = applyUnset(document, update["$unset"])
	if err != nil {
		return nil, err
	}

	return document, nil
}
