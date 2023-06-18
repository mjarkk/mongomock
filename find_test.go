package mongomock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindOneWithoutFilters(t *testing.T) {
	usersCollection := NewDB().Collection("users")

	mockData := NewMockuser()

	err := usersCollection.Insert(mockData)
	NoError(t, err)

	foundResult := MockUser{}
	err = usersCollection.FindOne(&foundResult, bson.M{})
	NoError(t, err)
	Equal(t, mockData.ID, foundResult.ID)
}

func TestFindWithoutFilters(t *testing.T) {
	usersCollection := NewDB().Collection("users")

	mockData := NewMockuser()

	err := usersCollection.Insert(mockData)
	NoError(t, err)

	foundResults := []MockUser{}
	err = usersCollection.Find(&foundResults, bson.M{})
	NoError(t, err)
	Len(t, foundResults, 1)
	Equal(t, mockData.ID, foundResults[0].ID)

	foundResultsPtrs := []*MockUser{}
	err = usersCollection.Find(&foundResultsPtrs, bson.M{})
	NoError(t, err)
	Len(t, foundResultsPtrs, 1)
	Equal(t, mockData.ID, foundResultsPtrs[0].ID)
}
