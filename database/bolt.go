package database

import (
	"encoding/json"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

var bucketName = []byte("File2URL-Bot")

type BoltDatabase struct {
	database *bbolt.DB
}

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

func (db BoltDatabase) Close() error {
	return db.database.Close()
}

func (db BoltDatabase) Store(file File) (string, error) {
	uid := uuid.New()                                    // create an ID for this
	fileJson, _ := json.Marshal(file)                    // convert the file as json
	err := db.database.Update(func(tx *bbolt.Tx) error { // save in database
		return tx.Bucket(bucketName).Put(uid[:], fileJson)
	})
	return uid.String(), err
}

func (db BoltDatabase) Load(id string) (File, bool) {
	// Parse the UID
	uid, err := uuid.Parse(id)
	if err != nil { // invalid UID, just ignore this
		return File{}, false
	}
	// Get data from database
	var resultBytes []byte
	err = db.database.View(func(tx *bbolt.Tx) error {
		resultBytes = tx.Bucket(bucketName).Get(uid[:])
		return nil
	})
	if err != nil || resultBytes == nil { // db error or id does not exists
		return File{}, false
	}
	// Parse json
	var result File
	err = json.Unmarshal(resultBytes, &result)
	return result, err == nil
}
