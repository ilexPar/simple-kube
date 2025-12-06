package namespaced_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	sk "github.com/ilexPar/simple-kube/pkg"
	skerr "github.com/ilexPar/simple-kube/pkg/errors"
	skres "github.com/ilexPar/simple-kube/pkg/namespaced/resources"
	kt "github.com/ilexPar/simple-kube/tests/k8sutil"
)

func TestConfigMapCreate(t *testing.T) {
	client := sk.Client{}
	new := skres.ConfigMap{
		Name: "my-config",
		Data: map[string]string{
			"key": "value",
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.ConfigMap](t, kt.Create, func(k8s *fake.Clientset) {
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
				ConfigMap().
				Create(new).
				DataHandler(func(res *api.ConfigMap) error {
					assert.Equal(t, new.Name, res.Name)
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
		client.Config(context.Background(), k8s)

		query := client.InNamespace("default").
			ConfigMap().
			Create(new).
			DataHandler(func(res *api.ConfigMap) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestConfigMapUpdate(t *testing.T) {
	client := sk.Client{}
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
				ConfigMap().
				Update(new).
				DataHandler(func(res *api.ConfigMap) error {
					assert.Equal(t, new.Name, res.Name)
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
		client.Config(context.Background(), k8s)

		query := client.InNamespace("default").
			ConfigMap().
			Update(new).
			DataHandler(func(res *api.ConfigMap) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestConfigMapGet(t *testing.T) {
	client := sk.Client{}
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
	client.Config(context.Background(), fake.NewSimpleClientset(kubeSvc))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			Get("my-config")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			Get("my-config").
			DataHandler(func(res *api.ConfigMap) error {
				res.Data["key"] = "override" // override
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "override", result.Data["key"])
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			Get("my-config").
			DataHandler(func(res *api.ConfigMap) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestConfigMapList(t *testing.T) {
	client := sk.Client{}
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
	client.Config(context.Background(), fake.NewSimpleClientset(cmap1, cmap2))

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
	t.Run("should handle filtering labels by negating value", func(t *testing.T) {
		query := client.InNamespace("default").
			ConfigMap().
			List().
			FilterByLabels(map[string]string{
				"app": "!nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestConfigMapDelete(t *testing.T) {
	client := sk.Client{}
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
		client.Config(context.Background(), k8s)

		query := client.InNamespace("default").
			ConfigMap().
			Delete("my-config")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "configmaps"))
	})
}
