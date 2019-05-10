package client

const (
	QueryIstioProjectGraphOutputType      = "queryIstioProjectGraphOutput"
	QueryIstioProjectGraphOutputFieldData = "data"
	QueryIstioProjectGraphOutputFieldType = "type"
)

type QueryIstioProjectGraphOutput struct {
	Data []QueryIstioProjectGraph `json:"data,omitempty" yaml:"data,omitempty"`
	Type string                   `json:"type,omitempty" yaml:"type,omitempty"`
}
