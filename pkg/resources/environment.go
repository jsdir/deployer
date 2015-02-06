package resources

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type EnvironmentType interface {
	Deploy(*Deploy, interface{}) error
}

type Environment struct {
	ReleaseId    int    `json:'release'`
	Updated      string `json:'updated'`
	DeployActive bool   `json:'active'`
}

func DeployToEnv(db *bolt.Db, name string, release *Release) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("environments"))
		if err != nil {
			return err
		}

		key := []byte(name)
		data := b.Get(key)
		if data == nil {
			// Create the environment.
			env := Environment{}
		} else {
			env := new(Environment)
			json.Unmarshal(data, &env)
		}

		// Start the deploy.
		prevReleaseId := env.ReleaseId

		// Update the environment within the same transation we loaded it from.
		env.ReleaseId = release.Id
		env.Updated = "now"
		env.DeployActive = true

		data, err = json.Marshal(env)
		if err != nil {
			return err
		}

		err = b.Put(key, data)
		if err != nil {
			return err
		}

		// Load the previous release for service comparison.
	})

	if err != nil {
		return err
	}
}

func GetEnvironment(db *bolt.Db, name string) (*Environment, error) {
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
	})

	if err != nil {
		return nil, err
	}

	return env, nil
}
