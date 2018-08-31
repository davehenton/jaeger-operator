package controller

import (
	"context"
	"strings"

	"github.com/jaegertracing/jaeger-operator/pkg/apis/io/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

// Controller knows what type of deployments to build based on a given spec
type Controller interface {
	Create() []sdk.Object
	Update() []sdk.Object
}

// NewController build a new controller object for the given spec
func NewController(ctx context.Context, jaeger *v1alpha1.Jaeger) Controller {
	// we need a name!
	if jaeger.ObjectMeta.Name == "" {
		logrus.Infof("This Jaeger instance was created without a name. Setting it to 'my-jaeger'")
		jaeger.ObjectMeta.Name = "my-jaeger"
	}

	// normalize the storage type
	if jaeger.Spec.Storage.Type == "" {
		logrus.Infof(
			"Storage type wasn't provided for the Jaeger instance '%v'. Falling back to 'memory'",
			jaeger.ObjectMeta.Name,
		)
		jaeger.Spec.Storage.Type = "memory"
	}

	if unknownStorage(jaeger.Spec.Storage.Type) {
		logrus.Infof(
			"The provided storage type for the Jaeger instance '%v' is unknown ('%v'). Falling back to 'memory'",
			jaeger.ObjectMeta.Name,
			jaeger.Spec.Storage.Type,
		)
		jaeger.Spec.Storage.Type = "memory"
	}

	// normalize the deployment strategy
	if strings.ToLower(jaeger.Spec.Strategy) != "production" {
		jaeger.Spec.Strategy = "all-in-one"
	}

	logrus.Debugf("Jaeger strategy: %s", jaeger.Spec.Strategy)
	if jaeger.Spec.Strategy == "all-in-one" {
		return newAllInOneController(ctx, jaeger)
	}

	// check for incompatible options
	if strings.ToLower(jaeger.Spec.Storage.Type) == "memory" {
		logrus.Warnf(
			"No suitable storage was provided for the Jaeger instance '%v'. "+
				"Falling back to all-in-one. Storage type: '%v'",
			jaeger.ObjectMeta.Name,
			jaeger.Spec.Storage.Type,
		)
		return newAllInOneController(ctx, jaeger)
	}

	return newProductionController(ctx, jaeger)
}

func unknownStorage(typ string) bool {
	known := []string{
		"memory",
		"kafka",
		"elasticsearch",
		"cassandra",
	}

	for _, k := range known {
		if strings.ToLower(typ) == k {
			return false
		}
	}

	return true
}
