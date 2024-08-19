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
	scaling "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestHPACreate(t *testing.T) {
	new := skres.HPA{
		Name: "my-deployment",
		Min:  1,
		Max:  10,
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.HPA](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				HPA().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.HPA](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				HPA().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*scaling.HorizontalPodAutoscaler)
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
			HPA().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestHPAUpdate(t *testing.T) {
	minReplicas := int32(1)
	old := &scaling.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
		},
		Spec: scaling.HorizontalPodAutoscalerSpec{
			MinReplicas: &minReplicas,
			MaxReplicas: 10,
		},
	}
	new := skres.HPA{
		Name: "my-deployment",
		Min:  1,
		Max:  11,
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.HPA](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				HPA().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.HPA](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				HPA().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*scaling.HorizontalPodAutoscaler)
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
			HPA().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestHPAGet(t *testing.T) {
	minReplicas := int32(1)
	utilization := int32(70)
	kubeHPA := &scaling.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
		},
		Spec: scaling.HorizontalPodAutoscalerSpec{
			MinReplicas: &minReplicas,
			MaxReplicas: 10,
			Metrics: []scaling.MetricSpec{{
				Type: "Resource",
				Resource: &scaling.ResourceMetricSource{
					Name: "cpu",
					Target: scaling.MetricTarget{
						Type:               "Utilization",
						AverageUtilization: &utilization,
					},
				},
			}},
		},
	}
	expected := skres.HPA{
		Name: "my-deployment",
		Min:  1,
		Max:  10,
		Metrics: []skres.HPAMetric{{
			Type: "Resource",
			Resource: skres.HPAResourceMetric{
				Name:        "cpu",
				Type:        "Utilization",
				Utilization: 70,
			},
		}},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeHPA))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			HPA().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			HPA().
			Get("my-deployment")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			HPA().
			Get("my-deployment").
			DataHandler(func(res interface{}) error {
				deployment := res.(*scaling.HorizontalPodAutoscaler)
				deployment.Spec.MaxReplicas = 11
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 11, result.Max)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			HPA().
			Get("my-deployment").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestHPAList(t *testing.T) {
	minReplicas := int32(1)
	hpa1 := &scaling.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: scaling.HorizontalPodAutoscalerSpec{
			MinReplicas: &minReplicas,
			MaxReplicas: 10,
		},
	}
	hpa2 := &scaling.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment2",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Spec: scaling.HorizontalPodAutoscalerSpec{
			MinReplicas: &minReplicas,
			MaxReplicas: 10,
		},
	}

	client := sk.NewClient(
		context.Background(),
		fake.NewSimpleClientset(hpa1, hpa2),
	)

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			HPA().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			HPA().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestHPADelete(t *testing.T) {
	hpa := &scaling.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: scaling.HorizontalPodAutoscalerSpec{
			MaxReplicas: 10,
		},
	}
	t.Run("should return no errors when calling delete on an object", func(t *testing.T) {
		k8s := fake.NewSimpleClientset(hpa)
		client := sk.NewClient(context.Background(), k8s)

		query := client.NamespacedQuery("default").
			HPA().
			Delete("my-deployment")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "horizontalpodautoscalers"))
	})
}
