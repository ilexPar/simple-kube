package crud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	scaling "k8s.io/api/autoscaling/v2"
	batch "k8s.io/api/batch/v1"
	api "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	sk "github.com/ilexPar/simple-kube/pkg"
	"github.com/ilexPar/simple-kube/pkg/base"
	skcl "github.com/ilexPar/simple-kube/pkg/cluster/resources"
	"github.com/ilexPar/simple-kube/pkg/core"
	skerr "github.com/ilexPar/simple-kube/pkg/errors"
	skns "github.com/ilexPar/simple-kube/pkg/namespaced/resources"
	kt "github.com/ilexPar/simple-kube/tests/k8sutil"
)

const (
	ns      = "default"
	objName = "obj"
)

// RunCRUD exercises the full action surface for a resource through the public
// fluent API. Because the CRUD flow now lives in one generic core, this single
// suite replaces the previous per-resource test files; adding a resource is one
// entry in TestCRUD plus a fixture.
func RunCRUD[T any, R base.KubernetesResources](
	t *testing.T,
	label string,
	scope func(*sk.Client) core.ScopeAction[T, R],
	sample T,
	namespaced bool,
) {
	// Seed the fake API with the native encoding of the sample, so Get/Update/
	// Delete operate on a well-formed object regardless of the resource shape.
	seed := encodeSeed[T, R](t, sample, namespaced)

	t.Run(label+"/create hits the API", func(t *testing.T) {
		kt.WithInformedClient[T](t, kt.Create, func(k8s *fake.Clientset) {
			c := sk.Client{}
			c.Config(context.Background(), k8s)
			assert.NoError(t, scope(&c).Create(sample).Run())
		})
	})
	t.Run(label+"/create runs DataHandler", func(t *testing.T) {
		c := sk.Client{}
		c.Config(context.Background(), fake.NewSimpleClientset())
		ran := false
		err := scope(&c).Create(sample).
			DataHandler(func(*R) error { ran = true; return nil }).
			Run()
		assert.NoError(t, err)
		assert.True(t, ran)
	})
	t.Run(label+"/create cancels on callback error", func(t *testing.T) {
		k8s := fake.NewSimpleClientset()
		c := sk.Client{}
		c.Config(context.Background(), k8s)
		err := scope(&c).Create(sample).
			DataHandler(func(*R) error { return errors.New("boom") }).
			Run()
		assert.EqualError(t, err, "boom")
		assert.Empty(t, k8s.Actions())
	})
	t.Run(label+"/get returns the object", func(t *testing.T) {
		c := sk.Client{}
		c.Config(context.Background(), fake.NewSimpleClientset(toRuntime(seed)))
		_, err := scope(&c).Get(objName).Run()
		assert.NoError(t, err)
	})
	t.Run(label+"/get normalizes not found", func(t *testing.T) {
		c := sk.Client{}
		c.Config(context.Background(), fake.NewSimpleClientset())
		_, err := scope(&c).Get("missing").Run()
		assert.EqualError(t, err, skerr.ERROR_NOT_FOUND)
	})
	t.Run(label+"/update hits the API", func(t *testing.T) {
		c := sk.Client{}
		c.Config(context.Background(), fake.NewSimpleClientset(toRuntime(seed)))
		assert.NoError(t, scope(&c).Update(sample).Run())
	})
	t.Run(label+"/delete hits the API", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(toRuntime(seed))
		c := sk.Client{}
		c.Config(context.Background(), k8s)
		assert.NoError(t, scope(&c).Delete(objName).Run())
	})
}

func encodeSeed[T any, R base.KubernetesResources](t *testing.T, sample T, namespaced bool) *R {
	t.Helper()
	raw, err := core.SMCodec[T, R]{}.Encode(sample)
	require.NoError(t, err)
	acc, err := meta.Accessor(raw)
	require.NoError(t, err)
	acc.SetName(objName)
	if namespaced {
		acc.SetNamespace(ns)
	}
	return raw
}

func toRuntime[R any](seed *R) k8sruntime.Object {
	return any(seed).(k8sruntime.Object)
}

func TestCRUD(t *testing.T) {
	RunCRUD(t, "deployment",
		func(c *sk.Client) core.ScopeAction[skns.Deployment, apps.Deployment] {
			return c.InNamespace(ns).Deployment()
		},
		skns.Deployment{
			Name:       objName,
			Containers: []skns.Container{{Name: "main", Image: "nginx"}},
		},
		true,
	)
	RunCRUD(t, "service",
		func(c *sk.Client) core.ScopeAction[skns.Service, api.Service] {
			return c.InNamespace(ns).Service()
		},
		skns.Service{Name: objName, Port: 80},
		true,
	)
	RunCRUD(t, "job",
		func(c *sk.Client) core.ScopeAction[skns.Job, batch.Job] {
			return c.InNamespace(ns).Job()
		},
		skns.Job{
			Name:       objName,
			Containers: []skns.Container{{Name: "main", Image: "busybox"}},
		},
		true,
	)
	RunCRUD(t, "cronjob",
		func(c *sk.Client) core.ScopeAction[skns.CronJob, batch.CronJob] {
			return c.InNamespace(ns).CronJob()
		},
		skns.CronJob{
			Name:       objName,
			Schedule:   "* * * * *",
			Containers: []skns.Container{{Name: "main", Image: "busybox"}},
		},
		true,
	)
	RunCRUD(t, "configmap",
		func(c *sk.Client) core.ScopeAction[skns.ConfigMap, api.ConfigMap] {
			return c.InNamespace(ns).ConfigMap()
		},
		skns.ConfigMap{Name: objName, Data: map[string]string{"k": "v"}},
		true,
	)
	RunCRUD(t, "ingress",
		func(c *sk.Client) core.ScopeAction[skns.Ingress, net.Ingress] {
			return c.InNamespace(ns).Ingress()
		},
		skns.Ingress{
			Name:   objName,
			Domain: "example.com",
			Paths: []skns.IngressPathDef{{
				Path: "/", Type: net.PathTypePrefix, Service: "web", Port: 80,
			}},
		},
		true,
	)
	RunCRUD(t, "hpa",
		func(c *sk.Client) core.ScopeAction[skns.HPA, scaling.HorizontalPodAutoscaler] {
			return c.InNamespace(ns).HPA()
		},
		skns.HPA{
			Name:   objName,
			Min:    1,
			Max:    5,
			Target: skns.HPATarget{Name: "web", Kind: "Deployment", APIVersion: "apps/v1"},
		},
		true,
	)
	RunCRUD(t, "namespace",
		func(c *sk.Client) core.ScopeAction[skcl.Namespace, api.Namespace] {
			return c.Namespace()
		},
		skcl.Namespace{Name: objName},
		false,
	)
}

func TestListFiltering(t *testing.T) {
	d1 := &apps.Deployment{ObjectMeta: metav1.ObjectMeta{
		Name: "a", Namespace: ns, Labels: map[string]string{"app": "nginx"},
	}}
	d2 := &apps.Deployment{ObjectMeta: metav1.ObjectMeta{
		Name: "b", Namespace: ns, Labels: map[string]string{"app": "httpd"},
	}}
	c := sk.Client{}
	c.Config(context.Background(), fake.NewSimpleClientset(d1, d2))

	t.Run("returns all objects", func(t *testing.T) {
		out, err := c.InNamespace(ns).Deployment().List().Run()
		assert.NoError(t, err)
		assert.Len(t, out, 2)
	})
	t.Run("filters by label", func(t *testing.T) {
		out, err := c.InNamespace(ns).Deployment().List().
			FilterByLabels(map[string]string{"app": "nginx"}).Run()
		assert.NoError(t, err)
		assert.Len(t, out, 1)
	})
	t.Run("filters by negated label", func(t *testing.T) {
		out, err := c.InNamespace(ns).Deployment().List().
			FilterByLabels(map[string]string{"app": "!nginx"}).Run()
		assert.NoError(t, err)
		assert.Len(t, out, 1)
	})
}
