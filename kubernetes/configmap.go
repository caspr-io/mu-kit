package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s *K8s) CreateConfigMap(namespace string, name string, labels map[string]string, key string, payload string) error {
	_, err := k8s.CoreV1().ConfigMaps(namespace).Create(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Data: map[string]string{
			key: payload,
		},
	})

	return err
}

func (k8s *K8s) DeleteConfigMap(namespace string, name string) error {
	err := k8s.CoreV1().ConfigMaps(namespace).Delete(name, nil)
	return err
}
