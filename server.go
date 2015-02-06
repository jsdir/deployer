package deployer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mholt/binding"
)

type ServerError struct {
	err  error
	code int
}

type ContextHandler func(*Context, http.ResponseWriter, *http.Request) (int, error)

func CreateServer(ctx *Context) {
	http.HandleFunc("/builds", makeHandler(ctx, buildHandler))
	http.HandleFunc("/releases", makeHandler(ctx, releaseHandler))
	// http.HandleFunc("/deploys", makeHandler(ctx, deployHandler))

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
	build := new(Build)
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
	build := new(Build)
	errs := binding.Bind(req, build)
	if errs != nil {
		return http.StatusBadRequest, errs
	}

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
			if build.Equals(newBuild.(*Build)) {
				ctx.NewBuilds.Unregister(listenerId)
				break
			}
		}
	}

	// At this point, since the build exists, it is now safe to create and
	// deploy a release.
	release, err := NewRelease(ctx.Db, build)
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

/*
func deployHandler(ctx *Context, rw http.ResponseWriter, req *http.Request) {
	// Ensure both exist
	//req.from req.to
	// ensure to (if integer) is valid release
	// else if string, ensure valid environment
	// ensure to is valid environment
	// set environment deploys...
}
*/
