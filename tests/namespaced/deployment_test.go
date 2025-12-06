package namespaced_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	sk "github.com/ilexPar/simple-kube/pkg"
	skerr "github.com/ilexPar/simple-kube/pkg/errors"
	skres "github.com/ilexPar/simple-kube/pkg/namespaced/resources"
	kt "github.com/ilexPar/simple-kube/tests/k8sutil"
)

func TestDeploymentCreate(t *testing.T) {
	client := sk.Client{}
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
				Deployment().
				Create(new).
				DataHandler(func(deployment *apps.Deployment) error {
					assert.Equal(t, new.Name, deployment.Name)
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
			Deployment().
			Create(new).
			DataHandler(func(*apps.Deployment) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestDeploymentUpdate(t *testing.T) {
	client := sk.Client{}
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
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
			client.Config(context.Background(), k8s)

			query := client.InNamespace("default").
				Deployment().
				Update(new).
				DataHandler(func(deployment *apps.Deployment) error {
					assert.Equal(t, new.Name, deployment.Name)
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
			Deployment().
			Update(new).
			DataHandler(func(*apps.Deployment) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestDeploymentGet(t *testing.T) {
	client := sk.Client{}
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
	client.Config(context.Background(), fake.NewSimpleClientset(kubeDeployment))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.InNamespace("default").
			Deployment().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.InNamespace("default").
			Deployment().
			Get("my-deployment")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.InNamespace("default").
			Deployment().
			Get("my-deployment").
			DataHandler(func(deployment *apps.Deployment) error {
				deployment.Spec.Template.Spec.Containers[0].Image = "overrided"
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "overrided", result.Containers[0].Image)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.InNamespace("default").
			Deployment().
			Get("my-deployment").
			DataHandler(func(*apps.Deployment) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestDeploymentList(t *testing.T) {
	client := sk.Client{}
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

	client.Config(
		context.Background(),
		fake.NewSimpleClientset(dpl1, dpl2),
	)

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.InNamespace("default").
			Deployment().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.InNamespace("default").
			Deployment().
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
			Deployment().
			List().
			FilterByLabels(map[string]string{
				"app": "!nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestDeploymentDelete(t *testing.T) {
	client := sk.Client{}
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
		client.Config(context.Background(), k8s)

		query := client.InNamespace("default").
			Deployment().
			Delete("my-deployment")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "deployments"))
	})
}
