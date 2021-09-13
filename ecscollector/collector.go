// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License").
// You may not use this file except in compliance with the License.
// A copy of the License is located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package ecscollector implements a Prometheus collector for Amazon ECS
// metrics available at the ECS metadata server.
package ecscollector

import (
	"context"
	"log"

	"github.com/prometheus-community/ecs_exporter/ecsmetadata"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	numProcsDesc = prometheus.NewDesc(
		"ecs_num_procs",
		"Number of processes.",
		labelKeys, nil)

	cpuTotalDesc = prometheus.NewDesc(
		"ecs_cpu_total",
		"Total CPU usage.",
		labelKeys, nil)

	cpuUserDesc = prometheus.NewDesc(
		"ecs_cpu_user",
		"Total CPU usage by user space.",
		labelKeys, nil)

	cpuKernelDesc = prometheus.NewDesc(
		"ecs_cpu_kernel",
		"Total CPU usage by kernel space.",
		labelKeys, nil)

	cpuSystemDesc = prometheus.NewDesc(
		"ecs_cpu_system",
		"Total system CPU usage.",
		labelKeys, nil)

	cpuNumOnlineDesc = prometheus.NewDesc(
		"ecs_cpu_num_online",
		"Number of online CPUs.",
		labelKeys, nil)

	memUsageDesc = prometheus.NewDesc(
		"ecs_mem_total",
		"Total memory usage.",
		labelKeys, nil)

	memMaxUsageDesc = prometheus.NewDesc(
		"ecs_mem_usage",
		"Maximum memory usage.",
		labelKeys, nil)

	memLimitDesc = prometheus.NewDesc(
		"ecs_mem_limit",
		"Memory limit.",
		labelKeys, nil)

	networkRxBytesDesc = prometheus.NewDesc(
		"ecs_network_rx_bytes",
		"Network recieved in bytes.",
		labelKeysWithNetwork, nil)

	networkRxPacketsDesc = prometheus.NewDesc(
		"ecs_network_rx_packets",
		"Network packets recieved.",
		labelKeysWithNetwork, nil)

	networkRxDroppedDesc = prometheus.NewDesc(
		"ecs_network_rx_dropped",
		"Network packets dropped in recieving.",
		labelKeysWithNetwork, nil)

	networkRxErrorsDesc = prometheus.NewDesc(
		"ecs_network_rx_errors",
		"Network errors in recieving.",
		labelKeysWithNetwork, nil)

	networkTxBytesDesc = prometheus.NewDesc(
		"ecs_network_tx_bytes",
		"Network transmitted in bytes.",
		labelKeysWithNetwork, nil)

	networkTxPacketsDesc = prometheus.NewDesc(
		"ecs_network_tx_packets",
		"Network packets transmitted.",
		labelKeysWithNetwork, nil)

	networkTxDroppedDesc = prometheus.NewDesc(
		"ecs_network_tx_dropped",
		"Network packets dropped in transmit.",
		labelKeysWithNetwork, nil)

	networkTxErrorsDesc = prometheus.NewDesc(
		"ecs_network_tx_errors",
		"Network errors in transmit.",
		labelKeysWithNetwork, nil)

	networkRxRateDesc = prometheus.NewDesc(
		"ecs_network_rx_rate",
		"Network recieved rate per second.",
		labelKeys, nil)

	networkTxRateDesc = prometheus.NewDesc(
		"ecs_network_tx_rate",
		"Network transmitted rate per second.",
		labelKeys, nil)
)

var labelKeys = []string{
	"cluster",
	"task_az",
	"task_family",
	"task_family_revision",
	// TODO(jbd): Add service name as a label.
	"container",
}

var labelKeysWithNetwork = append(labelKeys,
	"network_interface",
)

// NewCollector returns a new Collector that queries ECS metadata server
// for ECS task and container metrics.
func NewCollector(client *ecsmetadata.Client) prometheus.Collector {
	return &collector{client: client}
}

type collector struct {
	client *ecsmetadata.Client
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- numProcsDesc
	ch <- cpuTotalDesc
	ch <- cpuUserDesc
	ch <- cpuKernelDesc
	ch <- cpuSystemDesc
	ch <- cpuNumOnlineDesc
	ch <- memUsageDesc
	ch <- memMaxUsageDesc
	ch <- memLimitDesc
	ch <- networkRxBytesDesc
	ch <- networkRxPacketsDesc
	ch <- networkRxDroppedDesc
	ch <- networkRxErrorsDesc
	ch <- networkTxBytesDesc
	ch <- networkTxPacketsDesc
	ch <- networkTxDroppedDesc
	ch <- networkTxErrorsDesc
	ch <- networkRxRateDesc
	ch <- networkTxRateDesc
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	metadata, err := c.client.RetrieveTaskMetadata(ctx)
	if err != nil {
		log.Printf("Failed to retrieve metadata: %v", err)
		return
	}
	stats, err := c.client.RetrieveTaskStats(ctx)
	if err != nil {
		log.Printf("Failed to retrieve container stats: %v", err)
		return
	}
	for _, container := range metadata.Containers {
		s := stats[container.DockerID]
		if s == nil {
			log.Printf("Couldn't find container with ID %q in stats", container.DockerID)
			continue
		}

		labelVals := []string{
			metadata.Cluster,
			metadata.AvailabilityZone,
			metadata.Family,
			metadata.Revision,
			container.Name,
		}
		metrics := map[*prometheus.Desc]float64{
			numProcsDesc:      s.NumProcs,
			cpuTotalDesc:      s.CPUStats.CPUUsage.TotalUsage,
			cpuKernelDesc:     s.CPUStats.CPUUsage.UsageInKernelmode,
			cpuUserDesc:       s.CPUStats.CPUUsage.UsageInUsermode,
			cpuSystemDesc:     s.CPUStats.SystemCPUUsage,
			cpuNumOnlineDesc:  s.CPUStats.OnlineCPUs,
			memUsageDesc:      s.MemoryStats.Usage,
			memMaxUsageDesc:   s.MemoryStats.MaxUsage,
			memLimitDesc:      s.MemoryStats.Limit,
			networkRxRateDesc: s.NetworkRateStats.RxBytesPerSec,
			networkTxRateDesc: s.NetworkRateStats.TxBytesPerSec,
		}
		for desc, value := range metrics {
			ch <- prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				value,
				labelVals...,
			)
		}

		// Network metrics per inteface.
		for iface, netStats := range s.Networks {
			networkLabelVals := append(labelVals, iface)

			ch <- prometheus.MustNewConstMetric(
				networkRxBytesDesc,
				prometheus.GaugeValue,
				netStats.RxBytes,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkRxPacketsDesc,
				prometheus.GaugeValue,
				netStats.RxPackets,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkRxDroppedDesc,
				prometheus.GaugeValue,
				netStats.RxDropped,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkRxErrorsDesc,
				prometheus.GaugeValue,
				netStats.RxErrors,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkTxBytesDesc,
				prometheus.GaugeValue,
				netStats.TxBytes,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkTxPacketsDesc,
				prometheus.GaugeValue,
				netStats.TxPackets,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkTxDroppedDesc,
				prometheus.GaugeValue,
				netStats.TxDropped,
				networkLabelVals...,
			)
			ch <- prometheus.MustNewConstMetric(
				networkTxErrorsDesc,
				prometheus.GaugeValue,
				netStats.TxErrors,
				networkLabelVals...,
			)
		}
	}
}
