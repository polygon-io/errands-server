
package main

import (
	"os"
	"log"
	"os/signal"
)


var server *ErrandsServer
func main(){

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	server = NewErrandsServer()
	log.Println("listening for signals")
	for {
		select {
			case <-signals:
				// Logger.Info("main: done. exiting")
				log.Println("Exiting")
				server.kill()
				return
		}
	}

}






