// Library to simplify querying Kubernetes objects
//
// # Usage
//
// First create the client and configure it:
//
//	client := simplekube.Client{}
//	client.Config(ctx, clientset)
//
// This gives basic cluster level objects, for example creating a new namespace:
//
//	ns := skcl.Namespace{ Name: "my-namespace" }
//	err := client.Namespace().Create(ns).Run()
//
// Note: skcl being "github.com/ilexPar/simple-kube/pkg/cluster/resources"
//
// ## Resource Actions
//
// These are valid resource actions for each resource:
//
// - Get
// - List
// - Create
// - Update
// - Delete
//
// Select any aditional options for your query and then call `Run()`
//
// In this example we list Cronjobs in "my-namespace" filtered by a map of labels:
//
//	filter := map[string]string{"key": "value"}
//	cronjobs, err := client.InNamespace("my-ns").
//	  CronJob().
//	  List().
//	  FilterByLabels(filter).
//	  Run()
//
// # Advanced usage
//
// Objects are simplified for basic use cases. But you can have access to the
// raw kubernetes resource by providing a callback to `DataHandler`.
//
// `Get` actions will execute the callback after getting Kubernetes API objects
// and before loading them into library ones. While `Create` and `Update` execute
// the callback after populating Kubernetes API objects but before calling
// Kubernetes API. This way you should be able to tweak any aditional configuration
// not yet available or supported by the library.
//
// Example:
//
//	// Create a CronJob with a custom termination grace period
//	cron, err := client.InNamespace("my-namespace").
//	  CronJob().
//	  Get("my-cron").
//	  DataHandler(func(cronjob *batch.CronJob) error {
//	    grace := int64(10)
//	    cronjob.Spec.JobTemplate.Spec.Template.Spec.TerminationGracePeriodSeconds = grace
//	    return nil
//	  }).
//	  Run()
package simplekube

import (
	"context"

	"github.com/ilexPar/simple-kube/pkg/cluster"
	"github.com/ilexPar/simple-kube/pkg/namespaced"

	"k8s.io/client-go/kubernetes"
)

type clusterQuery = *cluster.Query

type Client struct {
	ctx    context.Context
	client kubernetes.Interface
	clusterQuery
}

func (c *Client) Config(ctx context.Context, client kubernetes.Interface) *Client {
	query := (&cluster.Query{}).Config(ctx, client)
	c.clusterQuery = query
	c.ctx = ctx
	c.client = client
	return c
}

func (c *Client) InNamespace(
	namespace string,
) namespaced.QueryNamespace {
	query := (&namespaced.Query{}).Config(c.ctx, c.client, namespace)
	return query
}
