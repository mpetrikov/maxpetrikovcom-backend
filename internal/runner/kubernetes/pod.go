package kubernetes

import (
	"fmt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
	k8scorev1 "k8s.io/api/core/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	labelApp          = "app.kubernetes.io/name"
	labelPartOf       = "app.kubernetes.io/part-of"
	labelLabSessionID = "maxpetrikov.com/lab-session-id"
	labelLabID        = "maxpetrikov.com/lab-id"
	labelUserID       = "maxpetrikov.com/user-id"

	labelValueApp    = "maxpetrikov-lab"
	labelValuePartOf = "maxpetrikovcom"
)

func makePodName(
	session labsession.Session,
) string {
	return fmt.Sprintf(
		"lab-%s",
		session.ID.String(),
	)
}

func makePodLabels(
	session labsession.Session,
) map[string]string {
	return map[string]string{
		labelApp:          labelValueApp,
		labelPartOf:       labelValuePartOf,
		labelLabSessionID: session.ID.String(),
		labelLabID:        session.LabID.String(),
		labelUserID:       session.UserID.String(),
	}
}

func buildPod(
	namespace string,
	name string,
	session labsession.Session,
	lab lab.Lab,
) (*k8scorev1.Pod, error) {
	cpuLimit, err := k8sresource.ParseQuantity(lab.CPULimit)
	if err != nil {
		return nil, fmt.Errorf(
			"parse lab CPU limit: %w",
			err,
		)
	}

	memoryLimit, err := k8sresource.ParseQuantity(lab.MemoryLimit)
	if err != nil {
		return nil, fmt.Errorf(
			"parse lab memory limit: %w",
			err,
		)
	}

	return &k8scorev1.Pod{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    makePodLabels(session),
		},
		Spec: k8scorev1.PodSpec{
			RestartPolicy: k8scorev1.RestartPolicyNever,
			Containers: []k8scorev1.Container{
				{
					Name:            "lab",
					Image:           lab.Image,
					ImagePullPolicy: k8scorev1.PullIfNotPresent,
					Command: []string{
						"/bin/sh",
						"-c",
						"sleep infinity",
					},
					Resources: k8scorev1.ResourceRequirements{
						Limits: k8scorev1.ResourceList{
							k8scorev1.ResourceCPU:    cpuLimit,
							k8scorev1.ResourceMemory: memoryLimit,
						},
					},
				},
			},
		},
	}, nil
}

func isPodReady(
	pod *k8scorev1.Pod,
) (bool, error) {
	if pod.Status.Phase == k8scorev1.PodFailed ||
		pod.Status.Phase == k8scorev1.PodSucceeded {
		return false, fmt.Errorf(
			"pod reached terminal phase before ready: phase=%s reason=%s message=%s",
			pod.Status.Phase,
			pod.Status.Reason,
			pod.Status.Message,
		)
	}

	for _, condition := range pod.Status.Conditions {
		if condition.Type == k8scorev1.PodReady &&
			condition.Status == k8scorev1.ConditionTrue {
			return true, nil
		}
	}

	return false, nil
}
