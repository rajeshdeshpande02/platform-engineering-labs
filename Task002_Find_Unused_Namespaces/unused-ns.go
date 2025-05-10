package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// get ns
	// check resources exist
	// if at least one resource exists, its used ns
	// if no resources exist, its unused ns
	clientSet, err := getClientSet()
	if err != nil {
		fmt.Println("Error connecting to Kubernetes cluster:", err)
		os.Exit(1)
	}
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error getting namespaces:", err)
		return
	}
	processNamespaces(namespaces, clientSet)

}

func processNamespaces(namespaces *v1.NamespaceList, clientSet *kubernetes.Clientset) {
	var wg sync.WaitGroup

	for _, ns := range namespaces.Items {
		wg.Add(1)
		go checkUnusedNs(ns.Name, clientSet, &wg)
	}
	wg.Wait()
}
func checkUnusedNs(nsName string, clientSet *kubernetes.Clientset, wg *sync.WaitGroup) {
	defer wg.Done()
	deploys, err := clientSet.AppsV1().Deployments(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments in namespace %s: %v\n", nsName, err)
		return
	}
	if len(deploys.Items) > 0 {
		return
	}

	daemons, err := clientSet.AppsV1().StatefulSets(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing StatefulSets in namespace %s: %v\n", nsName, err)
		return
	}
	if len(daemons.Items) > 0 {
		return
	}
	/*
		Continue Checking required resources
	*/
	fmt.Printf("Unused namespace: %s\n", nsName)
}

// getClientSet initializes and returns a Kubernetes clientset by loading the kubeconfig file.
func getClientSet() (*kubernetes.Clientset, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home dir: %w", err)
	}

	kubeConfigPath := filepath.Join(homeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error getting kube config: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating client set: %w", err)
	}
	return clientSet, nil
}
