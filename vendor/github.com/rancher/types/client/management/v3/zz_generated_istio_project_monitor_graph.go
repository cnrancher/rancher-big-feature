package client

import (
	"github.com/rancher/norman/types"
)

const (
	IstioProjectMonitorGraphType                        = "istioProjectMonitorGraph"
	IstioProjectMonitorGraphFieldAnnotations            = "annotations"
	IstioProjectMonitorGraphFieldCreated                = "created"
	IstioProjectMonitorGraphFieldCreatorID              = "creatorId"
	IstioProjectMonitorGraphFieldDescription            = "description"
	IstioProjectMonitorGraphFieldDetailsMetricsSelector = "detailsMetricsSelector"
	IstioProjectMonitorGraphFieldDisplayResourceType    = "displayResourceType"
	IstioProjectMonitorGraphFieldGraphType              = "graphType"
	IstioProjectMonitorGraphFieldLabels                 = "labels"
	IstioProjectMonitorGraphFieldMetricsSelector        = "metricsSelector"
	IstioProjectMonitorGraphFieldName                   = "name"
	IstioProjectMonitorGraphFieldNamespaceId            = "namespaceId"
	IstioProjectMonitorGraphFieldOwnerReferences        = "ownerReferences"
	IstioProjectMonitorGraphFieldPriority               = "priority"
	IstioProjectMonitorGraphFieldProjectID              = "projectId"
	IstioProjectMonitorGraphFieldRemoved                = "removed"
	IstioProjectMonitorGraphFieldResourceType           = "resourceType"
	IstioProjectMonitorGraphFieldUUID                   = "uuid"
	IstioProjectMonitorGraphFieldYAxis                  = "yAxis"
)

type IstioProjectMonitorGraph struct {
	types.Resource
	Annotations            map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
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
	ProjectID              string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Removed                string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	UUID                   string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}

type IstioProjectMonitorGraphCollection struct {
	types.Collection
	Data   []IstioProjectMonitorGraph `json:"data,omitempty"`
	client *IstioProjectMonitorGraphClient
}

type IstioProjectMonitorGraphClient struct {
	apiClient *Client
}

type IstioProjectMonitorGraphOperations interface {
	List(opts *types.ListOpts) (*IstioProjectMonitorGraphCollection, error)
	Create(opts *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error)
	Update(existing *IstioProjectMonitorGraph, updates interface{}) (*IstioProjectMonitorGraph, error)
	Replace(existing *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error)
	ByID(id string) (*IstioProjectMonitorGraph, error)
	Delete(container *IstioProjectMonitorGraph) error

	CollectionActionQuery(resource *IstioProjectMonitorGraphCollection, input *QueryGraphInput) (*QueryIstioProjectGraphOutput, error)
}

func newIstioProjectMonitorGraphClient(apiClient *Client) *IstioProjectMonitorGraphClient {
	return &IstioProjectMonitorGraphClient{
		apiClient: apiClient,
	}
}

func (c *IstioProjectMonitorGraphClient) Create(container *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error) {
	resp := &IstioProjectMonitorGraph{}
	err := c.apiClient.Ops.DoCreate(IstioProjectMonitorGraphType, container, resp)
	return resp, err
}

func (c *IstioProjectMonitorGraphClient) Update(existing *IstioProjectMonitorGraph, updates interface{}) (*IstioProjectMonitorGraph, error) {
	resp := &IstioProjectMonitorGraph{}
	err := c.apiClient.Ops.DoUpdate(IstioProjectMonitorGraphType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *IstioProjectMonitorGraphClient) Replace(obj *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error) {
	resp := &IstioProjectMonitorGraph{}
	err := c.apiClient.Ops.DoReplace(IstioProjectMonitorGraphType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *IstioProjectMonitorGraphClient) List(opts *types.ListOpts) (*IstioProjectMonitorGraphCollection, error) {
	resp := &IstioProjectMonitorGraphCollection{}
	err := c.apiClient.Ops.DoList(IstioProjectMonitorGraphType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *IstioProjectMonitorGraphCollection) Next() (*IstioProjectMonitorGraphCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &IstioProjectMonitorGraphCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *IstioProjectMonitorGraphClient) ByID(id string) (*IstioProjectMonitorGraph, error) {
	resp := &IstioProjectMonitorGraph{}
	err := c.apiClient.Ops.DoByID(IstioProjectMonitorGraphType, id, resp)
	return resp, err
}

func (c *IstioProjectMonitorGraphClient) Delete(container *IstioProjectMonitorGraph) error {
	return c.apiClient.Ops.DoResourceDelete(IstioProjectMonitorGraphType, &container.Resource)
}

func (c *IstioProjectMonitorGraphClient) CollectionActionQuery(resource *IstioProjectMonitorGraphCollection, input *QueryGraphInput) (*QueryIstioProjectGraphOutput, error) {
	resp := &QueryIstioProjectGraphOutput{}
	err := c.apiClient.Ops.DoCollectionAction(IstioProjectMonitorGraphType, "query", &resource.Collection, input, resp)
	return resp, err
}
