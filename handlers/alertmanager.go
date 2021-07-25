package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	conf "k8s-scale-incremental/config"
	k8s "k8s-scale-incremental/kubernetes"

	"github.com/prometheus/alertmanager/template"
)

// http post request handler and decode request body to json and return status 200
func AlertmanagerWebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	config := *conf.AppConfig

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := template.Data{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding json:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, alert := range data.Alerts {
		action := alert.Labels["action"]
		namespace := alert.Labels["namespace"]
		deployment := alert.Labels["deployment"]

		// if valid action, namespace and deployment, then scale
		if action != "scale-up" && action != "scale-down" {
			log.Error("Invalid action:", action)
			http.Error(w, "Invalid value for field action in body", http.StatusBadRequest)
		}

		if config.IfNamespaceExists(namespace) && config.IfDeploymentExists(namespace, deployment) {
			err := k8s.ScaleDeployment(namespace, deployment, action)
			if err != nil {
				log.Error("Error scaling deployment:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			log.Error("Namespace or deployment does not exist:", namespace, deployment)
			http.Error(w, "Namespace or Deployment not configured", http.StatusBadRequest)
		}
	}
	// Set response code to 200
	http.Error(w, "Successfully scaled the deployment", http.StatusOK)
}
