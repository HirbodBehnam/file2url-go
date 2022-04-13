package database

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

var bucketName = []byte("File2URL-Bot")

type BoltDatabase struct {
	database *bbolt.DB
}

// NewBoltDatabase will open a bolt database in specified path
func NewBoltDatabase(path string) (BoltDatabase, error) {
	db, err := bbolt.Open(path, 666, nil)
	if err != nil {
		return BoltDatabase{}, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	return BoltDatabase{db}, err
}

// Close closes the database
func (db BoltDatabase) Close() error {
	return db.database.Close()
}

// Store will store the file in database overwriting any old entries.
// It will also create a random uuid and return its ID
func (db BoltDatabase) Store(file File) (string, error) {
	uid := uuid.New() // create an ID for this
	// Encode data
	var fileData bytes.Buffer
	err := gob.NewEncoder(&fileData).Encode(file)
	if err != nil {
		return "", fmt.Errorf("cannot gob the file: %w", err)
	}
	// Save in database
	err = db.database.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketName).Put(uid[:], fileData.Bytes())
	})
	return uid.String(), err
}

// Load will load an entry from database
func (db BoltDatabase) Load(id string) (File, bool) {
	// Parse the UID
	uid, err := uuid.Parse(id)
	if err != nil { // invalid UID, just ignore this
		return File{}, false
	}
	// Get data from database
	var file File
	err = db.database.View(func(tx *bbolt.Tx) error {
		resultBytes := tx.Bucket(bucketName).Get(uid[:])
		if resultBytes == nil { // id does not exists
			return errors.New("not found")
		}
		// Decode data with a buffer
		buffer := bytes.NewReader(resultBytes)
		return gob.NewDecoder(buffer).Decode(&file)
	})
	return file, err == nil
}
