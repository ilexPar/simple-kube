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

func TestServiceCreate(t *testing.T) {
	new := skres.Service{
		Name: "my-svc",
		Port: 80,
		Selector: map[string]string{
			"app": "nginx",
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Service](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Service().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.Service](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Service().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*api.Service)
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
			Service().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestServiceUpdate(t *testing.T) {
	old := &api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-svc",
			Namespace: "default",
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []api.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}
	new := skres.Service{
		Name: "my-svc",
		Port: 81, // changed
		Selector: map[string]string{
			"app": "nginx",
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Service](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Service().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.Service](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Service().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*api.Service)
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
			Service().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestServiceGet(t *testing.T) {
	kubeSvc := &api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-svc",
			Namespace: "default",
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []api.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}
	expected := skres.Service{
		Name: "my-svc",
		Port: 80,
		Selector: map[string]string{
			"app": "nginx",
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeSvc))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Service().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Service().
			Get("my-svc")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Service().
			Get("my-svc").
			DataHandler(func(res interface{}) error {
				svc := res.(*api.Service)
				svc.Spec.Ports[0].Port = 81 // override
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 81, result.Port)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Service().
			Get("my-svc").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestServiceList(t *testing.T) {
	dpl1 := &api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-svc",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []api.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}
	dpl2 := &api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-svc2",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []api.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}

	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(dpl1, dpl2))

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Service().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Service().
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
	dpl := &api.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-svc",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []api.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}
	t.Run("should return no errors when calling delete on an object", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(dpl)
		client := sk.NewClient(context.Background(), k8s)

		query := client.NamespacedQuery("default").
			Service().
			Delete("my-svc")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "services"))
	})
}
