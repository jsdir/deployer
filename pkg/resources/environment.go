package resources

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type EnvironmentType interface {
	Deploy(*Deploy) error
}

type Environment struct {
	ReleaseId    int    `json:'release'`
	Updated      string `json:'updated'`
	DeployActive bool   `json:'active'`
}

func GetEnvironment(db *bolt.DB, name string) (*Environment, error) {
	var env *Environment
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("environments"))
		if err != nil {
			return err
		}

		key := []byte(name)
		data := b.Get(key)
		if data != nil {
			json.Unmarshal(data, env)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return env, nil
}
