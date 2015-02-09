package resources

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/jsdir/deployer/pkg/names"
)

type Release struct {
	Id       int               `json:"id"`
	Name     string            `json:"name"`
	Config   interface{}       `json:"config"`
	Services map[string]string `json:"services"`
}

func NewRelease(db *bolt.DB, build *Build) (*Release, error) {
	// Create a release
	release := new(Release)
	release.Id = 0
	release.Name = names.NewRandomName("-")
	release.Config = false
	release.Services = make(map[string]string)

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("releases"))
		if err != nil {
			return err
		}

		// Extend the params of the previous release if it exists.
		c := b.Cursor()
		k, v := c.Last()
		if k != nil {
			lastRelease := new(Release)
			err := json.Unmarshal(v, &lastRelease)
			if err != nil {
				return err
			}

			release.Id = lastRelease.Id + 1
			release.Services = lastRelease.Services
		}

		// Set the new build tag.
		release.Services[build.Service] = build.Tag

		// Save the release.
		data, err := json.Marshal(&release)
		if err != nil {
			return err
		}

		return b.Put(convertIdToBytes(release.Id), data)
	})

	if err != nil {
		return nil, err
	}

	return release, nil
}

func GetRelease(db *bolt.DB, id int) (*Release, error) {
	var release Release

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("releases"))
		if err != nil {
			return err
		}

		data := b.Get(convertIdToBytes(id))
		return json.Unmarshal(data, &release)
	})

	if err != nil {
		return nil, err
	}

	return &release, nil
}

func (r *Release) Deploy(db *bolt.DB, dest string, envConfig interface{}, envType EnvironmentType) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("environments"))
		if err != nil {
			return err
		}

		// Load the destination environment.
		key := []byte(dest)
		data := b.Get(key)
		env := new(Environment)
		if data != nil {
			err = json.Unmarshal(data, &env)

			if err != nil {
				return err
			}
		}

		// Change and save the environment data.
		lastReleaseId := env.ReleaseId

		// Update the environment in the same transation it was loaded from.
		env.ReleaseId = r.Id
		env.Updated = "now"
		env.DeployActive = true

		newData, err := json.Marshal(env)
		if err != nil {
			return err
		}

		err = b.Put(key, newData)
		if err != nil {
			return err
		}

		// Begin the deploy by first loading the last release for service
		// comparision.
		// TODO: handle errors in deploy

		releaseBucket := tx.Bucket([]byte("releases"))
		releaseData := releaseBucket.Get(convertIdToBytes(lastReleaseId))

		if releaseData == nil {
			return errors.New("could not load last release")
		}

		lastRelease := new(Release)
		err = json.Unmarshal(releaseData, &lastRelease)
		if err != nil {
			return err
		}

		changedServices := r.getChangedServices(lastRelease)

		if len(changedServices) == 0 {
			// TODO: no services changed: deploy successful
		}

		// Create deploy and give
		return envType.Deploy(&Deploy{
			Env:             env,
			LastRelease:     lastRelease,
			Release:         r,
			ChangedServices: changedServices,
			EnvConfig:       envConfig,
		})
	})
}

func (r *Release) getChangedServices(lastRelease *Release) []string {
	changedServices := []string{}
	for service, tag := range r.Services {
		if lastRelease.Services[service] != tag {
			changedServices = append(changedServices, tag)
		}
	}
	return changedServices
}

func convertIdToBytes(id int) []byte {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint16(result, uint16(id))
	return result
}
