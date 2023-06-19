package mongomock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	usersCollection := NewDB().Collection("users")

	Equal(t, "users", usersCollection.name)

	mockData := NewMockuser()

	// Insert dummy data
	err := usersCollection.Insert(mockData)
	NoError(t, err)
	documentsCount, err := usersCollection.Count(nil)
	NoError(t, err)
	Equal(t, uint64(1), documentsCount)

	// Create a new mockuser so we don't update the value behind the pointer
	// what might cause a fake positive when checking the updated data in the database
	newMockData := NewMockuser()
	realname := "John Doe"
	newMockData.Realname = &realname
	newMockData.ID = mockData.ID

	err = usersCollection.ReplaceFirstByID(mockData.ID, newMockData)
	NoError(t, err)

	// Check if the data in the database is actually replaced
	collectionData := usersCollection.documents
	Equal(t, 1, len(collectionData))

	firstItem := MockUser{}
	err = usersCollection.FindFirst(&firstItem, nil)
	NoError(t, err)
	NotNil(t, firstItem.Realname)
	Equal(t, realname, *firstItem.Realname)
}
