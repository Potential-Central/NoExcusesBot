package NoExcusesBot

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

//Interface for a database object
type DataObject interface {
	bucket() string
	primaryKey() string
}

//Writes an object to the database
func WriteObject(db *bolt.DB, object DataObject) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(object.bucket()))
		encoded, _ := json.Marshal(object)
		err := b.Put([]byte(object.primaryKey()), encoded)
		return err
	})
	return err
}

//Reads an object from the database
func ReadObject(db *bolt.DB, object DataObject) error {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(object.bucket()))
		v := b.Get([]byte(object.primaryKey()))
		err := json.Unmarshal(v, &object)
		return err
	})
	return err
}

//Returns all keys for a given bucket
func GetKeys(db *bolt.DB, bucket string) []string {
	var keys []string
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}

		return nil
	})
	return keys
}

//Creates database if not exists and sets up buckets
func CreateDB(buckets ...string) (*bolt.DB, error) {
	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	for _, bucket := range buckets {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
