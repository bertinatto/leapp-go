package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/leapp-to/leapp-go/pkg/db"
)

func captureLog(f func()) string {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestStarts(t *testing.T) {
	up := make(chan struct{})
	//temporary - skip db.EnsureDbPath, db.Connect db.CreateTable
	handleDb = func() *db.Database { return nil }

	go Main(up)
	select {
	case <-up:
		t.Log("server is up")
	case <-time.After(5 * time.Second):
		t.Errorf("server took too long to start")
	}
}

func TestShutdown(t *testing.T) {
	//temporary - skip db.EnsureDbPath, db.Connect db.CreateTable
	handleDb = func() *db.Database { return nil }
	// This goroutine should start correctly
	go Main(nil)

	// This should return 1
	if Main(nil) != 1 {
		t.Errorf("server should not have started")
	}

	// Check listen error
	o := captureLog(func() {
		Main(nil)
	})

	if !strings.Contains(o, "address already in use") {
		t.Errorf("did not catch bind error")
	}

}
