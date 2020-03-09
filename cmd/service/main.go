package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/otesk/AiposHTML/pkg/handler"
)

func main() {
	logFile, err := os.Create("log.txt")
	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	h := &handler.Handler{C: client}

	mux := handler.BuildRouter(h)
	log.Println(http.ListenAndServe(":8000", mux))
}
