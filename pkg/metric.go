package metric

import (
	"fmt"

	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func int32Ptr(i int32) *int32 { return &i }

// Gen serialize k8s objects to YAML
func Gen(name, namespace, project, cluster, location string, metrics []string) {

	// Namespace
	if namespace != "default" {
		ns, err := yaml.Marshal(newNamespace(namespace))
		if err != nil {
			panic(err)
		}
		fmt.Println("---")
		fmt.Println("apiVersion: v1")
		fmt.Println("kind: Namespace")
		fmt.Println(string(ns))
	}

	// ServiceAccount
	sa, err := yaml.Marshal(newServiceAccount(name, namespace))
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	fmt.Println("apiVersion: v1")
	fmt.Println("kind: ServiceAccount")
	fmt.Println(string(sa))

	// ClusterRole
	cr, err := yaml.Marshal(newClusterRole(name))
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	fmt.Println("apiVersion: rbac.authorization.k8s.io/v1")
	fmt.Println("kind: ClusterRole")
	fmt.Println(string(cr))

	// ClusterRoleBinding
	crb, err := yaml.Marshal(newClusterRoleBinding(name, namespace))
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	fmt.Println("apiVersion: rbac.authorization.k8s.io/v1")
	fmt.Println("kind: ClusterRoleBinding")
	fmt.Println(string(crb))

	// ConfigMap
	cm, err := yaml.Marshal(newConfigMap(name, namespace))
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	fmt.Println("apiVersion: v1")
	fmt.Println("kind: ConfigMap")
	fmt.Println(string(cm))

	// Deployment
	deployment, err := yaml.Marshal(newDeployment(name, namespace, project, cluster, location, metrics))
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	fmt.Println("apiVersion: apps/v1")
	fmt.Println("kind: Deployment")
	fmt.Println(string(deployment))
}

func newNamespace(namespace string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
}

func newServiceAccount(name, namespace string) *corev1.ServiceAccount {
	resourceName := name + "-prom"

	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: namespace,
		},
	}
}

func newClusterRole(name string) *rbacv1.ClusterRole {
	resourceName := name + "-prom"

	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
					"nodes/proxy",
					"services",
					"endpoints",
					"pods",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				APIGroups: []string{
					"extentions",
				},
				Resources: []string{
					"ingresses",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
			{
				NonResourceURLs: []string{
					"/metrics",
				},
				Verbs: []string{
					"get",
				},
			},
		},
	}
}

func newClusterRoleBinding(name, namespace string) *rbacv1.ClusterRoleBinding {
	resourceName := name + "-prom"

	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     resourceName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      resourceName,
				Namespace: namespace,
			},
		},
	}
}

func newConfigMap(name, namespace string) *corev1.ConfigMap {
	resourceName := name + "-prom"

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"prometheus.yml": `
scrape_configs:
- job_name: 'kubernetes-pods'
  metrics_path: /metrics
  kubernetes_sd_configs:
    - role: pod
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_annotation_cm_example_com_scrape]
      action: keep
      regex: true
    - source_labels:
      - __meta_kubernetes_pod_annotationpresent_cm_example_com_path
      - __meta_kubernetes_pod_annotation_cm_example_com_path
      action: replace
      target_label: __metrics_path__
      regex: true;(.+)
    - source_labels: [__address__, __meta_kubernetes_pod_annotation_cm_example_com_port]
      action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      target_label: __address__
    - action: labelmap
      regex: __meta_kubernetes_pod_label_(.+)
    - source_labels: [__meta_kubernetes_namespace]
      action: replace
      target_label: kubernetes_namespace
    - source_labels: [__meta_kubernetes_pod_name]
      action: replace
      target_label: kubernetes_pod_name

- job_name: 'kubernetes-nodes'
  metrics_path: /metrics
  kubernetes_sd_configs:
    - role: node
  relabel_configs:
    - source_labels: [__meta_kubernetes_node_annotation_cm_example_com_scrape]
      action: keep
      regex: true
    - source_labels: [__address__]
      action: replace
      regex: ([^:]+)(?::\d+)?
      replacement: $1:10255
      target_label: __address__
    - action: labelmap
      regex: __meta_kubernetes_node_label_(.+)

- job_name: 'kubernetes-nodes-cadvisor'
  metrics_path: /metrics/cadvisor
  kubernetes_sd_configs:
    - role: node
  relabel_configs:
    - source_labels: [__meta_kubernetes_node_annotation_cm_example_com_scrape]
      action: keep
      regex: true
    - source_labels: [__address__]
      action: replace
      regex: ([^:]+)(?::\d+)?
      replacement: $1:10255
      target_label: __address__
    - action: labelmap
      regex: __meta_kubernetes_node_label_(.+)

- job_name: 'kubernetes-apiservers'
  kubernetes_sd_configs:
    - role: endpoints
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
  relabel_configs:
    - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
      action: keep
      regex: default;kubernetes;https
      `,
		},
	}
}

func newDeployment(name, namespace, project, cluster, location string, metrics []string) *appsv1.Deployment {
	resourceName := name + "-prom"

	labels := map[string]string{
		"app":        "prometheus-server",
		"controller": name,
	}

	sidecarArgs := []string{
		fmt.Sprintf("--stackdriver.project-id=%s", project),
		fmt.Sprintf("--stackdriver.kubernetes.cluster-name=%s", cluster),
		fmt.Sprintf("--stackdriver.kubernetes.location=%s", location),
		"--prometheus.wal-directory=/prometheus/wal",
		"--log.level=debug",
	}

	for _, m := range metrics {
		sidecarArgs = append(sidecarArgs, fmt.Sprintf("--include={__name__=~\"%s\"}", m))
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: namespace,
			// OwnerReferences: []metav1.OwnerReference{
			//   *metav1.NewControllerRef(cr, cmv1alpha1.SchemeGroupVersion.WithKind("CustomMetric")),
			// },
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: resourceName,
					Containers: []corev1.Container{
						{
							Name:  "prometheus",
							Image: "prom/prometheus:v2.6.1",
							Args: []string{
								"--config.file=/etc/prometheus/prometheus.yml",
								"--storage.tsdb.path=/prometheus/",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9090,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "prometheus-config-volume",
									MountPath: "/etc/prometheus/",
								},
								{
									Name:      "prometheus-storage-volume",
									MountPath: "/prometheus/",
								},
							},
						},
						{
							Name:  "sidecar",
							Image: "gcr.io/stackdriver-prometheus/stackdriver-prometheus-sidecar:0.8.0",
							Args:  sidecarArgs,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9091,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "prometheus-storage-volume",
									MountPath: "/prometheus/",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "prometheus-config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: resourceName,
									},
									DefaultMode: int32Ptr(420),
								},
							},
						},
						{
							Name: "prometheus-storage-volume", // default to emptyDir
						},
					},
				},
			},
		},
	}
}
