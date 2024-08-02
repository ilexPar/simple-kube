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
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestJobCreate(t *testing.T) {
	new := skres.Job{
		Name: "my-job",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa",
			},
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Job](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Job().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.Job](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Job().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*batch.Job)
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
			Job().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestJobUpdate(t *testing.T) {
	old := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-job",
			Namespace: "default",
		},
		Spec: batch.JobSpec{
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
	new := skres.Job{
		Name: "my-job",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa2",
			},
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Job](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Job().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.Job](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Job().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*batch.Job)
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
			Job().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestJobGet(t *testing.T) {
	kubeDeployment := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-job",
			Namespace: "default",
		},
		Spec: batch.JobSpec{
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
	expected := skres.Job{
		Name: "my-job",
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
			Job().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Job().
			Get("my-job")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Job().
			Get("my-job").
			DataHandler(func(res interface{}) error {
				job := res.(*batch.Job)
				job.Spec.Template.Spec.Containers[0].Image = "overrided"
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "overrided", result.Containers[0].Image)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Job().
			Get("my-job").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestJobList(t *testing.T) {
	job1 := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-job",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: batch.JobSpec{
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
	job2 := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-job2",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Spec: batch.JobSpec{
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
		fake.NewSimpleClientset(job1, job2),
	)

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Job().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Job().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestJobDelete(t *testing.T) {
	job := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-job",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: batch.JobSpec{
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
		k8s := fake.NewSimpleClientset(job)
		client := sk.NewClient(context.Background(), k8s)

		query := client.NamespacedQuery("default").
			Job().
			Delete("my-job")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "jobs"))
	})
}
