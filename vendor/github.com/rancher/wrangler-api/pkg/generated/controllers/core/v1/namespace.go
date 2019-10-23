/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type NamespaceHandler func(string, *v1.Namespace) (*v1.Namespace, error)

type NamespaceController interface {
	generic.ControllerMeta
	NamespaceClient

	OnChange(ctx context.Context, name string, sync NamespaceHandler)
	OnRemove(ctx context.Context, name string, sync NamespaceHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() NamespaceCache
}

type NamespaceClient interface {
	Create(*v1.Namespace) (*v1.Namespace, error)
	Update(*v1.Namespace) (*v1.Namespace, error)
	UpdateStatus(*v1.Namespace) (*v1.Namespace, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Namespace, error)
	List(opts metav1.ListOptions) (*v1.NamespaceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

type NamespaceCache interface {
	Get(name string) (*v1.Namespace, error)
	List(selector labels.Selector) ([]*v1.Namespace, error)

	AddIndexer(indexName string, indexer NamespaceIndexer)
	GetByIndex(indexName, key string) ([]*v1.Namespace, error)
}

type NamespaceIndexer func(obj *v1.Namespace) ([]string, error)

type namespaceController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.NamespacesGetter
	informer          informers.NamespaceInformer
	gvk               schema.GroupVersionKind
}

func NewNamespaceController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.NamespacesGetter, informer informers.NamespaceInformer) NamespaceController {
	return &namespaceController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromNamespaceHandlerToHandler(sync NamespaceHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Namespace
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Namespace))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *namespaceController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Namespace))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateNamespaceDeepCopyOnChange(client NamespaceClient, obj *v1.Namespace, handler func(obj *v1.Namespace) (*v1.Namespace, error)) (*v1.Namespace, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *namespaceController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *namespaceController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *namespaceController) OnChange(ctx context.Context, name string, sync NamespaceHandler) {
	c.AddGenericHandler(ctx, name, FromNamespaceHandlerToHandler(sync))
}

func (c *namespaceController) OnRemove(ctx context.Context, name string, sync NamespaceHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromNamespaceHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *namespaceController) Enqueue(name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), "", name)
}

func (c *namespaceController) EnqueueAfter(name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), "", name, duration)
}

func (c *namespaceController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *namespaceController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *namespaceController) Cache() NamespaceCache {
	return &namespaceCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *namespaceController) Create(obj *v1.Namespace) (*v1.Namespace, error) {
	return c.clientGetter.Namespaces().Create(obj)
}

func (c *namespaceController) Update(obj *v1.Namespace) (*v1.Namespace, error) {
	return c.clientGetter.Namespaces().Update(obj)
}

func (c *namespaceController) UpdateStatus(obj *v1.Namespace) (*v1.Namespace, error) {
	return c.clientGetter.Namespaces().UpdateStatus(obj)
}

func (c *namespaceController) Delete(name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Namespaces().Delete(name, options)
}

func (c *namespaceController) Get(name string, options metav1.GetOptions) (*v1.Namespace, error) {
	return c.clientGetter.Namespaces().Get(name, options)
}

func (c *namespaceController) List(opts metav1.ListOptions) (*v1.NamespaceList, error) {
	return c.clientGetter.Namespaces().List(opts)
}

func (c *namespaceController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Namespaces().Watch(opts)
}

func (c *namespaceController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error) {
	return c.clientGetter.Namespaces().Patch(name, pt, data, subresources...)
}

type namespaceCache struct {
	lister  listers.NamespaceLister
	indexer cache.Indexer
}

func (c *namespaceCache) Get(name string) (*v1.Namespace, error) {
	return c.lister.Get(name)
}

func (c *namespaceCache) List(selector labels.Selector) ([]*v1.Namespace, error) {
	return c.lister.List(selector)
}

func (c *namespaceCache) AddIndexer(indexName string, indexer NamespaceIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Namespace))
		},
	}))
}

func (c *namespaceCache) GetByIndex(indexName, key string) (result []*v1.Namespace, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.Namespace))
	}
	return result, nil
}

type NamespaceStatusHandler func(obj *v1.Namespace, status v1.NamespaceStatus) (v1.NamespaceStatus, error)

type NamespaceGeneratingHandler func(obj *v1.Namespace, status v1.NamespaceStatus) ([]runtime.Object, v1.NamespaceStatus, error)

func RegisterNamespaceStatusHandler(ctx context.Context, controller NamespaceController, condition condition.Cond, name string, handler NamespaceStatusHandler) {
	statusHandler := &namespaceStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromNamespaceHandlerToHandler(statusHandler.sync))
}

func RegisterNamespaceGeneratingHandler(ctx context.Context, controller NamespaceController, apply apply.Apply,
	condition condition.Cond, name string, handler NamespaceGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &namespaceGeneratingHandler{
		NamespaceGeneratingHandler: handler,
		apply:                      apply,
		name:                       name,
		gvk:                        controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	RegisterNamespaceStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type namespaceStatusHandler struct {
	client    NamespaceClient
	condition condition.Cond
	handler   NamespaceStatusHandler
}

func (a *namespaceStatusHandler) sync(key string, obj *v1.Namespace) (*v1.Namespace, error) {
	if obj == nil {
		return obj, nil
	}

	status := obj.Status
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *status.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(obj, "", nil)
		} else {
			a.condition.SetError(obj, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(status, newStatus) {
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type namespaceGeneratingHandler struct {
	NamespaceGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *namespaceGeneratingHandler) Handle(obj *v1.Namespace, status v1.NamespaceStatus) (v1.NamespaceStatus, error) {
	objs, newStatus, err := a.NamespaceGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	apply := a.apply

	if !a.opts.DynamicLookup {
		apply = apply.WithStrictCaching()
	}

	if !a.opts.AllowCrossNamespace && !a.opts.AllowClusterScoped {
		apply = apply.WithSetOwnerReference(true, false).
			WithDefaultNamespace(obj.GetNamespace()).
			WithListerNamespace(obj.GetNamespace())
	}

	if !a.opts.AllowClusterScoped {
		apply = apply.WithRestrictClusterScoped()
	}

	return newStatus, apply.
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
