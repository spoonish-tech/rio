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

package v1alpha3

import (
	"context"

	v1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	clientset "github.com/knative/pkg/client/clientset/versioned/typed/istio/v1alpha3"
	informers "github.com/knative/pkg/client/informers/externalversions/istio/v1alpha3"
	listers "github.com/knative/pkg/client/listers/istio/v1alpha3"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ServiceEntryHandler func(string, *v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error)

type ServiceEntryController interface {
	ServiceEntryClient

	OnChange(ctx context.Context, name string, sync ServiceEntryHandler)
	OnRemove(ctx context.Context, name string, sync ServiceEntryHandler)
	Enqueue(namespace, name string)

	Cache() ServiceEntryCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type ServiceEntryClient interface {
	Create(*v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error)
	Update(*v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha3.ServiceEntry, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha3.ServiceEntryList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.ServiceEntry, err error)
}

type ServiceEntryCache interface {
	Get(namespace, name string) (*v1alpha3.ServiceEntry, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha3.ServiceEntry, error)

	AddIndexer(indexName string, indexer ServiceEntryIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha3.ServiceEntry, error)
}

type ServiceEntryIndexer func(obj *v1alpha3.ServiceEntry) ([]string, error)

type serviceEntryController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.ServiceEntriesGetter
	informer          informers.ServiceEntryInformer
	gvk               schema.GroupVersionKind
}

func NewServiceEntryController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.ServiceEntriesGetter, informer informers.ServiceEntryInformer) ServiceEntryController {
	return &serviceEntryController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromServiceEntryHandlerToHandler(sync ServiceEntryHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha3.ServiceEntry
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha3.ServiceEntry))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *serviceEntryController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha3.ServiceEntry))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateServiceEntryOnChange(updater generic.Updater, handler ServiceEntryHandler) ServiceEntryHandler {
	return func(key string, obj *v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1alpha3.ServiceEntry)
			}
		}

		return copyObj, err
	}
}

func (c *serviceEntryController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *serviceEntryController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *serviceEntryController) OnChange(ctx context.Context, name string, sync ServiceEntryHandler) {
	c.AddGenericHandler(ctx, name, FromServiceEntryHandlerToHandler(sync))
}

func (c *serviceEntryController) OnRemove(ctx context.Context, name string, sync ServiceEntryHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromServiceEntryHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *serviceEntryController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *serviceEntryController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *serviceEntryController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *serviceEntryController) Cache() ServiceEntryCache {
	return &serviceEntryCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *serviceEntryController) Create(obj *v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error) {
	return c.clientGetter.ServiceEntries(obj.Namespace).Create(obj)
}

func (c *serviceEntryController) Update(obj *v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error) {
	return c.clientGetter.ServiceEntries(obj.Namespace).Update(obj)
}

func (c *serviceEntryController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.ServiceEntries(namespace).Delete(name, options)
}

func (c *serviceEntryController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha3.ServiceEntry, error) {
	return c.clientGetter.ServiceEntries(namespace).Get(name, options)
}

func (c *serviceEntryController) List(namespace string, opts metav1.ListOptions) (*v1alpha3.ServiceEntryList, error) {
	return c.clientGetter.ServiceEntries(namespace).List(opts)
}

func (c *serviceEntryController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.ServiceEntries(namespace).Watch(opts)
}

func (c *serviceEntryController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.ServiceEntry, err error) {
	return c.clientGetter.ServiceEntries(namespace).Patch(name, pt, data, subresources...)
}

type serviceEntryCache struct {
	lister  listers.ServiceEntryLister
	indexer cache.Indexer
}

func (c *serviceEntryCache) Get(namespace, name string) (*v1alpha3.ServiceEntry, error) {
	return c.lister.ServiceEntries(namespace).Get(name)
}

func (c *serviceEntryCache) List(namespace string, selector labels.Selector) ([]*v1alpha3.ServiceEntry, error) {
	return c.lister.ServiceEntries(namespace).List(selector)
}

func (c *serviceEntryCache) AddIndexer(indexName string, indexer ServiceEntryIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha3.ServiceEntry))
		},
	}))
}

func (c *serviceEntryCache) GetByIndex(indexName, key string) (result []*v1alpha3.ServiceEntry, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha3.ServiceEntry))
	}
	return result, nil
}
