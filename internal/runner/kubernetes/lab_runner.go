package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
	runnercontracts "github.com/maxpetrikov/maxpetrikovcom-backend/internal/runner/contracts"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8swait "k8s.io/apimachinery/pkg/util/wait"
	k8sclient "k8s.io/client-go/kubernetes"
)

var errNotImplemented = errors.New(
	"kubernetes lab runner is not implemented yet",
)

var _ runnercontracts.LabRunner = (*LabRunner)(nil)

type LabRunner struct {
	client          k8sclient.Interface
	namespace       string
	podReadyTimeout time.Duration
	logger          *slog.Logger
}

func NewLabRunner(
	client k8sclient.Interface,
	namespace string,
	podReadyTimeout time.Duration,
	logger *slog.Logger,
) *LabRunner {
	return &LabRunner{
		client:          client,
		namespace:       namespace,
		podReadyTimeout: podReadyTimeout,
		logger:          logger,
	}
}

func (lr *LabRunner) Start(
	ctx context.Context,
	session labsession.Session,
	lab lab.Lab,
) (runnercontracts.StartResult, error) {
	podName := makePodName(session)

	lr.logger.Info(
		"starting Kubernetes lab environment",
		"lab_session_id", session.ID,
		"lab_id", lab.ID,
		"namespace", lr.namespace,
		"pod_name", podName,
		"pod_ready_timeout", lr.podReadyTimeout.String(),
	)

	if err := lr.createPod(ctx, session, lab, podName); err != nil {
		return runnercontracts.StartResult{}, err
	}

	if err := lr.waitPodReady(ctx, podName); err != nil {
		return runnercontracts.StartResult{}, err
	}

	return runnercontracts.StartResult{
		Namespace: lr.namespace,
		PodName:   podName,
	}, nil
}

func (lr *LabRunner) Stop(
	ctx context.Context,
	session labsession.Session,
) error {
	lr.logger.Info(
		"stopping Kubernetes lab environment",
		"lab_session_id", session.ID,
		"namespace", lr.namespace,
		"pod_name", makePodName(session),
	)

	return errNotImplemented
}

func (lr *LabRunner) createPod(
	ctx context.Context,
	session labsession.Session,
	lab lab.Lab,
	podName string,
) error {
	pod, err := buildPod(
		lr.namespace,
		podName,
		session,
		lab,
	)
	if err != nil {
		return err
	}

	_, err = lr.client.CoreV1().
		Pods(lr.namespace).
		Create(
			ctx,
			pod,
			k8smetav1.CreateOptions{},
		)
	if k8serrors.IsAlreadyExists(err) {
		lr.logger.Info(
			"Kubernetes lab pod already exists",
			"lab_session_id", session.ID,
			"namespace", lr.namespace,
			"pod_name", podName,
		)

		return nil
	}
	if err != nil {
		return fmt.Errorf(
			"create Kubernetes lab pod: %w",
			err,
		)
	}

	return nil
}

func (lr *LabRunner) waitPodReady(
	ctx context.Context,
	podName string,
) error {
	err := k8swait.PollUntilContextTimeout(
		ctx,
		time.Second,
		lr.podReadyTimeout,
		true,
		func(ctx context.Context) (bool, error) {
			return lr.checkPodReady(ctx, podName)
		},
	)
	if err != nil {
		return fmt.Errorf(
			"wait for Kubernetes lab pod ready: %w",
			err,
		)
	}

	return nil
}

func (lr *LabRunner) checkPodReady(
	ctx context.Context,
	podName string,
) (bool, error) {
	pod, err := lr.client.CoreV1().
		Pods(lr.namespace).
		Get(
			ctx,
			podName,
			k8smetav1.GetOptions{},
		)
	if err != nil {
		return false, fmt.Errorf(
			"get Kubernetes lab pod: %w",
			err,
		)
	}

	return isPodReady(pod)
}
