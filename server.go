package deployer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jsdir/deployer/pkg/resources"

	"github.com/mholt/binding"
)

type ContextHandler func(*Context, http.ResponseWriter, *http.Request) (int, error)

func CreateServer(ctx *Context) {
	http.HandleFunc("/builds", makeHandler(ctx, buildHandler))
	http.HandleFunc("/releases", makeHandler(ctx, releaseHandler))
	http.HandleFunc("/deploys", makeHandler(ctx, deployHandler))

	log.Printf(fmt.Sprintf("Listening on port %v", ctx.Config.Port))
	http.ListenAndServe(fmt.Sprintf(":%v", ctx.Config.Port), nil)
}

func makeHandler(ctx *Context, fn ContextHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		status, err := fn(ctx, rw, req)

		// Send errors through the api.
		if err != nil {
			log.Printf("HTTP %d: %q", status, err)
			switch status {
			case http.StatusNotFound:
				http.NotFound(rw, req)
			case http.StatusInternalServerError:
				http.Error(rw, http.StatusText(status), status)
			default:
				http.Error(rw, http.StatusText(status), status)
			}
		}
	}
}

func buildHandler(ctx *Context, rw http.ResponseWriter, req *http.Request) (int, error) {
	// Create new build
	build := new(resources.Build)
	errs := binding.Bind(req, build)
	if errs != nil {
		return http.StatusBadRequest, errs
	}

	err := build.Save(ctx.Db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Notify any pending releases that their dependent build may now exist.
	ctx.NewBuilds.Write(build)
	return 201, nil
}

func releaseHandler(ctx *Context, rw http.ResponseWriter, req *http.Request) (int, error) {
	// Create new build for comparison
	build := new(resources.Build)
	errs := binding.Bind(req, build)
	if errs != nil {
		return http.StatusBadRequest, errs
	}

	/*
		exists, err := build.Exists(ctx.Db)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// If the build does not already exist, wait for it to be created.
		//var newBuild *Build
			if !exists {
				listenerId, newBuilds := ctx.NewBuilds.Listen()
				for {
					newBuild := <-newBuilds
					if build.Equals(newBuild.(*resources.Build)) {
						ctx.NewBuilds.Unregister(listenerId)
						break
					}
				}
			}
	*/

	// At this point, since the build exists, it is now safe to create and
	// deploy a release.
	release, err := resources.NewRelease(ctx.Db, build)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Respond with the release object.
	data, err := json.Marshal(release)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
	return 201, nil
}

func deployHandler(ctx *Context, rw http.ResponseWriter, req *http.Request) (int, error) {
	// Create a deploy request
	deploy := new(resources.DeployRequest)
	errs := binding.Bind(req, deploy)
	if errs != nil {
		return http.StatusBadRequest, errs
	}

	// Get the environment config
	envConfig, exists := ctx.Config.Environments[deploy.Dest]
	if !exists {
		return http.StatusBadRequest, errors.New("Invalid destination environment")
	}

	// Check if src is valid. First, we'll assume that src is a release id.
	var release *resources.Release
	id, err := strconv.Atoi(deploy.Src)
	if err == nil {
		release, err = resources.GetRelease(ctx.Db, id)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	if release == nil {
		// Since no release with the given id exists, it must belong to an
		// environment.
		env, err := resources.GetEnvironment(ctx.Db, deploy.Src)
		if err != nil {
			return http.StatusBadRequest, err
		}

		if env == nil {
			return http.StatusBadRequest, errors.New("Invalid source environment")
		}

		release, err = resources.GetRelease(ctx.Db, env.ReleaseId)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	// Get the environment type
	envType := GetEnvironmentType(envConfig["type"].(string))
	if envType == nil {
		return http.StatusBadRequest, errors.New("Invalid environment type")
	}

	// Deploy to destination environment.
	err = release.Deploy(ctx.Db, deploy.Dest, envConfig, envType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 201, nil
}
