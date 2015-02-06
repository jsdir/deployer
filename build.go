package deployer

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/mholt/binding"
)

type Build struct {
	Service string
	Tag     string
}

func (build *Build) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&build.Service: binding.Field{
			Form:     "service",
			Required: true,
		},
		&build.Tag: binding.Field{
			Form:     "tag",
			Required: true,
		},
	}
}

func (build *Build) Equals(cmpBuild *Build) bool {
	log.Printf("tx %v", build)
	log.Printf("rx %v", cmpBuild)
	if build.Service != cmpBuild.Service {
		return false
	}
	if build.Tag != cmpBuild.Tag {
		return false
	}
	return true
}

func (build *Build) Exists(db *bolt.DB) (bool, error) {
	exists := false
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("builds"))
		if err != nil {
			return err
		}

		key := build.getKey()
		v := b.Get([]byte(key))
		exists = v != nil
		return nil
	})
	return exists, err
}

func (build *Build) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("builds"))
		if err != nil {
			return err
		}

		key := build.getKey()
		return b.Put([]byte(key), []byte("ts"))
	})
}

func (build *Build) getKey() string {
	return build.Service + ":" + build.Tag
}
