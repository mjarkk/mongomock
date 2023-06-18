package mongomock

import (
	"errors"
	"reflect"

	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOne finds one document in the collection of placeInto
// The result can be filtered using filters
// The filters should work equal to MongoDB filters (https://docs.mongodb.com/manual/tutorial/query-documents/)
// tough this might miss features compared to mongoDB's filters
func (c *Collection) FindOne(placeInto any, filter bson.M) error {
	placeIntoReflection := reflect.ValueOf(placeInto)
	if placeIntoReflection.Kind() != reflect.Ptr {
		return errors.New("placeInto should be a pointer")
	}

	c.m.Lock()
	defer c.m.Unlock()

	for _, document := range c.documents {
		if match.Match(document.bson, filter) {
			err := bson.Unmarshal(document.bytes, placeInto)
			return err
		}
	}

	return mongo.ErrNoDocuments
}

// Find finds documents in the collection of the base
// The results can be filtered using filters
// The filters should work equal to MongoDB filters (https://docs.mongodb.com/manual/tutorial/query-documents/)
// tough this might miss features compared to mongoDB's filters
func (c *Collection) Find(results any, filter bson.M) error {
	c.m.Lock()
	defer c.m.Unlock()

	resultRefl := reflect.ValueOf(results)
	if resultRefl.Kind() != reflect.Ptr {
		return errors.New("requires pointer to slice as results argument")
	}

	resultRefl = resultRefl.Elem()
	if resultRefl.Kind() != reflect.Slice {
		return errors.New("requires pointer to slice as results argument")
	}

	resultsSliceContentType := resultRefl.Type().Elem()
	resultIsSliceOfPtrs := resultsSliceContentType.Kind() == reflect.Ptr
	if resultIsSliceOfPtrs {
		resultsSliceContentType = resultsSliceContentType.Elem()
	}

	for _, document := range c.documents {
		if !match.Match(document.bson, filter) {
			continue
		}

		newDocument := reflect.New(resultsSliceContentType)
		err := bson.Unmarshal(document.bytes, newDocument.Interface())
		if err != nil {
			return err
		}

		if resultIsSliceOfPtrs {
			resultRefl = reflect.Append(resultRefl, newDocument)
		} else {
			resultRefl = reflect.Append(resultRefl, newDocument.Elem())
		}
	}

	reflect.ValueOf(results).Elem().Set(resultRefl)

	return nil
}

// Cursor is a cursor for the testingdb implementing the db.Cursor
type Cursor struct {
	// should be set initially
	collection *Collection
	idx        int
	filter     bson.M
	// set after init
	document documentT
}

// Next returns the next item in the cursor
func (c *Cursor) Next() bool {
	c.collection.m.Lock()
	defer c.collection.m.Unlock()

	for c.idx < len(c.collection.documents) {
		c.document = c.collection.documents[c.idx]
		c.idx++
		if !match.Match(c.document.bson, c.filter) {
			continue
		}
		return true
	}

	return false
}

// Decode implements db.Cursor
func (c *Cursor) Decode(e any) error {
	eReflection := reflect.ValueOf(e)
	if eReflection.Kind() != reflect.Pointer {
		return errors.New("requires pointer as argument")
	}

	return bson.Unmarshal(c.document.bytes, e)
}

// FindCursor finds documents in the collection of the base
func (c *Collection) FindCursor(collectionName string, filter bson.M) (*Cursor, error) {

	c.m.Lock()
	cursor := &Cursor{
		collection: c,
		idx:        0,
		filter:     filter,
	}
	c.m.Unlock()

	return cursor, nil
}
