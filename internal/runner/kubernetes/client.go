package kubernetes

import (
	"fmt"

	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(
	kubeconfigPath string,
) (k8sclient.Interface, error) {
	var config *rest.Config
	var err error

	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags(
			"",
			kubeconfigPath,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"build Kubernetes config from kubeconfig: %w",
				err,
			)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf(
				"build in-cluster Kubernetes config: %w",
				err,
			)
		}
	}

	client, err := k8sclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf(
			"create Kubernetes client: %w",
			err,
		)
	}

	return client, nil
}
