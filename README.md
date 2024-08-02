# simple-kube

Library to simplify querying Kubernetes objects

# Usage

First select if you want to query cluster level objects or namespaced objects

```go
client := simplekube.NewClient(ctx, clientset)
clusterQuery := client.ClusterQuery()
namespacedQuery := client.NamespacedQuery(namespace)
```

Then you can choose one of the following actions:

- Get
- List
- Create
- Update
- Delete

Select any aditional options for your query and then call `Run()`

For  example:
```go
filter := map[string]string{"key": "value"}
cronjobs, err := namespacedQuery.CronJob().
    List().
    FilterByLabels(filter).
    Run()
```

# Advanced usage

Objects are simplified for basic use cases. But you can have access to the
raw kubernetes resource by providing a callback to `DataHandler`.

`Get` actions will execute the callback after getting Kubernetes API objects
and before loading them into library ones. While `Create` and `Update` execute
the callback after populating Kubernetes API objects but before calling
Kubernetes API. This way you should be able to tweak any aditional configuration
not yet available or supported by the library.

Example:

```go
// Create a CronJob with a custom termination grace period
cron, err := namespacedQuery.CronJob().Get("my-cron").
    DataHandler(func(obj interface{}) error {
        kc := res.(*batch.CronJob)
        grace := int64(10)
        kc.Spec.JobTemplate.Spec.Template.Spec.TerminationGracePeriodSeconds = grace
        return nil
    })
```