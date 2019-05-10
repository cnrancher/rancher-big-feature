package client

import (
	"github.com/rancher/norman/types"
)

const (
	IstioClusterMonitorGraphType                        = "istioClusterMonitorGraph"
	IstioClusterMonitorGraphFieldAnnotations            = "annotations"
	IstioClusterMonitorGraphFieldClusterID              = "clusterId"
	IstioClusterMonitorGraphFieldCreated                = "created"
	IstioClusterMonitorGraphFieldCreatorID              = "creatorId"
	IstioClusterMonitorGraphFieldDescription            = "description"
	IstioClusterMonitorGraphFieldDetailsMetricsSelector = "detailsMetricsSelector"
	IstioClusterMonitorGraphFieldDisplayResourceType    = "displayResourceType"
	IstioClusterMonitorGraphFieldGraphType              = "graphType"
	IstioClusterMonitorGraphFieldLabels                 = "labels"
	IstioClusterMonitorGraphFieldMetricsSelector        = "metricsSelector"
	IstioClusterMonitorGraphFieldName                   = "name"
	IstioClusterMonitorGraphFieldNamespaceId            = "namespaceId"
	IstioClusterMonitorGraphFieldOwnerReferences        = "ownerReferences"
	IstioClusterMonitorGraphFieldPriority               = "priority"
	IstioClusterMonitorGraphFieldRemoved                = "removed"
	IstioClusterMonitorGraphFieldResourceType           = "resourceType"
	IstioClusterMonitorGraphFieldUUID                   = "uuid"
	IstioClusterMonitorGraphFieldYAxis                  = "yAxis"
)

type IstioClusterMonitorGraph struct {
	types.Resource
	Annotations            map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClusterID              string            `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	Created                string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID              string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description            string            `json:"description,omitempty" yaml:"description,omitempty"`
	DetailsMetricsSelector map[string]string `json:"detailsMetricsSelector,omitempty" yaml:"detailsMetricsSelector,omitempty"`
	DisplayResourceType    string            `json:"displayResourceType,omitempty" yaml:"displayResourceType,omitempty"`
	GraphType              string            `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	MetricsSelector        map[string]string `json:"metricsSelector,omitempty" yaml:"metricsSelector,omitempty"`
	Name                   string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId            string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences        []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Priority               int64             `json:"priority,omitempty" yaml:"priority,omitempty"`
	Removed                string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	UUID                   string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}

type IstioClusterMonitorGraphCollection struct {
	types.Collection
	Data   []IstioClusterMonitorGraph `json:"data,omitempty"`
	client *IstioClusterMonitorGraphClient
}

type IstioClusterMonitorGraphClient struct {
	apiClient *Client
}

type IstioClusterMonitorGraphOperations interface {
	List(opts *types.ListOpts) (*IstioClusterMonitorGraphCollection, error)
	Create(opts *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error)
	Update(existing *IstioClusterMonitorGraph, updates interface{}) (*IstioClusterMonitorGraph, error)
	Replace(existing *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error)
	ByID(id string) (*IstioClusterMonitorGraph, error)
	Delete(container *IstioClusterMonitorGraph) error

	CollectionActionQuery(resource *IstioClusterMonitorGraphCollection, input *QueryGraphInput) (*QueryIstioClusterGraphOutput, error)
}

func newIstioClusterMonitorGraphClient(apiClient *Client) *IstioClusterMonitorGraphClient {
	return &IstioClusterMonitorGraphClient{
		apiClient: apiClient,
	}
}

func (c *IstioClusterMonitorGraphClient) Create(container *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error) {
	resp := &IstioClusterMonitorGraph{}
	err := c.apiClient.Ops.DoCreate(IstioClusterMonitorGraphType, container, resp)
	return resp, err
}

func (c *IstioClusterMonitorGraphClient) Update(existing *IstioClusterMonitorGraph, updates interface{}) (*IstioClusterMonitorGraph, error) {
	resp := &IstioClusterMonitorGraph{}
	err := c.apiClient.Ops.DoUpdate(IstioClusterMonitorGraphType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *IstioClusterMonitorGraphClient) Replace(obj *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error) {
	resp := &IstioClusterMonitorGraph{}
	err := c.apiClient.Ops.DoReplace(IstioClusterMonitorGraphType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *IstioClusterMonitorGraphClient) List(opts *types.ListOpts) (*IstioClusterMonitorGraphCollection, error) {
	resp := &IstioClusterMonitorGraphCollection{}
	err := c.apiClient.Ops.DoList(IstioClusterMonitorGraphType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *IstioClusterMonitorGraphCollection) Next() (*IstioClusterMonitorGraphCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &IstioClusterMonitorGraphCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *IstioClusterMonitorGraphClient) ByID(id string) (*IstioClusterMonitorGraph, error) {
	resp := &IstioClusterMonitorGraph{}
	err := c.apiClient.Ops.DoByID(IstioClusterMonitorGraphType, id, resp)
	return resp, err
}

func (c *IstioClusterMonitorGraphClient) Delete(container *IstioClusterMonitorGraph) error {
	return c.apiClient.Ops.DoResourceDelete(IstioClusterMonitorGraphType, &container.Resource)
}

func (c *IstioClusterMonitorGraphClient) CollectionActionQuery(resource *IstioClusterMonitorGraphCollection, input *QueryGraphInput) (*QueryIstioClusterGraphOutput, error) {
	resp := &QueryIstioClusterGraphOutput{}
	err := c.apiClient.Ops.DoCollectionAction(IstioClusterMonitorGraphType, "query", &resource.Collection, input, resp)
	return resp, err
}
