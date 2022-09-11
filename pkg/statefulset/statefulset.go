package statefulset

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateStatefulsetObject(nr int, imageName string, tag string) (*appsv1.StatefulSet, error) {
	var count int32 = 1
	var name string
	name = fmt.Sprintf("test-csi-%d", nr)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &count,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test-csi", "name": name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "test-csi", "name": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: imageName + ":" + tag,
						Name:  "container",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http-80",
							Protocol:      "TCP",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      name,
							MountPath: "/etc/test",
						}},
					}},
					RestartPolicy:      "Always",
					DNSPolicy:          "ClusterFirst",
					ServiceAccountName: "default",
					Volumes: []corev1.Volume{{
						Name: name,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: name,
								ReadOnly:  false,
							},
						},
					}},
				},
			},
		},
	}
	return sts, nil
}
