package serverManagement

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

var action string

const (
	shutdownServer = "shutdownServer"
	restartServer  = "restartServer"
)

func ServeurAction(act string) {
	action = act
	// log.Println("ServerRestart")
	// cmd := exec.Command("./internal/serverManagement/restart.sh")
	// cmd.Run()
	// os.Exit(0)
}

func WaitServerOrder(serveur *http.Server) {
	for {
		if action != "" {
			switch action {
			case restartServer:
				log.Println("Serveur Restart")
				serveur.Shutdown(context.Background())
				action = ""
				os.Stdout.WriteString("1")
				return
			case shutdownServer:
				serveur.Shutdown(context.Background())
				log.Println("Serveur shutdown")
				os.Exit(0)
			}
		}
		time.Sleep(time.Second * 2)
	}
}
