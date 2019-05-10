package client

const (
	QueryIstioProjectGraphType           = "queryIstioProjectGraph"
	QueryIstioProjectGraphFieldGraphName = "graphID"
	QueryIstioProjectGraphFieldSeries    = "series"
)

type QueryIstioProjectGraph struct {
	GraphName string   `json:"graphID,omitempty" yaml:"graphID,omitempty"`
	Series    []string `json:"series,omitempty" yaml:"series,omitempty"`
}
