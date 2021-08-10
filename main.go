package main

import (
	"log"
	"os"

	"github.com/navidshaikh/test-webhook/pkg/apis/ghwebhook"
)

func main() {
	f, err := os.OpenFile("/tmp/tkg5360.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	ghwebhook.Listen()
}
