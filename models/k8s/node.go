package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Nodes struct {
	Id int `json:"id"`
	HostName string `json:"hostname"`
	Address string `json:"address"`
	Status string `json:"status"`
	Role string `json:"role"`
	OsImage string `json:"osImage"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	CpuUsage string `json:"cpuUsage"`
	CpuInfo string `json:"cpuInfo"`
	MemoryUsage string `json:"memoryUsage"`
	MemoryInfo string `json:"memoryInfo"`
	PodUsage string `json:"podUsage"`
	PodInfo string `json:"podInfo"`
	//StorageUsage string `json:"storageUsage"`
	StorageInfo string `json:"storageInfo"`
	Age metav1.Time `json:"age"`
}
