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
	net "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestIngressCreate(t *testing.T) {
	new := skres.Ingress{
		Name:   "my-ingress",
		Domain: "example.com",
		Paths: []skres.IngressPathDef{
			{
				Path:    "/",
				Service: "my-service",
				Port:    80,
			},
		},
	}

	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Ingress](t, kt.Create, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Ingress().
				Create(new)
			err := query.Run()

			assert.Nil(t, err)
		})
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		kt.WithInformedClient[skres.Ingress](t, kt.Create, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Ingress().
				Create(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*net.Ingress)
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
			Ingress().
			Create(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))

	})
}

func TestIngressUpdate(t *testing.T) {
	old := &net.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ingress",
			Namespace: "default",
		},
		Spec: net.IngressSpec{
			Rules: []net.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: net.IngressRuleValue{
						HTTP: &net.HTTPIngressRuleValue{
							Paths: []net.HTTPIngressPath{
								{
									Path: "/",
									Backend: net.IngressBackend{
										Service: &net.IngressServiceBackend{
											Name: "my-service",
											Port: net.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	new := skres.Ingress{
		Name:   "my-ingress",
		Domain: "example.com",
		Paths: []skres.IngressPathDef{
			{
				Path:    "/",
				Service: "my-service",
				Port:    81,
			},
		},
	}
	t.Run("should success without errors", func(t *testing.T) {
		kt.WithInformedClient[skres.Ingress](t, kt.Update, func(k8s *fake.Clientset) {
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Ingress().
				Update(new)

			err := query.Run()

			assert.Nil(t, err)
		}, old)
	})
	t.Run("should run DataHandler callback before updating object", func(t *testing.T) {
		kt.WithInformedClient[skres.Ingress](t, kt.Update, func(k8s *fake.Clientset) {
			hasCallbackRun := false
			baseKubeActions := 2 // kube fake clients with informers starts with 2 actions
			client := sk.NewClient(context.Background(), k8s)

			query := client.NamespacedQuery("default").
				Ingress().
				Update(new).
				DataHandler(func(res interface{}) error {
					obj := res.(*net.Ingress)
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
			Ingress().
			Update(new).
			DataHandler(func(res interface{}) error {
				return errors.New("test error")
			})
		err := query.Run()

		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, 0, len(k8s.Actions()))
	})
}

func TestIngressGet(t *testing.T) {
	kubeSvc := &net.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ingress",
			Namespace: "default",
		},
		Spec: net.IngressSpec{
			Rules: []net.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: net.IngressRuleValue{
						HTTP: &net.HTTPIngressRuleValue{
							Paths: []net.HTTPIngressPath{
								{
									Path: "/",
									Backend: net.IngressBackend{
										Service: &net.IngressServiceBackend{
											Name: "my-service",
											Port: net.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	expected := skres.Ingress{
		Name:   "my-ingress",
		Domain: "example.com",
		Paths: []skres.IngressPathDef{
			{
				Path:    "/",
				Service: "my-service",
				Port:    80,
			},
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(kubeSvc))

	t.Run("should return custom error when not found", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Ingress().
			Get("not-found")
		_, err := query.Run()

		assert.Equal(t, skerr.ERROR_NOT_FOUND, err.Error())
	})
	t.Run("should return expected object", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Ingress().
			Get("my-ingress")
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("should run DataHandler callback", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Ingress().
			Get("my-ingress").
			DataHandler(func(res interface{}) error {
				svc := res.(*net.Ingress)
				svc.Spec.Rules[0].Host = "overriden.com"
				return nil
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, "overriden.com", result.Domain)
	})
	t.Run("should cancel execution on callback error", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Ingress().
			Get("my-ingress").
			DataHandler(func(interface{}) error {
				return errors.New("test error")
			})
		_, err := query.Run()

		assert.Equal(t, "test error", err.Error())
	})
}

func TestIngressList(t *testing.T) {
	igr1 := &net.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ingress",
			Namespace: "default",
			Labels: map[string]string{
				"app":  "nginx",
				"some": "label",
			},
		},
		Spec: net.IngressSpec{
			Rules: []net.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: net.IngressRuleValue{
						HTTP: &net.HTTPIngressRuleValue{
							Paths: []net.HTTPIngressPath{
								{
									Path: "/",
									Backend: net.IngressBackend{
										Service: &net.IngressServiceBackend{
											Name: "my-service",
											Port: net.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	igr2 := &net.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-site",
			Namespace: "default",
			Labels: map[string]string{
				"some": "label",
			},
		},
		Spec: net.IngressSpec{
			Rules: []net.IngressRule{
				{
					Host: "example3.com",
					IngressRuleValue: net.IngressRuleValue{
						HTTP: &net.HTTPIngressRuleValue{
							Paths: []net.HTTPIngressPath{
								{
									Path: "/",
									Backend: net.IngressBackend{
										Service: &net.IngressServiceBackend{
											Name: "my-other-service",
											Port: net.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	client := sk.NewClient(context.Background(), fake.NewSimpleClientset(igr1, igr2))

	t.Run("should return expected objects", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Ingress().
			List()
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
	t.Run("should filter by label", func(t *testing.T) {
		query := client.NamespacedQuery("default").
			Ingress().
			List().
			FilterByLabels(map[string]string{
				"app": "nginx",
			})
		result, err := query.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestIngressDelete(t *testing.T) {
	dpl := &net.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-ingress",
			Namespace: "default",
		},
		Spec: net.IngressSpec{
			Rules: []net.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: net.IngressRuleValue{
						HTTP: &net.HTTPIngressRuleValue{
							Paths: []net.HTTPIngressPath{
								{
									Path: "/",
									Backend: net.IngressBackend{
										Service: &net.IngressServiceBackend{
											Name: "my-service",
											Port: net.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
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
			Ingress().
			Delete("my-ingress")

		err := query.Run()

		assert.Nil(t, err)
		assert.True(t, k8s.Actions()[0].Matches("delete", "ingresses"))
	})
}
