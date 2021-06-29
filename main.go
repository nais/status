package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type statusResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Team        string `json:"team"`
	Status      string `json:"status"`
	Dashboard   string `json:"dashboard"`
}

var (
	clusterName string
	bindAddr    string
)

func init() {
	flag.StringVar(&bindAddr, "bind-address", ":8080", "ip:port where http requests are served")
	flag.StringVar(&clusterName, "cluster-name", os.Getenv("NAIS_CLUSTER_NAME"), "which cluster we are running in")
	flag.Parse()
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := &statusResponse{
			Name:        fmt.Sprintf("nais-%s", clusterName),
			Description: fmt.Sprintf("NAIS status in cluster %s", clusterName),
			Team:        "aura",
			Status:      "ok",
			Dashboard:   fmt.Sprintf("https://grafana.nais.io/d/0nmGteKmz/nais-cluster?var-cluster=%s", clusterName),
		}
		if err := json.NewEncoder(w).Encode(status); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(fmt.Fprint(w, "uh oh"))
		}
	})

	fmt.Println("running @", bindAddr)
	go func() {
		log.Fatal((&http.Server{Addr: bindAddr}).ListenAndServe())
	}()

	<-interrupt

	fmt.Println("shutting down")
	fmt.Println((&http.Server{Addr: bindAddr}).Shutdown(context.Background()))
}
