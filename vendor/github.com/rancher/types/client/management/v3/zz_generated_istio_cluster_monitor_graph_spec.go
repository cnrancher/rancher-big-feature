package client

const (
	IstioClusterMonitorGraphSpecType                        = "istioClusterMonitorGraphSpec"
	IstioClusterMonitorGraphSpecFieldClusterID              = "clusterId"
	IstioClusterMonitorGraphSpecFieldDescription            = "description"
	IstioClusterMonitorGraphSpecFieldDetailsMetricsSelector = "detailsMetricsSelector"
	IstioClusterMonitorGraphSpecFieldDisplayResourceType    = "displayResourceType"
	IstioClusterMonitorGraphSpecFieldGraphType              = "graphType"
	IstioClusterMonitorGraphSpecFieldMetricsSelector        = "metricsSelector"
	IstioClusterMonitorGraphSpecFieldPriority               = "priority"
	IstioClusterMonitorGraphSpecFieldResourceType           = "resourceType"
	IstioClusterMonitorGraphSpecFieldYAxis                  = "yAxis"
)

type IstioClusterMonitorGraphSpec struct {
	ClusterID              string            `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	Description            string            `json:"description,omitempty" yaml:"description,omitempty"`
	DetailsMetricsSelector map[string]string `json:"detailsMetricsSelector,omitempty" yaml:"detailsMetricsSelector,omitempty"`
	DisplayResourceType    string            `json:"displayResourceType,omitempty" yaml:"displayResourceType,omitempty"`
	GraphType              string            `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	MetricsSelector        map[string]string `json:"metricsSelector,omitempty" yaml:"metricsSelector,omitempty"`
	Priority               int64             `json:"priority,omitempty" yaml:"priority,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}
