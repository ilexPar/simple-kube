package roundtrip_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	scaling "k8s.io/api/autoscaling/v2"
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	net "k8s.io/api/networking/v1"

	"github.com/ilexPar/simple-kube/pkg/base"
	skcl "github.com/ilexPar/simple-kube/pkg/cluster/resources"
	"github.com/ilexPar/simple-kube/pkg/core"
	skns "github.com/ilexPar/simple-kube/pkg/namespaced/resources"
)

// roundTrip encodes a simplified resource into its native Kubernetes type and
// back, asserting the value survives. This is the only per-resource test code:
// it exercises the `sm` tag mappings, the part that is genuinely
// resource-specific. Fixtures use a single port / single rule on purpose — the
// single-element mappings are a deliberate design choice.
func roundTrip[T any, R base.KubernetesResources](t *testing.T, sample T) {
	t.Helper()
	codec := core.SMCodec[T, R]{}

	raw, err := codec.Encode(sample)
	require.NoError(t, err)

	var back T
	require.NoError(t, codec.Decode(raw, &back))
	assert.Equal(t, sample, back)
}

func TestRoundTrip(t *testing.T) {
	t.Run("Deployment", func(t *testing.T) {
		roundTrip[skns.Deployment, apps.Deployment](t, skns.Deployment{
			Name:       "web",
			Labels:     map[string]string{"app": "web"},
			Containers: []skns.Container{{Name: "main", Image: "nginx"}},
		})
	})
	t.Run("Service", func(t *testing.T) {
		roundTrip[skns.Service, v1.Service](t, skns.Service{
			Name:     "web",
			Port:     80,
			Selector: map[string]string{"app": "web"},
		})
	})
	t.Run("Job", func(t *testing.T) {
		roundTrip[skns.Job, batch.Job](t, skns.Job{
			Name:       "batch",
			Containers: []skns.Container{{Name: "main", Image: "busybox"}},
			Behaviour:  skns.JobBehaviour{RestartPolicy: v1.RestartPolicyNever, FinishedTTL: 60},
		})
	})
	t.Run("CronJob", func(t *testing.T) {
		roundTrip[skns.CronJob, batch.CronJob](t, skns.CronJob{
			Name:       "cron",
			Schedule:   "* * * * *",
			Containers: []skns.Container{{Name: "main", Image: "busybox"}},
			Behaviour:  skns.CronJobBehaviour{RestartPolicy: v1.RestartPolicyNever, SuccessHistory: 3},
		})
	})
	t.Run("ConfigMap", func(t *testing.T) {
		roundTrip[skns.ConfigMap, v1.ConfigMap](t, skns.ConfigMap{
			Name: "cfg",
			Data: map[string]string{"key": "value"},
		})
	})
	t.Run("Ingress", func(t *testing.T) {
		roundTrip[skns.Ingress, net.Ingress](t, skns.Ingress{
			Name:   "ing",
			Domain: "example.com",
			Paths: []skns.IngressPathDef{{
				Path:    "/",
				Type:    net.PathTypePrefix,
				Service: "web",
				Port:    80,
			}},
		})
	})
	t.Run("HPA", func(t *testing.T) {
		roundTrip[skns.HPA, scaling.HorizontalPodAutoscaler](t, skns.HPA{
			Name:   "hpa",
			Min:    1,
			Max:    5,
			Target: skns.HPATarget{Name: "web", Kind: "Deployment", APIVersion: "apps/v1"},
			Metrics: []skns.HPAMetric{{
				Type: scaling.ResourceMetricSourceType,
				Resource: skns.HPAResourceMetric{
					Name:        "cpu",
					Type:        scaling.UtilizationMetricType,
					Utilization: 80,
				},
			}},
		})
	})
	t.Run("Namespace", func(t *testing.T) {
		roundTrip[skcl.Namespace, v1.Namespace](t, skcl.Namespace{
			Name:   "team-a",
			Labels: map[string]string{"env": "prod"},
		})
	})
}
