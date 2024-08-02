package namespaced_test

import (
	"context"
	"errors"
	"testing"

	sk "github.com/ilexPar/simple-kube/pkg"
	skerr "github.com/ilexPar/simple-kube/pkg/errors"
	skres "github.com/ilexPar/simple-kube/pkg/namespaced/resources"
	kt "github.com/ilexPar/simple-kube/tests/k8sutil"

	"github.com/stretchr/testify/assert"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestConfigMapCreate(t *testing.T) {
	new := skres.ConfigMap{
		Name: "my-config",
		Data: map[string]string{
			"key": "value",
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.ConfigMap](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				ConfigMap().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.ConfigMap](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				ConfigMap().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*api.ConfigMap)
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

		query := client.NamespacedQuery("default").
			ConfigMap().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestConfigMapUpdate(t *testing.T) {
	old := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"key": "value",
		},
	}
	new := skres.ConfigMap{
		Name: "my-config",
		Data: map[string]string{
			"key": "new value",
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.ConfigMap](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				ConfigMap().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.ConfigMap](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				ConfigMap().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*api.ConfigMap)
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

		query := client.NamespacedQuery("default").
			ConfigMap().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestConfigMapGet(t *testing.T) {
	kubeSvc := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"key": "value",
		},
	}
	expected := skres.ConfigMap{
		Name: "my-config",
		Data: map[string]string{
			"key": "value",
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeSvc))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			ConfigMap().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			ConfigMap().
			Get("my-config")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			ConfigMap().
			Get("my-config").
			DataHandler(func(res interface{}) error {
				cmap := res.(*api.ConfigMap)
				cmap.Data["key"] = "override" // override
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "override", result.Data["key"])
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			ConfigMap().
			Get("my-config").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestConfigMapList(t *testing.T) {
	cmap1 := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-config",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Data: map[string]string{
			"key": "value",
		},
	}
	cmap2 := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-config2",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Data: map[string]string{
			"key2": "different value",
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(cmap1, cmap2))

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			ConfigMap().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			ConfigMap().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestConfigMapDelete(t *testing.T) {
	dpl := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-config",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Data: map[string]string{
			"key": "value",
		},
	}
	t.Run("should return no errors when calling delete on an object", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(dpl)
		client := sk.NewClient(context.Background(), k8s)

		query := client.NamespacedQuery("default").
			ConfigMap().
			Delete("my-config")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "configmaps"))
	})
}
