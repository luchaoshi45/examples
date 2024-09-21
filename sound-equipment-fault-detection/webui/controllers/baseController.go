package controllers

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// var kubeconfig *string

// func init() {
// 	// kubeconfig = flag.String("kubeconfig", "./config", "(optional) absolute path to the kubeconfig file")
// }

// This is a placeholder for your LED status
var ledStatus = true

// Handler function to get the LED status
func GetLEDStatus(c *gin.Context) {
	value := GetValue()

	if value == "0" {
		ledStatus = false
	} else if value == "1" {
		ledStatus = true
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isOn": ledStatus,
	})
}

func GetValue() string {
	flag.Parse()

	var config *rest.Config
	var err error

	// If a kubeconfig file is specified, use external config
	// if *kubeconfig != "" {
	// 	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// } else {
	// 	// Otherwise use the cluster internal configuration
	// 	config, err = rest.InClusterConfig()
	// }

	config, err = rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	// Create a dynamic client to access custom resources
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Define GVR (GroupVersionResource) to locate device resources
	deviceGVR := schema.GroupVersionResource{
		Group:    "devices.kubeedge.io",
		Version:  "v1beta1",
		Resource: "devices",
	}

	// Specify the namespace and device name
	namespace := "default"
	deviceName := "vled-instance-01"

	// Get device resources
	device, err := dynamicClient.Resource(deviceGVR).Namespace(namespace).Get(context.TODO(), deviceName, v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Print device information
	// fmt.Printf("Device Name: %s\n", device.GetName())
	// fmt.Printf("Namespace: %s\n", device.GetNamespace())

	// Get the device spec part
	// spec, found, err := unstructured.NestedMap(device.Object, "spec")
	// if err != nil || !found {
	//	 panic("Spec not found!")
	// }

	// fmt.Printf("Device Spec: %v\n", spec)

	// Get the device status part
	status, found, err := unstructured.NestedMap(device.Object, "status")
	if err != nil || !found {
		panic("Status not found!")
	}

	// fmt.Printf("Device Status: %v\n", status)

	// Get the twins list
	twins, found, err := unstructured.NestedSlice(status, "twins")
	if err != nil || !found {
		log.Fatalf("Twins not found or error occurred: %v", err)
	}

	// Iterate over the twins list
	for _, twin := range twins {
		if twinMap, ok := twin.(map[string]interface{}); ok {
			// Get the observedDesired field
			observedDesired, found, err := unstructured.NestedMap(twinMap, "observedDesired")
			if err != nil || !found {
				log.Printf("observedDesired not found or an error occurred: %v", err)
				continue
			}

			// Get the value
			value, found, err := unstructured.NestedString(observedDesired, "value")
			if err != nil || !found {
				log.Printf("value not found or error occurred: %v", err)
				continue
			} else {
				return value
			}
		}
	}
	return "-1"
}
