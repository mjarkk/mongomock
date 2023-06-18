package mongomock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestDelete(t *testing.T) {
	usersCollection := NewDB().Collection("users")

	mockData := NewMockuser()

	// Insert dummy data
	err := usersCollection.Insert(mockData)
	NoError(t, err)
	documentsCount, _ := usersCollection.Count(nil)
	Equal(t, uint64(1), documentsCount)

	// Delete entry and check if the collection is now empty
	usersCollection.DeleteByID(mockData.ID)
	documentsCount, _ = usersCollection.Count(nil)
	Equal(t, uint64(0), documentsCount)

	// Should result in no panics/errors if there is nothing to delete
	err = usersCollection.DeleteByID(mockData.ID)
	Equal(t, mongo.ErrNoDocuments, err)
	documentsCount, _ = usersCollection.Count(nil)
	Equal(t, uint64(0), documentsCount)
}
