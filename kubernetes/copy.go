package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/caspr-io/mu-kit/util"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

type IllegalPath struct {
	path string
}

func (e *IllegalPath) Error() string { return fmt.Sprintf("Illegal path: %s", e.path) }

type StdErrOutput struct {
	output string
}

func (e *StdErrOutput) Error() string { return fmt.Sprintf("STDERR: %s", e.output) }

// UploadToK8s uploads a single file to Kubernetes.
func (k8s *K8s) UploadToK8s(ctx context.Context, src, dest string, reader io.Reader) error {
	logger := log.Ctx(ctx).With().Logger()

	pSplit := strings.Split(dest, "/")
	if err := validateK8sPath(pSplit); err != nil {
		return err
	}

	if len(pSplit) == 3 {
		_, fileName := filepath.Split(src)
		pSplit = append(pSplit, fileName)
	}

	namespace, podName, containerName, pathToCopy := initK8sVariables(pSplit)
	logger = logger.With().Str("namespace", namespace).Str("pod", podName).Str("container", containerName).Logger()
	ctx = logger.WithContext(ctx)

	dir, _ := filepath.Split(pathToCopy)
	command := []string{"mkdir", "-p", dir}

	log.Debug().Str("path", dir).Msg("Creating directory")

	if err := k8s.doExec(ctx, namespace, podName, containerName, command, nil, nil); err != nil {
		logger.Error().Err(err).Str("path", dir).Msg("Error creating directory")
		return err
	}

	command = []string{"touch", pathToCopy}

	log.Debug().Str("path", pathToCopy).Msg("Creating file")

	if err := k8s.doExec(ctx, namespace, podName, containerName, command, nil, nil); err != nil {
		logger.Error().Err(err).Str("path", dir).Msg("Error creating file")
		return err
	}

	command = []string{"cp", "/dev/stdin", pathToCopy}

	log.Debug().Str("path", pathToCopy).Msg("Copying file contents")

	if err := k8s.doExec(ctx, namespace, podName, containerName, command, reader, nil); err != nil {
		logger.Error().Err(err).Str("path", dir).Msg("Error copying file contents")
		return err
	}

	return nil
}

func (k8s *K8s) doExec(ctx context.Context, namespace, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer) error {
	attempt, attempts := 0, 3

	for attempt < attempts {
		attempt++

		collector := &util.ErrorCollector{}
		stderr, err := k8s.Exec(ctx, namespace, podName, containerName, command, stdin, stdout)

		collector.Collect(err)

		if len(stderr) != 0 {
			collector.Collect(&StdErrOutput{(string)(stderr)})
		}

		if collector.HasErrors() {
			if attempt == attempts {
				return collector
			}

			util.Sleep(attempt)

			continue
		}

		break
	}

	return nil
}

func (k8s *K8s) Exec(ctx context.Context, namespace, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer) ([]byte, error) {
	clientset, config := k8s.Clientset, k8s.Config

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()

	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("error adding to scheme: %w", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&corev1.PodExecOptions{
		Command:   command,
		Container: containerName,
		Stdin:     stdin != nil,
		Stdout:    stdout != nil,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("error while creating Executor: %w", err)
	}

	var stderr bytes.Buffer

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		return nil, fmt.Errorf("error in Stream: %w", err)
	}

	return stderr.Bytes(), nil
}

func validateK8sPath(pathSplit []string) error {
	if len(pathSplit) >= 3 {
		return nil
	}

	return &IllegalPath{filepath.Join(pathSplit...)}
}

func initK8sVariables(split []string) (string, string, string, string) {
	namespace := split[0]
	pod := split[1]
	container := split[2]
	path := getAbsPath(split[3:]...)

	return namespace, pod, container, path
}

func getAbsPath(path ...string) string {
	return filepath.Join("/", filepath.Join(path...))
}
