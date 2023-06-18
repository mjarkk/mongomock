package mongomock

// Insert inserts an item into the database
// Implements db.Connection
func (c *Collection) Insert(documents ...any) error {
	c.m.Lock()
	defer c.m.Unlock()

	return c.UnsafeInsert(documents...)
}

// UnsafeInsert inserts data directly into the database without locking it
func (c *Collection) UnsafeInsert(documents ...any) error {
	if len(documents) == 0 {
		return nil
	}

	additiveDocuments := make([]documentT, len(documents))
	for idx, document := range documents {
		doc, err := TryNewDocument(document)
		if err != nil {
			return err
		}
		additiveDocuments[idx] = doc
	}

	c.documents = append(c.documents, additiveDocuments...)
	return nil
}
