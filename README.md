# ecs_prometheus_exporter

🚧 🚧 🚧 This repo is still work in progress and is subject to change.

This repo contains a Prometheus exporter for
Amazon Elastic Container Service (ECS) that publishes
ECS task infra metrics in Prometheus format.

You need to run it as a sidecar on your ECS tasks.

By default, it publishes Prometheus metrics on ":9779/metrics". The exporter in this repo can be a useful complementary sidecar for the scenario described in [this blog post](https://aws.amazon.com/blogs/opensource/metrics-collection-from-amazon-ecs-using-amazon-managed-service-for-prometheus/). Adding this sidecar to the ECS task definition would export task-level metrics in addition to the custom metrics described in the blog.  

The sidecar process is also supported on [AWS App Runner](https://aws.amazon.com/apprunner/)
and can be used to publish infra metrics in Prometheus format
from App Runner services.

## Labels

* **cluster**: The ECS cluster the task is running.
* **task_family**: The ECS task family name.
* **task_family_revision**: The ECS task family revision.
* **task_az**: The availability zone the ECS task is running. Example: us-east-1a.
* **container**: Container associated with a metric.
* **network_interface**: Network interface associated with the metric. Only
  available for several network metrics.

## Example output

```
# HELP ecs_cpu_kernel Total CPU usage by kernel space.
# TYPE ecs_cpu_kernel gauge
ecs_cpu_kernel{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 7.5e+08
# HELP ecs_cpu_num_online Number of online CPUs.
# TYPE ecs_cpu_num_online gauge
ecs_cpu_num_online{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 2
# HELP ecs_cpu_system Total system CPU usage.
# TYPE ecs_cpu_system gauge
ecs_cpu_system{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 8.638668e+14
# HELP ecs_cpu_total Total CPU usage.
# TYPE ecs_cpu_total gauge
ecs_cpu_total{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 2.761418681e+09
# HELP ecs_cpu_user Total CPU usage by user space.
# TYPE ecs_cpu_user gauge
ecs_cpu_user{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 1.46e+09
# HELP ecs_mem_limit Memory limit.
# TYPE ecs_mem_limit gauge
ecs_mem_limit{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 9.223372036854772e+18
# HELP ecs_mem_total Total memory usage.
# TYPE ecs_mem_total gauge
ecs_mem_total{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 3.7888e+06
# HELP ecs_mem_usage Maximum memory usage.
# TYPE ecs_mem_usage gauge
ecs_mem_usage{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 8.404992e+06
# HELP ecs_network_rx_bytes Network recieved in bytes.
# TYPE ecs_network_rx_bytes gauge
ecs_network_rx_bytes{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 9.217898e+06
# HELP ecs_network_rx_dropped Network packets dropped in recieving.
# TYPE ecs_network_rx_dropped gauge
ecs_network_rx_dropped{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 0
# HELP ecs_network_rx_errors Network errors in recieving.
# TYPE ecs_network_rx_errors gauge
ecs_network_rx_errors{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 0
# HELP ecs_network_rx_packets Network packets recieved.
# TYPE ecs_network_rx_packets gauge
ecs_network_rx_packets{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 80907
# HELP ecs_network_rx_rate Network recieved rate per second.
# TYPE ecs_network_rx_rate gauge
ecs_network_rx_rate{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 896.9922746540345
# HELP ecs_network_tx_bytes Network transmitted in bytes.
# TYPE ecs_network_tx_bytes gauge
ecs_network_tx_bytes{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 7.630328e+06
# HELP ecs_network_tx_dropped Network packets dropped in transmit.
# TYPE ecs_network_tx_dropped gauge
ecs_network_tx_dropped{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 0
# HELP ecs_network_tx_errors Network errors in transmit.
# TYPE ecs_network_tx_errors gauge
ecs_network_tx_errors{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 0
# HELP ecs_network_tx_packets Network packets transmitted.
# TYPE ecs_network_tx_packets gauge
ecs_network_tx_packets{cluster="prometheus-fargate",container="ecs-metadata-proxy",network_interface="eth1",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 66668
# HELP ecs_network_tx_rate Network transmitted rate per second.
# TYPE ecs_network_tx_rate gauge
ecs_network_tx_rate{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 1551.1866404050595
# HELP ecs_num_procs Number of processes.
# TYPE ecs_num_procs gauge
ecs_num_procs{cluster="prometheus-fargate",container="ecs-metadata-proxy",task_az="us-west-2c",task_family="ecs-metadata-proxy",task_family_revision="1"} 0
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 10
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.16.3"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 3.781384e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 3.781384e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.444934e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 629
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 4.20848e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 3.781384e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.1693952e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.825088e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 8693
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.1693952e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.651904e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 9322
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 9600
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 54808
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 65536
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 5.257248e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.392994e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 557056
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 557056
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.4204424e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 9
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 4
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
