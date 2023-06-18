package mongomock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	testingDB := NewDB()
	usersCollection := testingDB.Collection("users")

	mockData := NewMockuser()

	err := usersCollection.Insert(mockData)
	NoError(t, err)

	v, ok := testingDB.collections["users"]
	True(t, ok)
	NotNil(t, v)
	Equal(t, "users", v.name)
	Len(t, v.documents, 1)
}
