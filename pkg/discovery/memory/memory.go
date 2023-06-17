package memory

import (
	"context"
	"errors"
	"github.com/mamalmaleki/go-movie/pkg/discovery"
	"sync"
	"time"
)

//type serviceNameType string
//type instanceIDType string

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// Registry defines an in-memory service registry
type Registry struct {
	sync.RWMutex
	//serviceAddress map[serviceName]map[instanceID]*serviceInstance
	serviceAddress map[string]map[string]*serviceInstance
}

// Register creates a service record in the registry
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string,
	hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddress[serviceName]; !ok {
		r.serviceAddress[serviceName] = map[string]*serviceInstance{}
	}
	r.serviceAddress[serviceName][instanceID] = &serviceInstance{hostPort: hostPort,
		lastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddress[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddress[serviceName], instanceID)
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddress[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddress[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddress[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddress[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddress[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddress: map[string]map[string]*serviceInstance{}}
}
