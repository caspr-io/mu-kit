package kubernetes

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewJobSpec(name string, labels map[string]string, containerName string, containerImage string) *batchv1.Job {
	backoffLimit := int32(0)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  containerName,
							Image: containerImage,
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func (k8s *K8s) CreateJob(namespace string, jobSpec *batchv1.JobSpec) {}
