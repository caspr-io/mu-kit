package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s *K8s) CreateConfigMap(ctx context.Context, namespace, name string, labels map[string]string, key string, payload string) (*v1.ConfigMap, error) {
	cm, err := k8s.CoreV1().ConfigMaps(namespace).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Data: map[string]string{
			key: payload,
		},
	}, metav1.CreateOptions{})

	return cm, err
}

func (k8s *K8s) DeleteConfigMap(ctx context.Context, namespace, name string) error {
	err := k8s.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	return err
}
