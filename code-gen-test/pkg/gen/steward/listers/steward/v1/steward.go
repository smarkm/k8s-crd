/*
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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/smarkm/k8s-crd/code-gen-test/pkg/apis/steward/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// StewardLister helps list Stewards.
type StewardLister interface {
	// List lists all Stewards in the indexer.
	List(selector labels.Selector) (ret []*v1.Steward, err error)
	// Stewards returns an object that can list and get Stewards.
	Stewards(namespace string) StewardNamespaceLister
	StewardListerExpansion
}

// stewardLister implements the StewardLister interface.
type stewardLister struct {
	indexer cache.Indexer
}

// NewStewardLister returns a new StewardLister.
func NewStewardLister(indexer cache.Indexer) StewardLister {
	return &stewardLister{indexer: indexer}
}

// List lists all Stewards in the indexer.
func (s *stewardLister) List(selector labels.Selector) (ret []*v1.Steward, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Steward))
	})
	return ret, err
}

// Stewards returns an object that can list and get Stewards.
func (s *stewardLister) Stewards(namespace string) StewardNamespaceLister {
	return stewardNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// StewardNamespaceLister helps list and get Stewards.
type StewardNamespaceLister interface {
	// List lists all Stewards in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Steward, err error)
	// Get retrieves the Steward from the indexer for a given namespace and name.
	Get(name string) (*v1.Steward, error)
	StewardNamespaceListerExpansion
}

// stewardNamespaceLister implements the StewardNamespaceLister
// interface.
type stewardNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Stewards in the indexer for a given namespace.
func (s stewardNamespaceLister) List(selector labels.Selector) (ret []*v1.Steward, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Steward))
	})
	return ret, err
}

// Get retrieves the Steward from the indexer for a given namespace and name.
func (s stewardNamespaceLister) Get(name string) (*v1.Steward, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("steward"), name)
	}
	return obj.(*v1.Steward), nil
}
