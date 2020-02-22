package pods

import (
	cachev1alpha1 "github.com/Sher-Chowdhury/mysql-operator/pkg/apis/cache/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPodForCR returns a busybox pod with the same name/namespace as the cr
func NewPodForCR(cr *cachev1alpha1.MySQL) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	mysqlEnvVars := cr.Spec.Environment

	mysqlPvcVolumeSource := &corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: cr.Name + "-pvc",
	}

	mysqlVolumeSource := corev1.VolumeSource{
		PersistentVolumeClaim: mysqlPvcVolumeSource,
	}

	mysqlVolumes := []corev1.Volume{
		{
			Name:         "mysql-pvc-provisioned-volume",
			VolumeSource: mysqlVolumeSource,
		},
	}

	mysqlVolumeMounts := []corev1.VolumeMount{
		{
			Name:      "mysql-pvc-provisioned-volume",
			MountPath: "/var/lib/mysql",
		},
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Volumes: mysqlVolumes,
			Containers: []corev1.Container{
				{
					Name:  "mysqldb",
					Image: "docker.io/mysql:latest",
					// Command: []string{"sleep", "3600"},
					Env: []corev1.EnvVar{
						{
							Name:  "MYSQL_ROOT_PASSWORD",
							Value: mysqlEnvVars.MysqlRootPassword,
						},
						{
							Name:  "MYSQL_DATABASE",
							Value: mysqlEnvVars.MysqlDatabase,
						},
						{
							Name:  "MYSQL_USER",
							Value: mysqlEnvVars.MysqlUser,
						},
						{
							Name:  "MYSQL_PASSWORD",
							Value: mysqlEnvVars.MysqlPassword,
						},
					},
					VolumeMounts: mysqlVolumeMounts,
				},
			},
		},
	}
}
