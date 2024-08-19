package k8sutil

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	skcl "github.com/ilexPar/simple-kube/pkg/cluster"
	skclres "github.com/ilexPar/simple-kube/pkg/cluster/resources"
	skns "github.com/ilexPar/simple-kube/pkg/namespaced"
	sknsres "github.com/ilexPar/simple-kube/pkg/namespaced/resources"

	"github.com/stretchr/testify/assert"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

type informedActions int

const (
	Create informedActions = iota
	Update
)

type k8sResources interface {
	skns.NamespacedResourcesConstrain | skcl.ClusterResourcesConstrain
}

func WithInformedClient[T k8sResources](
	t *testing.T,
	action informedActions,
	exec func(client *fake.Clientset),
	objects ...k8sruntime.Object,
) {
	kChan := make(chan any, 1)
	ctx, cancel := context.WithTimeoutCause(
		context.TODO(),
		5*time.Second,
		errors.New("timed out running client action"),
	)
	defer cancel()

	k8s := fakeClientWithInformer(
		ctx,
		func(i informers.SharedInformerFactory) []cache.SharedInformer {
			var resType T
			informer := resolveInformerByDef(resType, i)
			switch action {
			case Create:
				informer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
					AddFunc: func(obj interface{}) {
						kChan <- obj
					},
				})
			case Update:
				informer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
					UpdateFunc: func(oldObj, newObj interface{}) {
						assert.NotEqualf(
							t,
							newObj,
							oldObj,
							"old object and new object are the same",
						)
						kChan <- newObj
					},
				})
			}
			return []cache.SharedInformer{informer}

		},
		objects...,
	)

	exec(k8s)

	select {
	case <-kChan:
		break
	case <-time.After(time.Second * 3):
		t.Fatalf("test timed out")
	}
}

// fakeClientWithInformer creates a fake Kubernetes clientset that has
// an informer factory injected. The informer factory has informers
// registered based on the setInformers callback. The objects are used
// to seed the informers. This function is intended for testing code
// that uses Kubernetes informers.
func fakeClientWithInformer(
	ctx context.Context,
	setInformers func(i informers.SharedInformerFactory) []cache.SharedInformer,
	objects ...k8sruntime.Object,
) *fake.Clientset {
	watcher := make(chan struct{})
	client := fake.NewSimpleClientset(objects...)
	client.PrependWatchReactor(
		"*",
		func(action clienttesting.Action) (handled bool, ret watch.Interface, err error) {
			gvr := action.GetResource()
			ns := action.GetNamespace()
			watch, err := client.Tracker().Watch(gvr, ns)
			if err != nil {
				return false, nil, err
			}
			close(watcher)
			return true, watch, nil
		},
	)

	// We will create an informer that writes added pods to a channel.
	informers := informers.NewSharedInformerFactory(client, 0)

	settedInformers := setInformers(informers)

	// // Make sure informers are running.
	informers.Start(ctx.Done())
	// This is not required in tests, but it serves as a proof-of-concept by
	// ensuring that the informer goroutine have warmed up and called List before
	//
	//
	// we send any events to it.
	for _, i := range settedInformers {
		cache.WaitForCacheSync(ctx.Done(), i.HasSynced)
	}
	// Any writes to the client
	// after the informer's initial LIST and before the informer establishing the
	// watcher will be missed by the informer. Therefore we wait until the watcher
	// starts.
	<-watcher

	return client
}

// resolveInformerByDef returns the SharedInformer for the given definition.
// It panics if no case is provided for the definition type.
func resolveInformerByDef(def interface{}, i informers.SharedInformerFactory) cache.SharedInformer {
	switch def.(type) {
	case sknsres.Deployment:
		return i.Apps().V1().Deployments().Informer()
	case sknsres.ConfigMap:
		return i.Core().V1().ConfigMaps().Informer()
	case skclres.Namespace:
		return i.Core().V1().Namespaces().Informer()
	case sknsres.Service:
		return i.Core().V1().Services().Informer()
	case sknsres.Job:
		return i.Batch().V1().Jobs().Informer()
	case sknsres.CronJob:
		return i.Batch().V1().CronJobs().Informer()
	case sknsres.Ingress:
		return i.Networking().V1().Ingresses().Informer()
	case sknsres.HPA:
		return i.Autoscaling().V2().HorizontalPodAutoscalers().Informer()
	default:
		t := reflect.ValueOf(def).Type().Name()
		err := fmt.Sprintf("no case provided for %s", t)
		panic(err)
	}
}
