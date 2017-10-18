package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leapp-to/leapp-go/pkg/db"
	"github.com/leapp-to/leapp-go/pkg/web"
)

var (
	defaultReadTimeout = 5 // seconds

	flagHelp    = flag.Bool("help", false, "show usage")
	flagListen  = flag.String("listen", "127.0.0.1:8000", "host:port to listen on.")
	flagTimeout = flag.Int64("timeout", 10, "time range in which daemon has to send a response to the client.")
	flagVerbose = flag.Bool("verbose", false, "enable verbose mode")
	//temporary solution to be able to fake function in test
	handleDb = EnsureConnectCreateDb
	dbConn   *db.Database
)

func main() {
	os.Exit(Main(nil))
}

//EnsureConnectCreateDb ensures that dbDir is created, connects to DB and create table
//temporary function
func EnsureConnectCreateDb() *db.Database {
	db.EnsureDbPath()
	conn := db.Connect()
	conn.CreateTable()
	return conn
}

// Main controls the startup and shutdown of the daemon lifecycle
func Main(up chan<- struct{}) int {
	flag.Parse()

	if *flagHelp {
		flag.Usage()
		return 0
	}

	// Parse options
	options := web.Options{
		ListenAddress: *flagListen,
		ReadTimeout:   time.Duration(defaultReadTimeout),
		WriteTimeout:  time.Duration(*flagTimeout),

		Verbose: *flagVerbose,
	}

	// Start HTTP server
	webHandler := web.New(&options)
	log.Printf("Starting leapp-daemon at %s\n", options.ListenAddress)
	go webHandler.Run()

	//temporary - calls EnsureConnectCreateDb
	dbConn = handleDb()
	if dbConn != nil {
		defer dbConn.Close()
	}

	// Handle shutdown under different conditions
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Println("Received SIGTERM. Shutting down...")
	case up <- struct{}{}:
		log.Println("Up channel unblocked. Shutting down...")
	case err := <-webHandler.ErrorCh():
		log.Printf("Error starting service: %v\n", err)
		return 1
	}

	return 0
}
