package e2e

import (
	"testing"
	"time"

	"github.com/jaegertracing/jaeger-operator/pkg/apis/io/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	retryInterval = time.Second * 5
	timeout       = time.Minute * 5
)

func TestJaeger(t *testing.T) {
	jaegerList := &v1alpha1.JaegerList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Jaeger",
			APIVersion: "io.jaegertracing/v1alpha1",
		},
	}

	err := framework.AddToFrameworkScheme(v1alpha1.AddToScheme, jaegerList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	t.Run("jaeger-group", func(t *testing.T) {
		// t.Run("my-jaeger", JaegerAllInOne)
		// t.Run("my-other-jaeger", JaegerAllInOne)

		// t.Run("simplest", SimplestJaeger)
		t.Run("simple-prod", SimpleProd)
	})
}

func prepare(t *testing.T) framework.TestCtx {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	err := ctx.InitializeClusterResources()
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for memcached-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "jaeger-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	return ctx
}
