package client

const (
	IstioProjectMonitorGraphSpecType                        = "istioProjectMonitorGraphSpec"
	IstioProjectMonitorGraphSpecFieldDescription            = "description"
	IstioProjectMonitorGraphSpecFieldDetailsMetricsSelector = "detailsMetricsSelector"
	IstioProjectMonitorGraphSpecFieldDisplayResourceType    = "displayResourceType"
	IstioProjectMonitorGraphSpecFieldGraphType              = "graphType"
	IstioProjectMonitorGraphSpecFieldMetricsSelector        = "metricsSelector"
	IstioProjectMonitorGraphSpecFieldPriority               = "priority"
	IstioProjectMonitorGraphSpecFieldProjectID              = "projectId"
	IstioProjectMonitorGraphSpecFieldResourceType           = "resourceType"
	IstioProjectMonitorGraphSpecFieldYAxis                  = "yAxis"
)

type IstioProjectMonitorGraphSpec struct {
	Description            string            `json:"description,omitempty" yaml:"description,omitempty"`
	DetailsMetricsSelector map[string]string `json:"detailsMetricsSelector,omitempty" yaml:"detailsMetricsSelector,omitempty"`
	DisplayResourceType    string            `json:"displayResourceType,omitempty" yaml:"displayResourceType,omitempty"`
	GraphType              string            `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	MetricsSelector        map[string]string `json:"metricsSelector,omitempty" yaml:"metricsSelector,omitempty"`
	Priority               int64             `json:"priority,omitempty" yaml:"priority,omitempty"`
	ProjectID              string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}
