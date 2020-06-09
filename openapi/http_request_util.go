package openapi

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func AsGameServerStatus(deployment *appsv1.Deployment) *GameContainerStatus {
	status := UNKNOWN
	conditionsLength := len(deployment.Status.Conditions)
	if conditionsLength != 0 {
		latestCondition := deployment.Status.Conditions[conditionsLength-1]
		switch latestCondition.Status {
		case v1.ConditionTrue:
			switch latestCondition.Type {
			case appsv1.DeploymentReplicaFailure:
				status = ERROR
			case appsv1.DeploymentAvailable:
				status = RUNNING
			case appsv1.DeploymentProgressing:
				status = RESTARTING
			}
		case v1.ConditionFalse:
			status = ERROR
		case v1.ConditionUnknown:
			status = UNKNOWN
		}

		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 0 {
			status = STOPPED
		}
	}

	return &GameContainerStatus{
		Id:     deployment.Name,
		Status: status,
	}
}
