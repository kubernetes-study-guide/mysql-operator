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

	var pvcSpec corev1.PersistentVolumeClaimSpec
	pvcSpec.AccessModes = []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteOnce,
	}
	pvcSpec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceStorage: resource.MustParse(volumesize),
		},
	}

	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx Accessmodes xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	println(pvcSpec.AccessModes)
	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxx  StorageClassName xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	println(pvcSpec.StorageClassName)
	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	if cr.Spec.Volume.StorageClass != "" {
		pvcSpec.StorageClassName = &cr.Spec.Volume.StorageClass
	}

	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx Accessmodes - after xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	println(pvcSpec.AccessModes)
	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxx  StorageClassName - after xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	println(pvcSpec.StorageClassName)
	println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pvc",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: pvcSpec,
	}
}
