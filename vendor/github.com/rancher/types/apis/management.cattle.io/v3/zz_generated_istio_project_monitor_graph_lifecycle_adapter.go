package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type IstioProjectMonitorGraphLifecycle interface {
	Create(obj *IstioProjectMonitorGraph) (runtime.Object, error)
	Remove(obj *IstioProjectMonitorGraph) (runtime.Object, error)
	Updated(obj *IstioProjectMonitorGraph) (runtime.Object, error)
}

type istioProjectMonitorGraphLifecycleAdapter struct {
	lifecycle IstioProjectMonitorGraphLifecycle
}

func (w *istioProjectMonitorGraphLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *istioProjectMonitorGraphLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *istioProjectMonitorGraphLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*IstioProjectMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *istioProjectMonitorGraphLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*IstioProjectMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *istioProjectMonitorGraphLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*IstioProjectMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewIstioProjectMonitorGraphLifecycleAdapter(name string, clusterScoped bool, client IstioProjectMonitorGraphInterface, l IstioProjectMonitorGraphLifecycle) IstioProjectMonitorGraphHandlerFunc {
	adapter := &istioProjectMonitorGraphLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *IstioProjectMonitorGraph) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
