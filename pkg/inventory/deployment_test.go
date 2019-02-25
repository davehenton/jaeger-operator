package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeploymentInventory(t *testing.T) {
	toCreate := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "to-create",
		},
	}
	toUpdate := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "to-update",
		},
	}
	toDelete := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "to-delete",
		},
	}

	existing := []appsv1.Deployment{toUpdate, toDelete}
	desired := []appsv1.Deployment{toUpdate, toCreate}

	inv := ForDeployments(existing, desired)
	assert.Len(t, inv.Create, 1)
	assert.Equal(t, "to-create", inv.Create[0].Name)

	assert.Len(t, inv.Update, 1)
	assert.Equal(t, "to-update", inv.Update[0].Name)

	assert.Len(t, inv.Delete, 1)
	assert.Equal(t, "to-delete", inv.Delete[0].Name)
}