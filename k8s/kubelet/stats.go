package kubelet

import (
	"code.cloudfoundry.org/lager"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DiskMetricsClient struct {
	nodeClient    NodeAPI
	kubeletClient API
	namespace     string
	logger        lager.Logger
}

func NewDiskMetricsClient(nodeClient NodeAPI, kubeletClient API, namespace string, logger lager.Logger) DiskMetricsClient {
	return DiskMetricsClient{
		nodeClient:    nodeClient,
		kubeletClient: kubeletClient,
		namespace:     namespace,
		logger:        logger,
	}
}

func (d DiskMetricsClient) GetPodMetrics() (map[string]float64, error) {
	metrics := map[string]float64{}
	pods := []PodStats{}
	nodes, err := d.nodeClient.List(metav1.ListOptions{})
	if err != nil {
		return metrics, errors.Wrap(err, "failed to list nodes")
	}
	for _, n := range nodes.Items {
		statsSummary, err := d.kubeletClient.StatsSummary(n.Name)
		if err != nil {
			d.logger.Error("failed-to-get-stats-summary", err, lager.Data{"node-name": n.Name})
		}
		pods = append(pods, statsSummary.Pods...)
	}
	for _, p := range pods {
		if p.PodRef.Namespace == d.namespace && len(p.Containers) != 0 {
			logsBytes := getUsedBytes(p.Containers[0].Logs)
			rootfsBytes := getUsedBytes(p.Containers[0].Rootfs)
			metrics[p.PodRef.Name] = logsBytes + rootfsBytes
		}
	}
	return metrics, nil
}

func getUsedBytes(stats *FsStats) float64 {
	if stats == nil || stats.UsedBytes == nil {
		return 0
	}
	return float64(*stats.UsedBytes)
}
