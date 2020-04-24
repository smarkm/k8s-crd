package main

import (
	"flag"
	"time"

	stewardclientset "github.com/smarkm/k8s-crd/code-gen-test/pkg/gen/steward/clientset/versioned"
	stewardinformers "github.com/smarkm/k8s-crd/code-gen-test/pkg/gen/steward/informers/externalversions"
	kinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	klog.InitFlags(nil)
	var masterUrl string
	var kubeconfigPath string
	flag.StringVar(&masterUrl, "url", "", "master url")
	flag.StringVar(&kubeconfigPath, "config", "", "kube config")
	flag.Parse()

	stopCh := signals.SetupSignalHandler()
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		klog.Fatal("Error load kubeconfig: %s", err.Error())
	}
	kclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatal("Error build k8s client: %s", err.Error())
	}
	stewardclient, err := stewardclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatal("Error build steward client: %s", err.Error())
	}

	kInformerFactory := kinformers.NewSharedInformerFactory(kclient, time.Second*30)
	stewardInformerFactory := stewardinformers.NewSharedInformerFactory(stewardclient, time.Second*1)
	controller := NewController(kclient, stewardclient,
		kInformerFactory.Apps().V1().Deployments(),
		stewardInformerFactory.Oam().V1().Stewards())

	kInformerFactory.Start(stopCh)
	stewardInformerFactory.Start(stopCh)
	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatal("Error run controller : %s", err.Error())
	}
}
