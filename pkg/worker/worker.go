package worker

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/leapp-to/leapp-go/pkg/db"
	"github.com/leapp-to/leapp-go/pkg/executor"
)

var (
	Jobs    = make(chan Job, 100)
	Results = make(chan Result, 100)
)

type Job struct {
	ID  uint32
	Cmd *executor.Command
}

type Result struct {
	ID     uint32
	Result *executor.Result
}

func NewJob(c *executor.Command) uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Uint32()
	Jobs <- Job{ID: id, Cmd: c}
	return id
}

func Start(n int) {
	for i := 0; i < n; i++ {
		go collect()
		go process()
	}
}

func process() {
	for j := range Jobs {
		log.Printf("Scheduling job: %d\n", j.ID)
		r := j.Cmd.Execute()
		log.Printf("Finished job: %d\n", j.ID)
		Results <- Result{ID: j.ID, Result: r}
	}
}

func collect() {
	for res := range Results {
		r := res.Result
		log.Printf("Collecting job: %d\n", res.ID)
		if r.ExitCode != 0 {
			log.Printf("Error on job %d:\n%s\n", res.ID, r.Stderr)
			continue
		}

		// Decode result from actor and send it back to client
		var stdout interface{}
		if err := json.Unmarshal([]byte(r.Stdout), &stdout); err != nil {
			continue
		}

		log.Printf("Setting job %d:\n%s\n", res.ID, r.Stdout)
		db.Set(res.ID, r.Stdout)

		//if r.Stderr != "" {
		//log.Printf("Stderr of job %d:\n%s\n", res.ID, r.Stderr)
		//}

		//log.Printf("Stdout of job %d:\n%s\n", res.ID, r.Stdout)
	}
}
