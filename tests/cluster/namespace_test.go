package cluster_test

import (
	"context"
	"errors"
	"testing"

	sk "github.com/ilexPar/simple-kube/pkg"
	skres "github.com/ilexPar/simple-kube/pkg/cluster/resources"
	skerr "github.com/ilexPar/simple-kube/pkg/errors"
	kt "github.com/ilexPar/simple-kube/tests/k8sutil"

	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNamespaceCreate(t *testing.T) {
	new := skres.Namespace{
		Name: "some-ns",
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Namespace](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.ClusterQuery().
				Namespace().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.Namespace](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.ClusterQuery().
				Namespace().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*api.Namespace)
					assert.Equal(t, new.Name, obj.Name)
					assert.Equal(t, baseKubeActions, len(k8s.Actions()))
					hasCallbackRun = true
					return nil
				})
			err := query.Run()

			assert.Nil(t, err)
			assert.True(t, hasCallbackRun)
			assert.Equal(t, baseKubeActions+1, len(k8s.Actions()))
		})
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		k8s := fake.NewSimpleClientset()
		client := sk.NewClient(context.Background(), k8s)

		query := client.ClusterQuery().
			Namespace().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestNamespaceUpdate(t *testing.T) {
	old := &api.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "some-ns",
			Labels: map[string]string{
				"test": "test",
			},
		},
	}
	new := skres.Namespace{
		Name: "some-ns",
		Labels: map[string]string{
			"test": "new data",
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Namespace](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.ClusterQuery().
				Namespace().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.Namespace](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.ClusterQuery().
				Namespace().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*api.Namespace)
					assert.Equal(t, new.Name, obj.Name)
					assert.Equal(t, baseKubeActions, len(k8s.Actions()))
					hasCallbackRun = true
					return nil
				})
			err := query.Run()

			assert.Nil(t, err)
			assert.True(t, hasCallbackRun)
			assert.Equal(t, baseKubeActions+1, len(k8s.Actions()))
		}, old)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(old)
		client := sk.NewClient(context.Background(), k8s)

		query := client.ClusterQuery().
			Namespace().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestNamespaceGet(t *testing.T) {
	kubeSvc := &api.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-ns",
		},
	}
	expected := skres.Namespace{
		Name: "my-ns",
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeSvc))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.ClusterQuery().
			Namespace().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.ClusterQuery().
			Namespace().
			Get("my-ns")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		overrideLabels := map[string]string{
			"test": "test",
		}
		query := client.ClusterQuery().
			Namespace().
			Get("my-ns").
			DataHandler(func(res interface{}) error {
				svc := res.(*api.Namespace)
				svc.ObjectMeta.Labels = overrideLabels
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, overrideLabels, result.Labels)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.ClusterQuery().
			Namespace().
			Get("my-ns").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestNamespaceList(t *testing.T) {
	ns1 := &api.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-ns",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
	}
	ns2 := &api.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-ns2",
			Labels: map[string]string{
				"some": "label",
			},
		},
	}

	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(ns1, ns2))

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.ClusterQuery().
			Namespace().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.ClusterQuery().
			Namespace().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestServiceDelete(t *testing.T) {
	dpl := &api.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-ns",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
	}
	t.Run("should return no errors when calling delete on an object", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(dpl)
		client := sk.NewClient(context.Background(), k8s)

		query := client.ClusterQuery().
			Namespace().
			Delete("my-ns")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "namespaces"))
	})
}
