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

// Package ecsmetadata queries ECS Metadata Server for ECS task metrics.
// This package is currently experimental and is subject to change.
package ecsmetadata

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var endpoint string

func init() {
	const endpointEnv = "ECS_CONTAINER_METADATA_URI_V4"
	endpoint = os.Getenv(endpointEnv)
	if endpoint == "" {
		log.Fatalf("%q environmental variable is not set, are you running this on ECS?", endpointEnv)
	}
}

type Client struct {
	// httpClient is the client to use when making HTTP requests
	// when set. Otherwise, an internal default client will be used.
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{httpClient: httpClient}, nil
}

func (c *Client) RetrieveTaskStats(ctx context.Context) (map[string]*ContainerStats, error) {
	out := make(map[string]*ContainerStats)
	err := c.request(ctx, endpoint+"/task/stats", &out)
	return out, err
}

func (c *Client) RetrieveTaskMetadata(ctx context.Context) (*TaskMetadata, error) {
	var out TaskMetadata
	err := c.request(ctx, endpoint+"/task", &out)
	// Replace cluster ARN with cluster name.
	if index := strings.LastIndex(out.Cluster, "/"); index > 0 {
		out.Cluster = out.Cluster[index+1:]
	}
	// TODO(jbd): Consider using taskWithTags for task tags.
	return &out, err
}

func (c *Client) request(ctx context.Context, uri string, out interface{}) error {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("%s", body)
	return json.Unmarshal(body, out)
}

type ContainerStats struct {
	Name     string  `json:"name"`
	ID       string  `json:"id"`
	NumProcs float64 `json:"num_procs"`

	CPUStats struct {
		CPUUsage struct {
			TotalUsage        float64   `json:"total_usage"`
			PercpuUsage       []float64 `json:"percpu_usage"`
			UsageInKernelmode float64   `json:"usage_in_kernelmode"`
			UsageInUsermode   float64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage float64 `json:"system_cpu_usage"`
		OnlineCPUs     float64 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          float64 `json:"periods"`
			ThrottledPeriods float64 `json:"throttled_periods"`
			ThrottledTime    float64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`

	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        float64   `json:"total_usage"`
			PercpuUsage       []float64 `json:"percpu_usage"`
			UsageInKernelmode float64   `json:"usage_in_kernelmode"`
			UsageInUsermode   float64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage float64 `json:"system_cpu_usage"`
		OnlineCPUs     float64 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          float64 `json:"periods"`
			ThrottledPeriods float64 `json:"throttled_periods"`
			ThrottledTime    float64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`

	MemoryStats struct {
		Usage    float64 `json:"usage"`
		MaxUsage float64 `json:"max_usage"`
		Stats    struct {
			ActiveAnon              float64 `json:"active_anon"`
			ActiveFile              float64 `json:"active_file"`
			Cache                   float64 `json:"cache"`
			Dirty                   float64 `json:"dirty"`
			HierarchicalMemoryLimit float64 `json:"hierarchical_memory_limit"`
			HierarchicalMemswLimit  float64 `json:"hierarchical_memsw_limit"`
			InactiveAnon            float64 `json:"inactive_anon"`
			InactiveFile            float64 `json:"inactive_file"`
			MappedFile              float64 `json:"mapped_file"`
			Pgfault                 float64 `json:"pgfault"`
			Pgmajfault              float64 `json:"pgmajfault"`
			Pgpgin                  float64 `json:"pgpgin"`
			Pgpgout                 float64 `json:"pgpgout"`
			RSS                     float64 `json:"rss"`
			RSSHuge                 float64 `json:"rss_huge"`
			TotalActiveAnon         float64 `json:"total_active_anon"`
			TotalActiveFile         float64 `json:"total_active_file"`
			TotalCache              float64 `json:"total_cache"`
			TotalDirty              float64 `json:"total_dirty"`
			TotalInactiveAnon       float64 `json:"total_inactive_anon"`
			TotalInactiveFile       float64 `json:"total_inactive_file"`
			TotalMappedFile         float64 `json:"total_mapped_file"`
			TotalPgfault            float64 `json:"total_pgfault"`
			TotalPgmajfault         float64 `json:"total_pgmajfault"`
			TotalPgpgin             float64 `json:"total_pgpgin"`
			TotalPgpgout            float64 `json:"total_pgpgout"`
			TotalRSS                float64 `json:"total_rss"`
			TotalRSSHuge            float64 `json:"total_rss_huge"`
			TotalUnevictable        float64 `json:"total_unevictable"`
			TotalWriteback          float64 `json:"total_writeback"`
			Unevictable             float64 `json:"unevictable"`
			Writeback               float64 `json:"writeback"`
		} `json:"stats"`
		Limit float64 `json:"limit"`
	} `json:"memory_stats"`

	Networks map[string]struct {
		RxBytes   float64 `json:"rx_bytes"`
		RxPackets float64 `json:"rx_packets"`
		RxErrors  float64 `json:"rx_errors"`
		RxDropped float64 `json:"rx_dropped"`
		TxBytes   float64 `json:"tx_bytes"`
		TxPackets float64 `json:"tx_packets"`
		TxErrors  float64 `json:"tx_errors"`
		TxDropped float64 `json:"tx_dropped"`
	} `json:"networks"`

	NetworkRateStats struct {
		RxBytesPerSec float64 `json:"rx_bytes_per_sec"`
		TxBytesPerSec float64 `json:"tx_bytes_per_sec"`
	} `json:"network_rate_stats"`
}

// TODO(jbd): Add storage stats.

type TaskMetadata struct {
	Cluster          string `json:"Cluster"`
	TaskARN          string `json:"TaskARN"`
	Family           string `json:"Family"`
	Revision         string `json:"Revision"`
	DesiredStatus    string `json:"DesiredStatus"`
	KnownStatus      string `json:"KnownStatus"`
	AvailabilityZone string `json:"AvailabilityZone"`
	LaunchType       string `json:"LaunchType"`
	Containers       []struct {
		DockerID      string            `json:"DockerId"`
		Name          string            `json:"Name"`
		DockerName    string            `json:"DockerName"`
		Image         string            `json:"Image"`
		ImageID       string            `json:"ImageID"`
		Labels        map[string]string `json:"Labels"`
		DesiredStatus string            `json:"DesiredStatus"`
		KnownStatus   string            `json:"KnownStatus"`
		Type          string            `json:"Type"`
		ContainerARN  string            `json:"ContainerARN"`
	} `json:"Containers"`
}
