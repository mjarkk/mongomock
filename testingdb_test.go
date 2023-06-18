package mongomock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUser struct {
	ID       primitive.ObjectID `bson:"_id" json:"id" description:"The unique id of the entry in the MongoDB ObjectId format, for more info see: https://docs.mongodb.com/manual/reference/method/ObjectId/"`
	Realname *string            `bson:"real_name,omitempty"`
	Username string
}

func NewMockuser() *MockUser {
	return &MockUser{
		ID:       primitive.NewObjectID(),
		Realname: nil,
		Username: "Piet",
	}
}

func TestNewDB(t *testing.T) {
	testDB := NewDB()
	NotNil(t, testDB)
}
