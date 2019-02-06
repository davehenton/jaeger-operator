package inventory

import (
	appsv1 "k8s.io/api/apps/v1"
)

// Deployment represents the deployment inventory based on the current and desired states
type Deployment struct {
	Create []appsv1.Deployment
	Update []appsv1.Deployment
	Delete []appsv1.Deployment
}

// ForDeployments builds a new Deployment inventory based on the existing and desired states
func ForDeployments(existing []appsv1.Deployment, desired []appsv1.Deployment) Deployment {
	update := []appsv1.Deployment{}
	mcreate := deploymentMap(desired)
	mdelete := deploymentMap(existing)

	for k, v := range mcreate {
		if _, ok := mdelete[k]; ok {
			update = append(update, v)
			delete(mcreate, k)
			delete(mdelete, k)
		}
	}

	return Deployment{
		Create: deploymentList(mcreate),
		Update: update,
		Delete: deploymentList(mdelete),
	}
}

func deploymentMap(deps []appsv1.Deployment) map[string]appsv1.Deployment {
	m := map[string]appsv1.Deployment{}
	for _, d := range deps {
		m[d.Name] = d
	}
	return m
}

func deploymentList(m map[string]appsv1.Deployment) []appsv1.Deployment {
	l := []appsv1.Deployment{}
	for _, v := range m {
		l = append(l, v)
	}
	return l
}
