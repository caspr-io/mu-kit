package kubernetes

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sJob struct {
	Name string
	job  *batchv1.Job
}

func NewJobSpec(name string, labels map[string]string, containerName string, containerImage string) *K8sJob {
	backoffLimit := int32(0)
	j := &batchv1.Job{
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

	return &K8sJob{Name: name, job: j}
}

func (j *K8sJob) AddAndMountVolume(volume v1.Volume, containerName string, mountPath string) {
	if j.job.Spec.Template.Spec.Volumes == nil {
		j.job.Spec.Template.Spec.Volumes = []v1.Volume{}
	}

	j.job.Spec.Template.Spec.Volumes = append(j.job.Spec.Template.Spec.Volumes, volume)

	vm := v1.VolumeMount{
		Name:      volume.Name,
		MountPath: mountPath,
	}

	for i, c := range j.job.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			if c.VolumeMounts == nil {
				c.VolumeMounts = []v1.VolumeMount{vm}
			} else {
				c.VolumeMounts = append(c.VolumeMounts, vm)
			}

			j.job.Spec.Template.Spec.Containers[i] = c
		}
	}
}

func (j *K8sJob) AddConfigMapVolume(cm *v1.ConfigMap, containerName string, mountPath string) {
	v := v1.Volume{
		Name: strings.ReplaceAll(mountPath, "/", "-"),
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{Name: cm.Name},
			},
		},
	}

	j.AddAndMountVolume(v, containerName, mountPath)
}

func (j *K8sJob) AddSecretVolume(sec *v1.Secret, containerName string, mountPath string) {
	v := v1.Volume{
		Name: sec.Name,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: sec.Name,
			},
		},
	}

	j.AddAndMountVolume(v, containerName, mountPath)
}

func (j *K8sJob) AddEnvironment(containerName string, env map[string]string) {
	envVars := []v1.EnvVar{}
	for k, v := range env {
		envVars = append(envVars, v1.EnvVar{Name: k, Value: v})
	}

	for i, c := range j.job.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			c.Env = envVars
			j.job.Spec.Template.Spec.Containers[i] = c
		}
	}
}

func (j *K8sJob) AddPullSecret(secretName string) {
	j.job.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{{Name: secretName}}
}

func (k8s *K8s) CreateJob(namespace string, jobSpec *K8sJob) error {
	_, err := k8s.BatchV1().Jobs(namespace).Create(jobSpec.job)
	return err
}
