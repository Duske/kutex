package main

import (
	"errors"
	"fmt"
	cobra "github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	var namespace string
	var kubeconfig string

	rootCmd := &cobra.Command{
		Use:   "kutex",
		Short: "Kubernetes To External (Service)",
		Long:  `A command-line interface to replace a service with external one`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			svcName := args[0]
			ip := args[1]
			// use the current context in kubeconfig
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				panic(err.Error())
			}

			// create the clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				panic(errors.New("cannot connect to k8s"))
			}
			fmt.Printf("Using namespace %s\n", namespace)
			svc, err := clientset.CoreV1().Services(namespace).Get(svcName, metav1.GetOptions{})
			if err != nil {
				panic(fmt.Errorf("cannot connect to service %s \n", svcName))
			}
			var extSvc v1.Service
			svc.DeepCopyInto(&extSvc)
			extSvc.Spec.Selector = nil
			extSvc.ObjectMeta.ResourceVersion = ""
			extSvc.ResourceVersion = ""
			extSvc.ObjectMeta.Annotations = nil
			extSvc.Spec.Type = v1.ServiceTypeClusterIP

			var endpointPorts []v1.EndpointPort
			for _, port := range svc.Spec.Ports {
				endpointPorts = append(endpointPorts, v1.EndpointPort{
					Name:     port.Name,
					Port:     port.Port,
					Protocol: port.Protocol,
				})
			}

			fmt.Printf("Deleting service %s \n", svc.Name)
			err = clientset.CoreV1().Services(namespace).Delete(svc.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("Creating new service %s\n", svc.Name)
			_, err = clientset.CoreV1().Services(namespace).Create(&extSvc)
			if err != nil {
				panic(err.Error())
			}

			var endpoint = v1.Endpoints{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Endpoints",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: svcName,
				},
				Subsets: []v1.EndpointSubset{{Addresses: []v1.EndpointAddress{{IP: ip}}, Ports: endpointPorts}},
			}
			createdEndpoint, err := clientset.CoreV1().Endpoints(namespace).Create(&endpoint)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Endpoint %s created pointing to %s \n", createdEndpoint.Name, ip)
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "The namespace the current service is placed")
	if home := homeDir(); home != "" {
		rootCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		rootCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "absolute path to the kubeconfig file")
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
