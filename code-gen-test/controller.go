package main

import (
	"fmt"
	"time"

	"github.com/smarkm/k8s-crd/code-gen-test/pkg/gen/steward/clientset/versioned"
	"github.com/smarkm/k8s-crd/code-gen-test/pkg/gen/steward/clientset/versioned/scheme"
	informers "github.com/smarkm/k8s-crd/code-gen-test/pkg/gen/steward/informers/externalversions/steward/v1"
	listers "github.com/smarkm/k8s-crd/code-gen-test/pkg/gen/steward/listers/steward/v1"
	appsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	kgscheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Controller struct {
	kubeclientset    kubernetes.Interface
	stewardclientset versioned.Interface

	deployLister appslisters.DeploymentLister
	deploySynced cache.InformerSynced

	stewardLister listers.StewardLister
	stewardSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

//NewController RT
func NewController(
	kubeclientset kubernetes.Interface,
	stewardclientset versioned.Interface,
	deployInformer appsinformers.DeploymentInformer,
	stewardInformer informers.StewardInformer) *Controller {

	utilruntime.Must(scheme.AddToScheme(kgscheme.Scheme))
	klog.V(4).Info("Create event broadcaster")

	evenBroadcaster := record.NewBroadcaster()
	evenBroadcaster.StartLogging(klog.Infof)
	evenBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})

	recorder := evenBroadcaster.NewRecorder(kgscheme.Scheme, kcorev1.EventSource{Component: "Steward-Controller"})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		stewardclientset: stewardclientset,
		deployLister:     deployInformer.Lister(),
		deploySynced:     deployInformer.Informer().HasSynced,
		stewardLister:    stewardInformer.Lister(),
		stewardSynced:    stewardInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Steward"),
		recorder:         recorder,
	}

	klog.Info("create event handler")

	stewardInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueSteward,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueSteward(new)
		},
	})

	deployInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})
	return controller
}

func (c *Controller) enqueueSteward(obj interface{}) {
	klog.Info("Enqueue Object", obj)
}
func (c *Controller) handleObject(obj interface{}) {
	klog.Info("Handle Object", obj)
}

//Run RT
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Start Steward controller")

	klog.Info("Warting for informer cache to sync")
	if ok := cache.WaitForCacheSync(stopCh); !ok {
		return fmt.Errorf("failed wait for informer cache sync")
	}

	klog.Info("Start worker ....")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Worker started ...")
	<-stopCh
	klog.Info("Shutdown workers")
	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {

	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("except string but go %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(obj)
			return fmt.Errorf("error sync '%s':%s", key, err.Error())

		}
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true
}

func (c *Controller) syncHandler(key string) error {
	klog.Info("sync handler: %s", key)
	return nil
}
