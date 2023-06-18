package mongomock

import (
	"sync"
)

// Collection contains all the data for a collection
type Collection struct {
	m                     sync.Mutex
	underlayingCollection *TestConnection
	name                  string
	documents             []documentT
}
