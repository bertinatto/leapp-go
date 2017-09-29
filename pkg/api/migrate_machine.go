package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

// FIXME: using a wrapper for testing, but the runner should be available via env variable
const RUNNER = "/home/fjb/src/snactor/runner_wrapper.sh"

// MigrateParams represents the parameters sent by the client
type MigrateParams struct {
	SourceHost       string            `json:"source_host,omitempty"`
	TargetHost       string            `json:"target_host,omitempty"`
	ContainerName    string            `json:"container_name,omitempty"`
	SourceUser       string            `json:"source_user,omitempty"`
	TargetUser       string            `json:"target_user,omitempty"`
	ExcludePaths     []string          `json:"excluded_paths"`
	TcpPorts         map[uint16]uint16 `json:"tcp_ports,omitempty"`
	ExcludedTcpPorts []uint16          `json:"excluded_tcp_ports,omitempty"`
	ForceCreate      bool              `json:"force_create,omitempty"`
	DisableStart     bool              `json:"disable_start,omitempty"`
	Debug            bool              `json:"debug,omitempty"`
}

// buildActorInput translates the data sent by the client into data that the actor can interpret
// This is ugly and nonintuitive,
func buildActorInput(p *MigrateParams) (string, error) {
	data := make(map[string]interface{})

	var sc bool
	if p.DisableStart == true {
		sc = false
	} else {
		sc = true
	}

	data["start_container"] = map[string]interface{}{"value": sc}
	data["excluded_paths"] = map[string]interface{}{"value": []string{}}
	data["excluded_tcp_ports"] = map[string]interface{}{"tcp": make(map[string]string)}
	data["force_create"] = map[string]interface{}{"value": false}
	data["source_host"] = map[string]interface{}{"value": p.SourceHost}
	data["source_user_name"] = map[string]interface{}{"value": p.SourceUser}
	data["target_host"] = map[string]interface{}{"value": p.TargetHost}
	data["target_user_name"] = map[string]interface{}{"value": p.TargetUser}
	data["tcp_ports_user_mapping"] = map[string]interface{}{"ports": []string{}}
	data["use_default_port_map"] = map[string]interface{}{"value": true}
	data["container_name"] = map[string]interface{}{"value": p.ContainerName}

	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(j), nil

}

// MigrateMachine handles the /migrate-machine endpoint
func MigrateMachine(request *http.Request) (interface{}, error) {
	var params MigrateParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	// Translate data sent by client into data that actor can read
	actor_data, err := buildActorInput(&params)
	if err != nil {
		return nil, err
	}

	// Call the actor runner passing data to its stdin
	c := executor.Command{
		CmdLine: strings.Split(RUNNER, " "),
		Stdin:   actor_data,
	}
	log.Println(c.CmdLine)
	log.Println(c.Stdin)
	r := c.Execute()

	return r, nil
}
