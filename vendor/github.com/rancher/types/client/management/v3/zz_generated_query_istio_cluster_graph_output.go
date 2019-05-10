package client

const (
	QueryIstioClusterGraphOutputType      = "queryIstioClusterGraphOutput"
	QueryIstioClusterGraphOutputFieldData = "data"
	QueryIstioClusterGraphOutputFieldType = "type"
)

type QueryIstioClusterGraphOutput struct {
	Data []QueryIstioClusterGraph `json:"data,omitempty" yaml:"data,omitempty"`
	Type string                   `json:"type,omitempty" yaml:"type,omitempty"`
}
