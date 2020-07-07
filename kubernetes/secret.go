package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s *K8s) CreateSecret(ctx context.Context, namespace, name string, labels map[string]string, key string, payload string) (*v1.Secret, error) {
	return k8s.CoreV1().Secrets(namespace).Create(ctx, &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		StringData: map[string]string{
			key: payload,
		},
	}, metav1.CreateOptions{})
}

func (k8s *K8s) DeleteSecret(ctx context.Context, namespace string, name string) error {
	err := k8s.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	return err
}
