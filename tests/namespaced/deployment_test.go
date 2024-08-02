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
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeploymentCreate(t *testing.T) {
	new := skres.Deployment{
		Name: "my-deployment",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa",
			},
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Deployment](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Deployment().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.Deployment](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Deployment().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*apps.Deployment)
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
			Deployment().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestDeploymentUpdate(t *testing.T) {
	old := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
		},
		Spec: apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "sarasa",
						},
					},
				},
			},
		},
	}
	new := skres.Deployment{
		Name: "my-deployment",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa2",
			},
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Deployment](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Deployment().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.Deployment](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Deployment().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*apps.Deployment)
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
			Deployment().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestDeploymentGet(t *testing.T) {
	kubeDeployment := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
		},
		Spec: apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "sarasa",
						},
					},
				},
			},
		},
	}
	expected := skres.Deployment{
		Name: "my-deployment",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa",
			},
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeDeployment))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Deployment().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Deployment().
			Get("my-deployment")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Deployment().
			Get("my-deployment").
			DataHandler(func(res interface{}) error {
				deployment := res.(*apps.Deployment)
				deployment.Spec.Template.Spec.Containers[0].Image = "overrided"
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "overrided", result.Containers[0].Image)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Deployment().
			Get("my-deployment").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestDeploymentList(t *testing.T) {
	dpl1 := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
	dpl2 := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment2",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Spec: apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "httpd",
						},
					},
				},
			},
		},
	}

	client := sk.NewClient(
		context.Background(),
		fake.NewSimpleClientset(dpl1, dpl2),
	)

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Deployment().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Deployment().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestDeploymentDelete(t *testing.T) {
	dpl := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
	t.Run("should return no errors when calling delete on an object", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(dpl)
		client := sk.NewClient(context.Background(), k8s)

		query := client.NamespacedQuery("default").
			Deployment().
			Delete("my-deployment")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "deployments"))
	})
}
