package client

const (
	QueryIstioClusterGraphType           = "queryIstioClusterGraph"
	QueryIstioClusterGraphFieldGraphName = "graphID"
	QueryIstioClusterGraphFieldSeries    = "series"
)

type QueryIstioClusterGraph struct {
	GraphName string   `json:"graphID,omitempty" yaml:"graphID,omitempty"`
	Series    []string `json:"series,omitempty" yaml:"series,omitempty"`
}
