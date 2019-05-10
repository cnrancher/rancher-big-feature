package v3

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	IstioClusterMonitorGraphGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "IstioClusterMonitorGraph",
	}
	IstioClusterMonitorGraphResource = metav1.APIResource{
		Name:         "istioclustermonitorgraphs",
		SingularName: "istioclustermonitorgraph",
		Namespaced:   true,

		Kind: IstioClusterMonitorGraphGroupVersionKind.Kind,
	}
)

func NewIstioClusterMonitorGraph(namespace, name string, obj IstioClusterMonitorGraph) *IstioClusterMonitorGraph {
	obj.APIVersion, obj.Kind = IstioClusterMonitorGraphGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type IstioClusterMonitorGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IstioClusterMonitorGraph
}

type IstioClusterMonitorGraphHandlerFunc func(key string, obj *IstioClusterMonitorGraph) (runtime.Object, error)

type IstioClusterMonitorGraphChangeHandlerFunc func(obj *IstioClusterMonitorGraph) (runtime.Object, error)

type IstioClusterMonitorGraphLister interface {
	List(namespace string, selector labels.Selector) (ret []*IstioClusterMonitorGraph, err error)
	Get(namespace, name string) (*IstioClusterMonitorGraph, error)
}

type IstioClusterMonitorGraphController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() IstioClusterMonitorGraphLister
	AddHandler(ctx context.Context, name string, handler IstioClusterMonitorGraphHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler IstioClusterMonitorGraphHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type IstioClusterMonitorGraphInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*IstioClusterMonitorGraph, error)
	Get(name string, opts metav1.GetOptions) (*IstioClusterMonitorGraph, error)
	Update(*IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*IstioClusterMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() IstioClusterMonitorGraphController
	AddHandler(ctx context.Context, name string, sync IstioClusterMonitorGraphHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle IstioClusterMonitorGraphLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync IstioClusterMonitorGraphHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle IstioClusterMonitorGraphLifecycle)
}

type istioClusterMonitorGraphLister struct {
	controller *istioClusterMonitorGraphController
}

func (l *istioClusterMonitorGraphLister) List(namespace string, selector labels.Selector) (ret []*IstioClusterMonitorGraph, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*IstioClusterMonitorGraph))
	})
	return
}

func (l *istioClusterMonitorGraphLister) Get(namespace, name string) (*IstioClusterMonitorGraph, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    IstioClusterMonitorGraphGroupVersionKind.Group,
			Resource: "istioClusterMonitorGraph",
		}, key)
	}
	return obj.(*IstioClusterMonitorGraph), nil
}

type istioClusterMonitorGraphController struct {
	controller.GenericController
}

func (c *istioClusterMonitorGraphController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *istioClusterMonitorGraphController) Lister() IstioClusterMonitorGraphLister {
	return &istioClusterMonitorGraphLister{
		controller: c,
	}
}

func (c *istioClusterMonitorGraphController) AddHandler(ctx context.Context, name string, handler IstioClusterMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*IstioClusterMonitorGraph); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *istioClusterMonitorGraphController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler IstioClusterMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*IstioClusterMonitorGraph); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type istioClusterMonitorGraphFactory struct {
}

func (c istioClusterMonitorGraphFactory) Object() runtime.Object {
	return &IstioClusterMonitorGraph{}
}

func (c istioClusterMonitorGraphFactory) List() runtime.Object {
	return &IstioClusterMonitorGraphList{}
}

func (s *istioClusterMonitorGraphClient) Controller() IstioClusterMonitorGraphController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.istioClusterMonitorGraphControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(IstioClusterMonitorGraphGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &istioClusterMonitorGraphController{
		GenericController: genericController,
	}

	s.client.istioClusterMonitorGraphControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type istioClusterMonitorGraphClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   IstioClusterMonitorGraphController
}

func (s *istioClusterMonitorGraphClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *istioClusterMonitorGraphClient) Create(o *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*IstioClusterMonitorGraph), err
}

func (s *istioClusterMonitorGraphClient) Get(name string, opts metav1.GetOptions) (*IstioClusterMonitorGraph, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*IstioClusterMonitorGraph), err
}

func (s *istioClusterMonitorGraphClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*IstioClusterMonitorGraph, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*IstioClusterMonitorGraph), err
}

func (s *istioClusterMonitorGraphClient) Update(o *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*IstioClusterMonitorGraph), err
}

func (s *istioClusterMonitorGraphClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *istioClusterMonitorGraphClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *istioClusterMonitorGraphClient) List(opts metav1.ListOptions) (*IstioClusterMonitorGraphList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*IstioClusterMonitorGraphList), err
}

func (s *istioClusterMonitorGraphClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *istioClusterMonitorGraphClient) Patch(o *IstioClusterMonitorGraph, patchType types.PatchType, data []byte, subresources ...string) (*IstioClusterMonitorGraph, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*IstioClusterMonitorGraph), err
}

func (s *istioClusterMonitorGraphClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *istioClusterMonitorGraphClient) AddHandler(ctx context.Context, name string, sync IstioClusterMonitorGraphHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *istioClusterMonitorGraphClient) AddLifecycle(ctx context.Context, name string, lifecycle IstioClusterMonitorGraphLifecycle) {
	sync := NewIstioClusterMonitorGraphLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *istioClusterMonitorGraphClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync IstioClusterMonitorGraphHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *istioClusterMonitorGraphClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle IstioClusterMonitorGraphLifecycle) {
	sync := NewIstioClusterMonitorGraphLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type IstioClusterMonitorGraphIndexer func(obj *IstioClusterMonitorGraph) ([]string, error)

type IstioClusterMonitorGraphClientCache interface {
	Get(namespace, name string) (*IstioClusterMonitorGraph, error)
	List(namespace string, selector labels.Selector) ([]*IstioClusterMonitorGraph, error)

	Index(name string, indexer IstioClusterMonitorGraphIndexer)
	GetIndexed(name, key string) ([]*IstioClusterMonitorGraph, error)
}

type IstioClusterMonitorGraphClient interface {
	Create(*IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error)
	Get(namespace, name string, opts metav1.GetOptions) (*IstioClusterMonitorGraph, error)
	Update(*IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*IstioClusterMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() IstioClusterMonitorGraphClientCache

	OnCreate(ctx context.Context, name string, sync IstioClusterMonitorGraphChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync IstioClusterMonitorGraphChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync IstioClusterMonitorGraphChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() IstioClusterMonitorGraphInterface
}

type istioClusterMonitorGraphClientCache struct {
	client *istioClusterMonitorGraphClient2
}

type istioClusterMonitorGraphClient2 struct {
	iface      IstioClusterMonitorGraphInterface
	controller IstioClusterMonitorGraphController
}

func (n *istioClusterMonitorGraphClient2) Interface() IstioClusterMonitorGraphInterface {
	return n.iface
}

func (n *istioClusterMonitorGraphClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *istioClusterMonitorGraphClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *istioClusterMonitorGraphClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *istioClusterMonitorGraphClient2) Create(obj *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error) {
	return n.iface.Create(obj)
}

func (n *istioClusterMonitorGraphClient2) Get(namespace, name string, opts metav1.GetOptions) (*IstioClusterMonitorGraph, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *istioClusterMonitorGraphClient2) Update(obj *IstioClusterMonitorGraph) (*IstioClusterMonitorGraph, error) {
	return n.iface.Update(obj)
}

func (n *istioClusterMonitorGraphClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *istioClusterMonitorGraphClient2) List(namespace string, opts metav1.ListOptions) (*IstioClusterMonitorGraphList, error) {
	return n.iface.List(opts)
}

func (n *istioClusterMonitorGraphClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *istioClusterMonitorGraphClientCache) Get(namespace, name string) (*IstioClusterMonitorGraph, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *istioClusterMonitorGraphClientCache) List(namespace string, selector labels.Selector) ([]*IstioClusterMonitorGraph, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *istioClusterMonitorGraphClient2) Cache() IstioClusterMonitorGraphClientCache {
	n.loadController()
	return &istioClusterMonitorGraphClientCache{
		client: n,
	}
}

func (n *istioClusterMonitorGraphClient2) OnCreate(ctx context.Context, name string, sync IstioClusterMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &istioClusterMonitorGraphLifecycleDelegate{create: sync})
}

func (n *istioClusterMonitorGraphClient2) OnChange(ctx context.Context, name string, sync IstioClusterMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &istioClusterMonitorGraphLifecycleDelegate{update: sync})
}

func (n *istioClusterMonitorGraphClient2) OnRemove(ctx context.Context, name string, sync IstioClusterMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &istioClusterMonitorGraphLifecycleDelegate{remove: sync})
}

func (n *istioClusterMonitorGraphClientCache) Index(name string, indexer IstioClusterMonitorGraphIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*IstioClusterMonitorGraph); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *istioClusterMonitorGraphClientCache) GetIndexed(name, key string) ([]*IstioClusterMonitorGraph, error) {
	var result []*IstioClusterMonitorGraph
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*IstioClusterMonitorGraph); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *istioClusterMonitorGraphClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type istioClusterMonitorGraphLifecycleDelegate struct {
	create IstioClusterMonitorGraphChangeHandlerFunc
	update IstioClusterMonitorGraphChangeHandlerFunc
	remove IstioClusterMonitorGraphChangeHandlerFunc
}

func (n *istioClusterMonitorGraphLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *istioClusterMonitorGraphLifecycleDelegate) Create(obj *IstioClusterMonitorGraph) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *istioClusterMonitorGraphLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *istioClusterMonitorGraphLifecycleDelegate) Remove(obj *IstioClusterMonitorGraph) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *istioClusterMonitorGraphLifecycleDelegate) Updated(obj *IstioClusterMonitorGraph) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
