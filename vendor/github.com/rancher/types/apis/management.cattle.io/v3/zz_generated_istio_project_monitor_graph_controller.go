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
	IstioProjectMonitorGraphGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "IstioProjectMonitorGraph",
	}
	IstioProjectMonitorGraphResource = metav1.APIResource{
		Name:         "istioprojectmonitorgraphs",
		SingularName: "istioprojectmonitorgraph",
		Namespaced:   true,

		Kind: IstioProjectMonitorGraphGroupVersionKind.Kind,
	}
)

func NewIstioProjectMonitorGraph(namespace, name string, obj IstioProjectMonitorGraph) *IstioProjectMonitorGraph {
	obj.APIVersion, obj.Kind = IstioProjectMonitorGraphGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type IstioProjectMonitorGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IstioProjectMonitorGraph
}

type IstioProjectMonitorGraphHandlerFunc func(key string, obj *IstioProjectMonitorGraph) (runtime.Object, error)

type IstioProjectMonitorGraphChangeHandlerFunc func(obj *IstioProjectMonitorGraph) (runtime.Object, error)

type IstioProjectMonitorGraphLister interface {
	List(namespace string, selector labels.Selector) (ret []*IstioProjectMonitorGraph, err error)
	Get(namespace, name string) (*IstioProjectMonitorGraph, error)
}

type IstioProjectMonitorGraphController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() IstioProjectMonitorGraphLister
	AddHandler(ctx context.Context, name string, handler IstioProjectMonitorGraphHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler IstioProjectMonitorGraphHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type IstioProjectMonitorGraphInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*IstioProjectMonitorGraph, error)
	Get(name string, opts metav1.GetOptions) (*IstioProjectMonitorGraph, error)
	Update(*IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*IstioProjectMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() IstioProjectMonitorGraphController
	AddHandler(ctx context.Context, name string, sync IstioProjectMonitorGraphHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle IstioProjectMonitorGraphLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync IstioProjectMonitorGraphHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle IstioProjectMonitorGraphLifecycle)
}

type istioProjectMonitorGraphLister struct {
	controller *istioProjectMonitorGraphController
}

func (l *istioProjectMonitorGraphLister) List(namespace string, selector labels.Selector) (ret []*IstioProjectMonitorGraph, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*IstioProjectMonitorGraph))
	})
	return
}

func (l *istioProjectMonitorGraphLister) Get(namespace, name string) (*IstioProjectMonitorGraph, error) {
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
			Group:    IstioProjectMonitorGraphGroupVersionKind.Group,
			Resource: "istioProjectMonitorGraph",
		}, key)
	}
	return obj.(*IstioProjectMonitorGraph), nil
}

type istioProjectMonitorGraphController struct {
	controller.GenericController
}

func (c *istioProjectMonitorGraphController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *istioProjectMonitorGraphController) Lister() IstioProjectMonitorGraphLister {
	return &istioProjectMonitorGraphLister{
		controller: c,
	}
}

func (c *istioProjectMonitorGraphController) AddHandler(ctx context.Context, name string, handler IstioProjectMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*IstioProjectMonitorGraph); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *istioProjectMonitorGraphController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler IstioProjectMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*IstioProjectMonitorGraph); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type istioProjectMonitorGraphFactory struct {
}

func (c istioProjectMonitorGraphFactory) Object() runtime.Object {
	return &IstioProjectMonitorGraph{}
}

func (c istioProjectMonitorGraphFactory) List() runtime.Object {
	return &IstioProjectMonitorGraphList{}
}

func (s *istioProjectMonitorGraphClient) Controller() IstioProjectMonitorGraphController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.istioProjectMonitorGraphControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(IstioProjectMonitorGraphGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &istioProjectMonitorGraphController{
		GenericController: genericController,
	}

	s.client.istioProjectMonitorGraphControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type istioProjectMonitorGraphClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   IstioProjectMonitorGraphController
}

func (s *istioProjectMonitorGraphClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *istioProjectMonitorGraphClient) Create(o *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*IstioProjectMonitorGraph), err
}

func (s *istioProjectMonitorGraphClient) Get(name string, opts metav1.GetOptions) (*IstioProjectMonitorGraph, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*IstioProjectMonitorGraph), err
}

func (s *istioProjectMonitorGraphClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*IstioProjectMonitorGraph, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*IstioProjectMonitorGraph), err
}

func (s *istioProjectMonitorGraphClient) Update(o *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*IstioProjectMonitorGraph), err
}

func (s *istioProjectMonitorGraphClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *istioProjectMonitorGraphClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *istioProjectMonitorGraphClient) List(opts metav1.ListOptions) (*IstioProjectMonitorGraphList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*IstioProjectMonitorGraphList), err
}

func (s *istioProjectMonitorGraphClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *istioProjectMonitorGraphClient) Patch(o *IstioProjectMonitorGraph, patchType types.PatchType, data []byte, subresources ...string) (*IstioProjectMonitorGraph, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*IstioProjectMonitorGraph), err
}

func (s *istioProjectMonitorGraphClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *istioProjectMonitorGraphClient) AddHandler(ctx context.Context, name string, sync IstioProjectMonitorGraphHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *istioProjectMonitorGraphClient) AddLifecycle(ctx context.Context, name string, lifecycle IstioProjectMonitorGraphLifecycle) {
	sync := NewIstioProjectMonitorGraphLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *istioProjectMonitorGraphClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync IstioProjectMonitorGraphHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *istioProjectMonitorGraphClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle IstioProjectMonitorGraphLifecycle) {
	sync := NewIstioProjectMonitorGraphLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type IstioProjectMonitorGraphIndexer func(obj *IstioProjectMonitorGraph) ([]string, error)

type IstioProjectMonitorGraphClientCache interface {
	Get(namespace, name string) (*IstioProjectMonitorGraph, error)
	List(namespace string, selector labels.Selector) ([]*IstioProjectMonitorGraph, error)

	Index(name string, indexer IstioProjectMonitorGraphIndexer)
	GetIndexed(name, key string) ([]*IstioProjectMonitorGraph, error)
}

type IstioProjectMonitorGraphClient interface {
	Create(*IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error)
	Get(namespace, name string, opts metav1.GetOptions) (*IstioProjectMonitorGraph, error)
	Update(*IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*IstioProjectMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() IstioProjectMonitorGraphClientCache

	OnCreate(ctx context.Context, name string, sync IstioProjectMonitorGraphChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync IstioProjectMonitorGraphChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync IstioProjectMonitorGraphChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() IstioProjectMonitorGraphInterface
}

type istioProjectMonitorGraphClientCache struct {
	client *istioProjectMonitorGraphClient2
}

type istioProjectMonitorGraphClient2 struct {
	iface      IstioProjectMonitorGraphInterface
	controller IstioProjectMonitorGraphController
}

func (n *istioProjectMonitorGraphClient2) Interface() IstioProjectMonitorGraphInterface {
	return n.iface
}

func (n *istioProjectMonitorGraphClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *istioProjectMonitorGraphClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *istioProjectMonitorGraphClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *istioProjectMonitorGraphClient2) Create(obj *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error) {
	return n.iface.Create(obj)
}

func (n *istioProjectMonitorGraphClient2) Get(namespace, name string, opts metav1.GetOptions) (*IstioProjectMonitorGraph, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *istioProjectMonitorGraphClient2) Update(obj *IstioProjectMonitorGraph) (*IstioProjectMonitorGraph, error) {
	return n.iface.Update(obj)
}

func (n *istioProjectMonitorGraphClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *istioProjectMonitorGraphClient2) List(namespace string, opts metav1.ListOptions) (*IstioProjectMonitorGraphList, error) {
	return n.iface.List(opts)
}

func (n *istioProjectMonitorGraphClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *istioProjectMonitorGraphClientCache) Get(namespace, name string) (*IstioProjectMonitorGraph, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *istioProjectMonitorGraphClientCache) List(namespace string, selector labels.Selector) ([]*IstioProjectMonitorGraph, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *istioProjectMonitorGraphClient2) Cache() IstioProjectMonitorGraphClientCache {
	n.loadController()
	return &istioProjectMonitorGraphClientCache{
		client: n,
	}
}

func (n *istioProjectMonitorGraphClient2) OnCreate(ctx context.Context, name string, sync IstioProjectMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &istioProjectMonitorGraphLifecycleDelegate{create: sync})
}

func (n *istioProjectMonitorGraphClient2) OnChange(ctx context.Context, name string, sync IstioProjectMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &istioProjectMonitorGraphLifecycleDelegate{update: sync})
}

func (n *istioProjectMonitorGraphClient2) OnRemove(ctx context.Context, name string, sync IstioProjectMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &istioProjectMonitorGraphLifecycleDelegate{remove: sync})
}

func (n *istioProjectMonitorGraphClientCache) Index(name string, indexer IstioProjectMonitorGraphIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*IstioProjectMonitorGraph); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *istioProjectMonitorGraphClientCache) GetIndexed(name, key string) ([]*IstioProjectMonitorGraph, error) {
	var result []*IstioProjectMonitorGraph
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*IstioProjectMonitorGraph); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *istioProjectMonitorGraphClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type istioProjectMonitorGraphLifecycleDelegate struct {
	create IstioProjectMonitorGraphChangeHandlerFunc
	update IstioProjectMonitorGraphChangeHandlerFunc
	remove IstioProjectMonitorGraphChangeHandlerFunc
}

func (n *istioProjectMonitorGraphLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *istioProjectMonitorGraphLifecycleDelegate) Create(obj *IstioProjectMonitorGraph) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *istioProjectMonitorGraphLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *istioProjectMonitorGraphLifecycleDelegate) Remove(obj *IstioProjectMonitorGraph) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *istioProjectMonitorGraphLifecycleDelegate) Updated(obj *IstioProjectMonitorGraph) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
