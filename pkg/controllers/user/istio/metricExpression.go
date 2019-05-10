package istio

import (
	"encoding/json"
	"strings"

	"github.com/ghodss/yaml"
	managementv3 "github.com/rancher/types/apis/management.cattle.io/v3"
)

func getPredefinedIstioMetrics() []*managementv3.MonitorMetric {
	yamls := strings.Split(IstioMetricsTemplate, "\n---\n")
	var rtn []*managementv3.MonitorMetric
	for _, yml := range yamls {
		var tmp managementv3.MonitorMetric
		if err := yamlToObject(yml, &tmp); err != nil {
			panic(err)
		}
		if tmp.Name == "" {
			continue
		}
		rtn = append(rtn, &tmp)
	}

	return rtn
}

func getPredefinedIstioClusterGraph() []*managementv3.ClusterMonitorGraph {
	yamls := strings.Split(IstioClusterGraphTemplate, "\n---\n")
	var rtn []*managementv3.ClusterMonitorGraph
	for _, yml := range yamls {
		var tmp managementv3.ClusterMonitorGraph
		if err := yamlToObject(yml, &tmp); err != nil {
			panic(err)
		}
		if tmp.Name == "" {
			continue
		}
		rtn = append(rtn, &tmp)
	}

	return rtn
}

func getPredefinedIstioProjectGraph() []*managementv3.ClusterMonitorGraph {
	yamls := strings.Split(IstioProjectGraphTemplate, "\n---\n")
	var rtn []*managementv3.ClusterMonitorGraph
	for _, yml := range yamls {
		var tmp managementv3.ClusterMonitorGraph
		if err := yamlToObject(yml, &tmp); err != nil {
			panic(err)
		}
		if tmp.Name == "" {
			continue
		}
		rtn = append(rtn, &tmp)
	}

	return rtn
}

func yamlToObject(yml string, obj interface{}) error {
	jsondata, err := yaml.YAMLToJSON([]byte(yml))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsondata, obj); err != nil {
		return err
	}
	return nil
}

var (
	IstioMetricsTemplate = `
# Source: metric-expression-cluster/templates/expressionmesh.yaml
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-4xxs
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: cluster
    graph: 4xxs
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", response_code=~"4.*"}[1m])) 
  legendFormat: 4xx request count 
  description: the count of requests that response code is 4xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-5xxs
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: cluster
    graph: 5xxs
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", response_code=~"5.*"}[1m])) 
  legendFormat: 5xx request count
  description: the count of requests that response code is 5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-success
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: cluster
    graph: global-success-rate
    source: rancher-istio
spec:
  expression: sum(rate(istio_requests_total{reporter="destination", response_code!~"5.*"}[1m])) / 
    sum(rate(istio_requests_total{reporter="destination"}[1m])) 
  legendFormat: Success rate
  description: the count of requests that response code is non-5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-global-requests-volume
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: cluster
    graph: global-request-volume
    source: rancher-istio
spec:
  expression: round(sum(irate(istio_requests_total{reporter="destination"}[1m])), 0.001) 
  legendFormat: Request volume
  description: the global request of volume 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-4xxs-client
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 4xxs-client
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="source", destination_service_name=~"$service",response_code=~"4.*", destination_service_namespace=~"$namespace"}[5m])) by (source_workload, source_workload_namespace)
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]'
  description: the count of client requests that response code is 4xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-5xxs-client
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 5xxs-client
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="source", destination_service_name=~"$service",response_code=~"5.*", destination_service_namespace=~"$namespace"}[5m])) by (source_workload, source_workload_namespace)
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]'
  description: the count of client requests that response code is 5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-success-client
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-success-rate-client
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="source", destination_service_name=~"$service",response_code!~"5.*", destination_service_namespace=~"$namespace"}[5m])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter="source", destination_service_name=~"$service", destination_service_namespace=~"$namespace"}[5m])) by (source_workload, source_workload_namespace) 
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]'
  description: the count of requests that response code is non-5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-volume-client
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-request-volume-client
    source: rancher-istio
spec:
  expression: round(sum(irate(istio_requests_total{reporter="source",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (source_workload, source_workload_namespace), 0.001) 
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]' 
  description: the global client request volume 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-client-request-duration-p50
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: client-request-duration 
    source: rancher-istio
spec:
  expression: histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter="source",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (le, source_workload, source_workload_namespace)) 
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]' 
  description: the client request duration that quantile is 0.5 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-client-request-duration-p90
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: client-request-duration
    source: rancher-istio
spec:
  expression: histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter="source",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (le, source_workload, source_workload_namespace)) 
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]' 
  description: the client request duration that quantile is 0.9 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-client-request-duration-p99
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: client-request-duration
    source: rancher-istio
spec:
  expression: histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter="source",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (le, source_workload, source_workload_namespace)) 
  legendFormat: 'source:[[source_workload]].[[source_workload_namespace]]' 
  description: the client request duration that quantile is 0.99 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-4xxs-service
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 4xxs-service
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", destination_service_name=~"$service",response_code=~"4.*", destination_service_namespace=~"$namespace"}[5m])) by ( destination_workload, destination_workload_namespace)
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]'
  description: the count of service requests that response code is 4xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-5xxs-service
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 5xxs-service
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", destination_service_name=~"$service",response_code=~"5.*", destination_service_namespace=~"$namespace"}[5m])) by (destination_workload, destination_workload_namespace)
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]'
  description: the count of service requests that response code is 5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-success-service
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-success-rate-service
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", destination_service_name=~"$service",response_code!~"5.*", destination_service_namespace=~"$namespace"}[5m])) by (destination_workload, destination_workload_namespace) / sum(irate(istio_requests_total{reporter="destination", destination_service_name=~"$service", destination_service_namespace=~"$namespace"}[5m])) by (destination_workload, destination_workload_namespace) 
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]'
  description: the count of service requests that response code is non-5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-volume-service
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-request-volume-service
    source: rancher-istio
spec:
  expression: round(sum(irate(istio_requests_total{reporter="destination",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (destination_workload, destination_workload_namespace), 0.001) 
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]'
  description: the global service request volume 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-server-request-duration-p50
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: server-request-duration 
    source: rancher-istio
spec:
  expression: histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter="destination",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (le, destination_workload, destination_workload_namespace)) 
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]' 
  description: the server request duration that quantile is 0.5 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-server-request-duration-p90
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: server-request-duration 
    source: rancher-istio
spec:
  expression: histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter="destination",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (le, destination_workload, destination_workload_namespace)) 
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]' 
  description: the server request duration that quantile is 0.9 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-server-request-duration-p99
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: server-request-duration 
    source: rancher-istio
spec:
  expression: histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter="destination",destination_service_name=~"$service",destination_service_namespace=~"$namespace"}[1m])) by (le, destination_workload, destination_workload_namespace)) 
  legendFormat: 'destination:[[destination_workload]].[[destination_workload_namespace]]' 
  description: the server request duration that quantile is 0.99 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-4xxs-project
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 4xxs
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace", response_code=~"4.*"}[1m])) 
  legendFormat: 4xx request count 
  description: the count of requests that response code is 4xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-4xxs-project-top10
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 4xxs-top10
    source: rancher-istio
spec:
  expression: topk(10, sum(irate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace", response_code=~"4.*"}[1m])) by (destination_service_name, destination_service_namespace)) 
  legendFormat: '[[destination_service_name]].[[destination_service_namespace]]' 
  description: the count of requests that response code is 4xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-5xxs-project
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 5xxs
    source: rancher-istio
spec:
  expression: sum(irate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace", response_code=~"5.*"}[1m])) 
  legendFormat: 5xx request count 
  description: the count of requests that response code is 5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-5xxs-project-top10
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: 5xxs-top10
    source: rancher-istio
spec:
  expression: topk(10, sum(irate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace",  response_code=~"5.*"}[1m])) by (destination_service_name, destination_service_namespace)) 
  legendFormat: '[[destination_service_name]].[[destination_service_namespace]]' 
  description: the count of requests that response code is 5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-requests-total-success-project
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-success-rate
    source: rancher-istio
spec:
  expression: sum(rate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace", response_code!~"5.*"}[1m])) / sum(rate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace"}[1m])) 
  legendFormat: Success rate 
  description: the count of requests that response code is non-5xx
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-global-requests-volume-project
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-request-volume
    source: rancher-istio
spec:
  expression: round(sum(irate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace"}[1m])), 0.001) 
  legendFormat: Request volume 
  description: the global request of volume 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-global-requests-volume-project-top10
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: global-request-volume-top10
    source: rancher-istio
spec:
  expression: topk(10, round(sum(irate(istio_requests_total{reporter="destination", destination_service_namespace=~"$namespace"}[1m])) by (destination_service_name, destination_service_namespace) , 0.001)) 
  legendFormat: '[[destination_service_name]].[[destination_service_namespace]]' 
  description: the global request of volume 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-request-duration-p50-top10
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: request-duration-p50-top10 
    source: rancher-istio
spec:
  expression: topk(10, histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter="destination", destination_service_namespace=~"$namespace"}[1m])) by (le, destination_service_name, destination_service_namespace))) 
  legendFormat: '[[destination_service_name]].[[destination_service_namespace]]' 
  description: the server request duration that quantile is 0.5 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-request-duration-p90-top10
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: request-duration-p90-top10 
    source: rancher-istio
spec:
  expression: topk(10, histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter="destination", destination_service_namespace=~"$namespace"}[1m])) by (le, destination_service_name, destination_service_namespace))) 
  legendFormat: '[[destination_service_name]].[[destination_service_namespace]]' 
  description: the server request duration that quantile is 0.9 
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: istio-request-duration-p99-top10
  labels:
    app: metric-expression
    component: istio
    details: "false"
    level: project 
    graph: request-duration-p99-top10 
    source: rancher-istio
spec:
  expression: topk(10, histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter="destination", destination_service_namespace=~"$namespace"}[1m])) by (le, destination_service_name, destination_service_namespace))) 
  legendFormat: '[[destination_service_name]].[[destination_service_namespace]]' 
  description: the server request duration that quantile is 0.99 
---
`

	IstioClusterGraphTemplate = `
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-monitoring
    level: cluster
    component: istio
  name: istio-4xxs
spec:
  resourceType: istiocluster
  priority: 800
  title: istio-4xxs
  metricsSelector:
    details: "false"
    component: istio
    graph: 4xxs
    level: cluster
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 4xxs-details
    level: cluster
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-monitoring
    level: cluster
    component: istio
  name: istio-5xxs
spec:
  resourceType: istiocluster
  priority: 800
  title: istio-5xxs
  metricsSelector:
    details: "false"
    component: istio
    graph: 5xxs
    level: cluster
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 5xxs-details
    level: cluster
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-monitoring
    level: cluster
    component: istio
  name: istio-global-success-rate
spec:
  resourceType: istiocluster
  priority: 800
  title: istio-global-success-rate
  metricsSelector:
    details: "false"
    component: istio
    graph: global-success-rate
    level: cluster
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-success-rate-details
    level: cluster
---
`

	IstioProjectGraphTemplate = `
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-4xxs-client
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-4xxs-client
  metricsSelector:
    details: "false"
    component: istio
    graph: 4xxs-client
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 4xxs-details-client
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-5xxs-client
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-5xxs-client
  metricsSelector:
    details: "false"
    component: istio
    graph: 5xxs-client
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 5xxs-details-client
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-success-rate-client
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-global-success-rate-client
  metricsSelector:
    details: "false"
    component: istio
    graph: global-success-rate-client
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-success-rate-details-client
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-request-volume-client
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-global-request-volume-client
  metricsSelector:
    details: "false"
    component: istio
    graph: global-request-volume-client
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-request-volume-details-client
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-client-request-duration
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-client-request-duration
  metricsSelector:
    details: "false"
    component: istio
    graph: client-request-duration
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: client-request-duration-details
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-4xxs-service
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-4xxs-service
  metricsSelector:
    details: "false"
    component: istio
    graph: 4xxs-service
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 4xxs-details-service
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-5xxs-service
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-5xxs-service
  metricsSelector:
    details: "false"
    component: istio
    graph: 5xxs-service
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 5xxs-details-service
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-success-rate-service
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-global-success-rate-service
  metricsSelector:
    details: "false"
    component: istio
    graph: global-success-rate-service
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-success-rate-details-service
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-request-volume-service
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-global-request-volume-service
  metricsSelector:
    details: "false"
    component: istio
    graph: global-request-volume-service
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-request-volume-details-service
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-server-request-duration
spec:
  resourceType: istioservice 
  priority: 800
  title: istio-server-request-duration
  metricsSelector:
    details: "false"
    component: istio
    graph: server-request-duration
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: server-request-duration-details
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-4xxs
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-4xxs
  metricsSelector:
    details: "false"
    component: istio
    graph: 4xxs
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 4xxs-details
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-4xxs-top10
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-4xxs-top10
  metricsSelector:
    details: "false"
    component: istio
    graph: 4xxs-top10
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 4xxs-details-top10
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-5xxs
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-5xxs
  metricsSelector:
    details: "false"
    component: istio
    graph: 5xxs
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 5xxs-details
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project 
    component: istio
  name: istio-5xxs-top10
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-5xxs-top10
  metricsSelector:
    details: "false"
    component: istio
    graph: 5xxs-top10
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: 5xxs-details-top10
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-success-rate
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-global-success-rate
  metricsSelector:
    details: "false"
    component: istio
    graph: global-success-rate
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-success-rate-details
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-request-volume
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-global-request-volume
  metricsSelector:
    details: "false"
    component: istio
    graph: global-request-volume
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-request-volume-details
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-global-request-volume-top10
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-global-request-volume-top10
  metricsSelector:
    details: "false"
    component: istio
    graph: global-request-volume-top10
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: global-request-volume-details-top10
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-request-duration-p50-top10
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-request-duration-p50-top10
  metricsSelector:
    details: "false"
    component: istio
    graph: request-duration-p50-top10
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: request-duration-details-p50-top10
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-request-duration-p90-top10
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-request-duration-p90-top10
  metricsSelector:
    details: "false"
    component: istio
    graph: request-duration-p90-top10
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: request-duration-details-p90-top10
    level: project
---
apiVersion: management.cattle.io/v3
kind: ClusterMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-istio
    level: project
    component: istio
  name: istio-request-duration-p99-top10
spec:
  resourceType: istioproject 
  priority: 800
  title: istio-request-duration-p99-top10
  metricsSelector:
    details: "false"
    component: istio
    graph: request-duration-p99-top10
    level: project
  detailsMetricsSelector:
    details: "true"
    component: istio
    graph: request-duration-details-p99-top10
    level: project
---
`
)
