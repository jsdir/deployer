package deployer

import (
	"github.com/boltdb/bolt"
	"github.com/mholt/binding"
)

type Environment interface {
	Deploy(release Release, config string)
}

type Release struct {
	Id       int
	Name     string
	Config   interface{}
	Services map[string]string
}

type Build struct {
	Service string
	Tag     string
}

func (build *Build) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&build.Service: binding.Field{
			Form:     "service",
			required: true,
		},
		&build.Tag: binding.Field{
			Form:     "tag",
			required: true,
		},
	}
}

func (build *Build) Save(db *bolt.DB) error {
	b, err := db.CreateBucketIfNotExists "builds"
	if err != nil {
		return err
	}

	key := build.Service + ":" + build.Tag
	err := b.Put([]byte(key), []byte("ts"))
	return err
}
