package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

// actorStreamParams is a registry of executing actors.
var actorRunnerRegistry = NewActorRunner()

// logExecutorError logs the stderr of on executor.Result if verbose mode is enabled.
func logExecutorError(ctx context.Context, r *executor.Result) {
	if ctx.Value(CKey("Verbose")).(bool) {
		log.Printf("Actor stderr: %s\n", r.Stderr)
	}
}

// parseExecutorResult parses a *executor.Result and returns its stdout, and HTTP status code and and error, if any.
func parseExecutorResult(r *executor.Result) (interface{}, int, error) {
	if r.ExitCode != 0 {
		msg := fmt.Sprintf("actor execution failed with %d", r.ExitCode)
		return nil, http.StatusOK, NewApiError(nil, errActorExecution, msg)
	}

	if r.Stdout == "" {
		return nil, http.StatusOK, NewApiError(nil, errActorExecution, "actor didn't return any data")
	}

	var stdout interface{}
	if err := json.Unmarshal([]byte(r.Stdout), &stdout); err != nil {
		return nil, http.StatusOK, NewApiError(err, errActorExecution, "could not decode actor output")
	}
	return stdout, http.StatusOK, nil
}

// checkTargetParams verifies if s is valid and return and HTTP status code and and appropriate error.
func checkTaskStatus(s *ActorStatus) (int, error) {
	if s == nil {
		return http.StatusNotFound, NewApiError(nil, errTaskNotFound, "task not found")
	}

	if s.Result == nil {
		return http.StatusOK, NewApiError(nil, errTaskRunning, "task found, but there is no result yet")
	}
	return 0, nil
}

// respHandler is the final handler that builds the response to be sent to the clients.
func respHandler(fn respFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var r apiResult

		data, status, err := fn(rw, req)
		if err != nil {
			switch t := err.(type) {
			case apiError:
				r.Errors = append(r.Errors, t)
			default:
				http.Error(rw, "Internal error", http.StatusInternalServerError)
				return
			}
		} else {
			r.Data = data
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)

		err = json.NewEncoder(rw).Encode(r)
		if err != nil {
			log.Printf("could not encode response: %v\n", err)
		}
	}
}

// EndpointEntry represents an endpoint exposed by the daemon.
type EndpointEntry struct {
	Method      string
	Endpoint    string
	IsPrefix    bool
	NeedsStrip  bool
	HandlerFunc http.HandlerFunc
}

// GetEndpoints should return a slice of all endpoints that the daemon exposes.
func GetEndpoints() []EndpointEntry {
	return []EndpointEntry{
		{
			Method:      "POST",
			Endpoint:    "/migrate-machine",
			HandlerFunc: respHandler(migrateMachineStart),
		},
		{
			Method:      "GET",
			Endpoint:    "/migrate-machine/results/{id}",
			HandlerFunc: respHandler(migrateMachineResult),
		},
		{
			Method:      "POST",
			Endpoint:    "/port-inspect",
			HandlerFunc: respHandler(portInspect),
		},
		{
			Method:      "POST",
			Endpoint:    "/check-target",
			HandlerFunc: respHandler(checkTarget),
		},
		{
			Method:      "POST",
			Endpoint:    "/port-map",
			HandlerFunc: respHandler(portMap),
		},
		{
			Method:      "POST",
			Endpoint:    "/destroy-container",
			HandlerFunc: respHandler(destroyContainer),
		},
		{
			Method:      "GET",
			Endpoint:    "/doc",
			IsPrefix:    true,
			NeedsStrip:  true,
			HandlerFunc: http.FileServer(http.Dir("/usr/share/leapp/apidoc/")).ServeHTTP,
		},
	}
}
