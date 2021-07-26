package kubernetes

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	conf "k8s-scale-incremental/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeClient *kubernetes.Clientset

// Initialise Kubernetes Client
func InitKubeClient() {

	var kubeconfig *rest.Config

	// In production, we will use in-cluster config to communicate with api server. In development, we will use a kubeconfig file to connect to a local cluster.
	if val, e := os.LookupEnv("APP_ENV"); e {
		if val == "prod" {
			config, err := rest.InClusterConfig()
			if err != nil {
				log.Error("Unable to load in cluster config. Shutting down")
				os.Exit(1)
			}
			kubeconfig = config
		} else {
			log.Error("Invalid APP_ENV")
			os.Exit(1)
		}
	} else {
		var devconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		config, err := clientcmd.BuildConfigFromFlags("", *devconfig)
		if err != nil {
			log.Error("Unable to load kube config. Trying to load in cluster config")
		}
		kubeconfig = config
	}

	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Info("Successfully initialised kubernetes client")
	kubeClient = client
}

// A function which takes namespace and deployment name as arguments, reads the deployment from the Kubernetes API server and set replication count to 0 and update the deployment
func ScaleDeployment(namespace string, deploymentName string, action string) error {

	var config conf.Config = *conf.AppConfig

	deployment, err := kubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	currReplicas := *deployment.Spec.Replicas
	minReplicas := config.Namespaces[namespace].Deployment[deploymentName].MinReplicas
	maxReplicas := config.Namespaces[namespace].Deployment[deploymentName].MaxReplicas

	// Shouuld not scale below minReplicas or above maxReplicas
	if action == "scale-up" {
		delta := config.Namespaces[namespace].Deployment[deploymentName].ScaleUpBatchSize
		targetReplicas := int32(math.Min(float64(currReplicas+delta), float64(maxReplicas)))
		deployment.Spec.Replicas = &(targetReplicas)
	} else {
		delta := config.Namespaces[namespace].Deployment[deploymentName].ScaleDownBatchSize
		targetReplicas := int32(math.Max(float64(currReplicas-delta), float64(minReplicas)))
		deployment.Spec.Replicas = &(targetReplicas)
	}

	fmt.Println("I am going to set the target replicas to:", *deployment.Spec.Replicas)

	_, err = kubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return err
}
