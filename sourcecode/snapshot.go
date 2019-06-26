package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient kubernetes.Clientset

func snapshot(namespace string, outputfilename string) {

	test := K8sYaml{}
	var kubeconfig *string
	var n int
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace_array := strings.Split(namespace, ",")

	for n = 0; n < len(namespace_array); n++ {

		fmt.Printf("%d namespace: %s\n", n, namespace_array[n])

		//對deployment做處理
		deploy_array := KubectlGetDeployment(namespace_array[n])
		fmt.Println(len(deploy_array))
		for i := range deploy_array {
			if deploy_array[i] != "" && deploy_array[i] != "NAME" {
				fmt.Println("deployment name : " + deploy_array[i])
				imagename := GetDeploymentImage(clientSet, namespace_array[n], deploy_array[i])
				fmt.Println("Get deployment image name : " + imagename)
				gitbranch, modulename, moduletag := ImagenameSplitReturnTag(imagename)
				if IdentifyOpenfaas(namespace_array[n], deploy_array[i]) {
					(&test.Deployment).AddOpenfaasStruct(deploy_array[i], modulename, moduletag, gitbranch)
				} else {
					(&test.Deployment).AddK8sStruct(deploy_array[i], modulename, moduletag, gitbranch)
				}

			}
		}
		//對statefulset做處理
		statefulset_array := KubectlGetStefulset(namespace_array[n])
		fmt.Println(len(statefulset_array))
		for i := range statefulset_array {
			if statefulset_array[i] != "" && statefulset_array[i] != "NAME" {
				fmt.Println("statefulset name : " + statefulset_array[i])
				imagename := GetStatefulSetsImage(clientSet, namespace_array[n], statefulset_array[i])
				fmt.Println("Get deployment image name : " + imagename)
				gitbranch, modulename, moduletag := ImagenameSplitReturnTag(imagename)
				if IdentifyOpenfaas(namespace_array[n], statefulset_array[i]) {
					(&test.Deployment).AddOpenfaasStruct(statefulset_array[i], modulename, moduletag, gitbranch)
				} else {
					(&test.Deployment).AddK8sStruct(statefulset_array[i], modulename, moduletag, gitbranch)
				}

			}
		}
		//對daemonset做處理
		daemonset_array := KubectlGetDaemonset(namespace_array[n])
		fmt.Println(len(daemonset_array))
		for i := range daemonset_array {
			if daemonset_array[i] != "" && daemonset_array[i] != "NAME" {
				fmt.Println("daemonset name : " + daemonset_array[i])
				imagename := GetDaemonsetImage(clientSet, namespace_array[n], daemonset_array[i])
				fmt.Println("Get daemonset image name : " + imagename)
				gitbranch, modulename, moduletag := ImagenameSplitReturnTag(imagename)
				if IdentifyOpenfaas(namespace_array[n], daemonset_array[i]) {
					(&test.Deployment).AddOpenfaasStruct(daemonset_array[i], modulename, moduletag, gitbranch)
				} else {
					(&test.Deployment).AddK8sStruct(daemonset_array[i], modulename, moduletag, gitbranch)
				}

			}
		}

		//對cronjob做處理
		cronjob_array := KubectlGetCronJob(namespace_array[n])
		fmt.Println(len(cronjob_array))
		for i := range cronjob_array {
			if cronjob_array[i] != "" && cronjob_array[i] != "NAME" {
				fmt.Println("daemonset name : " + cronjob_array[i])
				imagename := GetCronjobImage(clientSet, namespace_array[n], cronjob_array[i])
				fmt.Println("Get daemonset image name : " + imagename)
				gitbranch, modulename, moduletag := ImagenameSplitReturnTag(imagename)
				if IdentifyOpenfaas(namespace_array[n], cronjob_array[i]) {
					(&test.Deployment).AddOpenfaasStruct(cronjob_array[i], modulename, moduletag, gitbranch)
				} else {
					(&test.Deployment).AddK8sStruct(cronjob_array[i], modulename, moduletag, gitbranch)
				}

			}
		}

	}
	d, err := yaml.Marshal(&test)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	WriteWithIoutil(outputfilename, string(d))
}

func IdentifyOpenfaas(i_namesapce string, i_deployment string) bool {
	var i_token bool
	i_cmd := "kubectl get deploy -l faas_function -n " + i_namesapce + " | grep " + i_deployment
	i_cmd_result := RunCommand(i_cmd)

	if strings.Contains(i_cmd_result, i_deployment) {
		i_token = true
	} else {
		i_token = false
	}

	return i_token
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func GetDeploymentImage(clientSet *kubernetes.Clientset, namespace string, deploymentName string) string {
	deployment, err := clientSet.AppsV1beta1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
	var getimage string
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("Deployment not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting deployment%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment\n")
		name := deployment.GetName()
		fmt.Println("name ->", name)
		containers := &deployment.Spec.Template.Spec.Containers
		//	found := false
		for i := range *containers {
			c := *containers
			getimage = c[i].Image
		}
		/*
							fmt.Println("Old image ->", c[i].Image)
							if c[i].Name == *appName {
								found = true
								fmt.Println("Old image ->", c[i].Image)
								fmt.Println("New image ->", *imageName)
								c[i].Image = *imageName
				}
			}
					if found == false {
						fmt.Println("The application container not exist in the deployment pods.")
						os.Exit(0)
					}
					_, err := clientset.AppsV1beta1().Deployments("default").Update(deployment)
					if err != nil {
						panic(err.Error())
					}*/
	}
	return getimage
}

func GetStatefulSetsImage(clientSet *kubernetes.Clientset, namespace string, statefulsetName string) string {
	statefulset, err := clientSet.AppsV1().StatefulSets(namespace).Get(statefulsetName, metav1.GetOptions{})
	var getimage string
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("Statefulset not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting statefulset%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found statefulset\n")
		name := statefulset.GetName()
		fmt.Println("name ->", name)
		containers := &statefulset.Spec.Template.Spec.Containers
		for i := range *containers {
			c := *containers
			getimage = c[i].Image

		}
	}
	return getimage
}

func GetDaemonsetImage(clientSet *kubernetes.Clientset, namespace string, daemonsetName string) string {
	daemonset, err := clientSet.ExtensionsV1beta1().DaemonSets(namespace).Get(daemonsetName, metav1.GetOptions{})
	var getimage string
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("Daemonset not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting daemonset%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found daemonset\n")
		name := daemonset.GetName()
		fmt.Println("name ->", name)
		containers := &daemonset.Spec.Template.Spec.Containers
		for i := range *containers {
			c := *containers
			getimage = c[i].Image

		}
	}
	return getimage
}

func GetCronjobImage(clientSet *kubernetes.Clientset, namespace string, cronjobName string) string {
	cronjob, err := clientSet.BatchV1beta1().CronJobs(namespace).Get(cronjobName, metav1.GetOptions{})
	//	daemonset, err := clientSet.ExtensionsV1beta1().DaemonSets(namespace).Get(daemonsetName, metav1.GetOptions{})
	var getimage string
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("CronJob not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting CronJob%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found CronJob\n")
		name := cronjob.GetName()
		fmt.Println("name ->", name)
		containers := &cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers
		for i := range *containers {
			c := *containers
			getimage = c[i].Image

		}
	}
	return getimage
}
