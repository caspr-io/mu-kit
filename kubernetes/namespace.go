package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s *K8s) CreateNamespace(namespace string, labels map[string]string) error {
	ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:   namespace,
		Labels: labels,
	},
	}

	_, err := k8s.CoreV1().Namespaces().Create(ns)

	return err
}

func (k8s *K8s) DeleteNamespace(namespace string) error {
	return k8s.CoreV1().Namespaces().Delete(namespace, nil)
}
