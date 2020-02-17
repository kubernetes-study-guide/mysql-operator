package pvcs

import (
	cachev1alpha1 "github.com/Sher-Chowdhury/mysql-operator/pkg/apis/cache/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPvcForCR returns a busybox pod with the same name/namespace as the cr
func NewPvcForCR(cr *cachev1alpha1.MySQL) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app": cr.Name,
	}

	volumesize := cr.Spec.Volume.VolumeSize
	storageclassname := cr.Spec.Volume.StorageClass

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pvc",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(volumesize),
				},
			},
			StorageClassName: &storageclassname,
		},
	}
}
