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

func TestCronJobCreate(t *testing.T) {
	new := skres.CronJob{
		Name: "my-cron",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa",
			},
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.CronJob](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				CronJob().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.CronJob](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				CronJob().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*batch.CronJob)
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
			CronJob().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestCronJobUpdate(t *testing.T) {
	old := &batch.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cron",
			Namespace: "default",
		},
		Spec: batch.CronJobSpec{
			Schedule: "*/5 * * * *",
			JobTemplate: batch.JobTemplateSpec{
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
			},
		},
	}
	new := skres.CronJob{
		Name:     "my-cron",
		Schedule: "*/5 * * * *",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa2",
			},
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.CronJob](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				CronJob().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.CronJob](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				CronJob().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*batch.CronJob)
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
			CronJob().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestCronJobGet(t *testing.T) {
	kubeCronJob := &batch.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cron",
			Namespace: "default",
		},
		Spec: batch.CronJobSpec{
			Schedule: "*/5 * * * *",
			JobTemplate: batch.JobTemplateSpec{
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
			},
		},
	}
	expected := skres.CronJob{
		Name:     "my-cron",
		Schedule: "*/5 * * * *",
		Containers: []skres.Container{
			{
				Name:  "main",
				Image: "sarasa",
			},
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeCronJob))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			CronJob().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			CronJob().
			Get("my-cron")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			CronJob().
			Get("my-cron").
			DataHandler(func(res interface{}) error {
				cron := res.(*batch.CronJob)
				cron.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = "overrided"
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "overrided", result.Containers[0].Image)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			CronJob().
			Get("my-cron").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestCronJobList(t *testing.T) {
	job1 := &batch.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cron",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: batch.CronJobSpec{
			Schedule: "*/5 * * * *",
			JobTemplate: batch.JobTemplateSpec{
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
			},
		},
	}
	job2 := &batch.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cron2",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Spec: batch.CronJobSpec{
			Schedule: "*/10 * * * *",
			JobTemplate: batch.JobTemplateSpec{
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
			},
		},
	}

	client := sk.NewClient(
		context.Background(),
		fake.NewSimpleClientset(job1, job2),
	)

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			CronJob().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			CronJob().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestCronJobDelete(t *testing.T) {
	job := &batch.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cron",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: batch.CronJobSpec{
			Schedule: "*/5 * * * *",
			JobTemplate: batch.JobTemplateSpec{
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
			},
		},
	}
	t.Run("should return no errors when calling delete on an object", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(job)
		client := sk.NewClient(context.Background(), k8s)

		query := client.NamespacedQuery("default").
			CronJob().
			Delete("my-cron")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "cronjobs"))
	})
}
