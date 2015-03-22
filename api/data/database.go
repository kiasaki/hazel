package data

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

var ErrEntityNotFound = errors.New("Error: Entity not found")
var ErrEntityFound = errors.New("Error: Entity found")

type Database struct {
	DBFile string
	DB     *bolt.DB
}

type Model interface {
	Bucket() []byte
}

var registeredModels []Model

func NewDatabase(dbFile string) *Database {
	return &Database{DBFile: dbFile}
}

func (d *Database) Open() error {
	db, err := bolt.Open(d.DBFile, 0644, &bolt.Options{Timeout: time.Second * 10})
	if err != nil {
		return err
	}
	d.DB = db
	return d.createBuckets()
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) createBuckets() error {
	for _, model := range registeredModels {
		err := d.DB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(model.Bucket())
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func registerModel(model Model) {
	registeredModels = append(registeredModels, model)
}

func (d *Database) Get(key string, entity Model) error {
	return d.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(entity.Bucket())
		rawEntity := bucket.Get([]byte(key))
		if rawEntity == nil {
			return ErrEntityNotFound
		}
		return json.Unmarshal(rawEntity, entity)
	})
}

func (d *Database) Save(key string, entity Model) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(entity.Bucket())
		entityBytes, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), entityBytes)
	})
}
