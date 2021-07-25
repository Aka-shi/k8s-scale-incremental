package main

import (
	conf "k8s-scale-incremental/config"
	handlers "k8s-scale-incremental/handlers"
	k8s "k8s-scale-incremental/kubernetes"
	"net/http"
)

func main() {

	conf.LoadConfigFromFile()
	k8s.InitKubeClient()

	http.HandleFunc("/scale", handlers.AlertmanagerWebhookHandler)

	http.ListenAndServe(":8000", nil)
}
