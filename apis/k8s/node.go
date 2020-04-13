package k8s

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kubetabs/models"
	"kubetabs/models/k8s"
	"kubetabs/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

func GetNodeList(c *gin.Context) {
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}
	hostNameForm := strings.TrimSpace(c.Request.FormValue("hostname"))
	addressForm := strings.TrimSpace(c.Request.FormValue("address"))
	statusForm := strings.TrimSpace(c.Request.FormValue("status"))

	clientset := pkg.KubeConfig()
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	count := len(nodes.Items)
	var data k8s.Nodes
	var result []k8s.Nodes = make([]k8s.Nodes, count, count)
	for key, node := range nodes.Items {

		data.Id = key + 1

		// node主机名称
		data.HostName = node.ObjectMeta.Name

		// node主机IP
		for _, address := range node.Status.Addresses {
			if address.Type == "InternalIP" {
				data.Address = address.Address
			}
		}

		// node主机状态
		for _, status := range node.Status.Conditions {
			if status.Type == "Ready" {
				data.Status = string(status.Status)
			}
		}

		// node主机角色
		for key, _ := range node.ObjectMeta.Labels {
			role := strings.Split(key, "/")
			if role[0] == "node-role.kubernetes.io" {
				data.Role = role[1]
			}
		}

		// node节点信息
		data.OsImage = node.Status.NodeInfo.OSImage
		data.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion

		// 返回当前node的metrics
		metricsUsage := GetNodeMetricsUsgae(node.ObjectMeta.Name)

		// 计算CPU（core）CPU的利用率
		cpu := node.Status.Allocatable.Cpu().Value()
		cpuUsed := metricsUsage.Cpu().Value()
		data.CpuUsage = fmt.Sprintf("%.1f%s", float64(cpuUsed) / float64(cpu), "%")
		data.CpuInfo = fmt.Sprintf("%s/%d core", metricsUsage.Cpu().String(), cpu)

		// 计算内存（GB）内存的利用率
		memory := node.Status.Allocatable.Memory().Value()
		memoryUsed := metricsUsage.Memory().Value()
		data.MemoryUsage = fmt.Sprintf("%.1f%s", float64(memoryUsed) / float64(memory) * 100, "%")
		data.MemoryInfo = fmt.Sprintf("%d/%d GB", memoryUsed / 1024 / 1024 / 1024, memory / 1024 / 1024 / 1024)

		// 计算Pod Pod的利用率
		pod := node.Status.Allocatable.Pods().Value()
		podUsed := GetPodsNoHttp(node.ObjectMeta.Name)
		data.PodUsage = fmt.Sprintf("%.1f%s", float64(podUsed) / float64(pod) * 100, "%")
		data.PodInfo = fmt.Sprintf("%d/%d", podUsed, pod)

		// 计算磁盘（GB）磁盘的利用率
		storage := node.Status.Allocatable.StorageEphemeral().Value()
		data.StorageInfo = fmt.Sprintf("%dGB",storage / 1024 / 1024 / 1024)

		// node的创建时间
		data.Age = node.ObjectMeta.CreationTimestamp
		result[key] = data
	}

	//var resultFilter []k8s.Nodes
	//for _, node := range result {
	//	if strings.Contains(node.HostName, hostnameForm) && strings.Contains(node.Address, addressForm) && strings.Contains(node.Status, statusForm) {
	//		resultFilter =  append(resultFilter, node)
	//	}
	//}

	var hostNameFilter []k8s.Nodes
	var addressFilter []k8s.Nodes
	var statusFilter []k8s.Nodes

	if hostNameForm != "" {
		for _, node := range result {
			if strings.Contains(node.HostName, hostNameForm) {
				hostNameFilter = append(hostNameFilter, node)
			}
		}
	} else {
		hostNameFilter = result
	}

	if addressForm != "" {
		for _, node := range hostNameFilter {
			if strings.Contains(node.Address, addressForm) {
				addressFilter = append(addressFilter, node)
			}
		}
	} else {
		addressFilter = hostNameFilter
	}

	if statusForm != "" {
		for _, node := range addressFilter {
			if node.Status == statusForm {
				statusFilter = append(statusFilter, node)
			}
		}
	} else {
		statusFilter = addressFilter
	}

	countFilter := len(statusFilter)
	var mp = make(map[string]interface{}, 3)
	pStart := (pageIndex - 1) * pageSize
	pEnd := pageIndex * pageSize
	if pStart >= countFilter {
		mp["list"] = ""
	} else if pEnd >= countFilter {
		mp["list"] = statusFilter[pStart:countFilter]
	} else {
		mp["list"] = statusFilter[pStart:pEnd]
	}
	mp["count"] = countFilter
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

func GetNode(c *gin.Context) {
	var err error
	hostname := c.Param("hostname")

	clientset := pkg.KubeConfig()
	node, err := clientset.CoreV1().Nodes().Get(hostname, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	var mp = make(map[string]interface{}, 3)
	mp["list"] = node

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

func GetNodeMetricsUsgae(hostname string) v1.ResourceList {
	clientset := pkg.MetricsConfig()

	nodeMetric, err := clientset.MetricsV1beta1().NodeMetricses().Get(hostname, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return nodeMetric.Usage
}

func GetNodeMetricsList(c *gin.Context) {
	clientset := pkg.MetricsConfig()

	nodeMetrics, err := clientset.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var mp = make(map[string]interface{}, 3)
	mp["count"] = len(nodeMetrics.Items)
	mp["list"] = nodeMetrics.Items

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

func GetNodeMetrics(c *gin.Context) {
	hostname := c.Param("hostname")

	clientset := pkg.MetricsConfig()

	nodeMetric, err := clientset.MetricsV1beta1().NodeMetricses().Get(hostname, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var mp = make(map[string]interface{}, 3)
	mp["list"] = nodeMetric

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}
