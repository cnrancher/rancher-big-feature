package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type IstioClusterMonitorGraphLifecycle interface {
	Create(obj *IstioClusterMonitorGraph) (runtime.Object, error)
	Remove(obj *IstioClusterMonitorGraph) (runtime.Object, error)
	Updated(obj *IstioClusterMonitorGraph) (runtime.Object, error)
}

type istioClusterMonitorGraphLifecycleAdapter struct {
	lifecycle IstioClusterMonitorGraphLifecycle
}

func (w *istioClusterMonitorGraphLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *istioClusterMonitorGraphLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *istioClusterMonitorGraphLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*IstioClusterMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *istioClusterMonitorGraphLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*IstioClusterMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *istioClusterMonitorGraphLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*IstioClusterMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewIstioClusterMonitorGraphLifecycleAdapter(name string, clusterScoped bool, client IstioClusterMonitorGraphInterface, l IstioClusterMonitorGraphLifecycle) IstioClusterMonitorGraphHandlerFunc {
	adapter := &istioClusterMonitorGraphLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *IstioClusterMonitorGraph) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
