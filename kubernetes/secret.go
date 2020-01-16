package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s *K8s) CreateSecret(namespace string, name string, labels map[string]string, key string, payload string) error {
	_, err := k8s.CoreV1().Secrets(namespace).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		StringData: map[string]string{
			key: payload,
		},
	})

	return err
}

func (k8s *K8s) DeleteSecret(namespace string, name string) error {
	err := k8s.CoreV1().Secrets(namespace).Delete(name, nil)
	return err
}
