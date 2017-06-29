package api

import (
	"fmt"
	"log"
	"net/http"

	"context"

	"github.com/Azure/azure-service-broker/pkg/async"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
	"github.com/gorilla/mux"
)

// Server is an interface for components that respond to HTTP requests on behalf
// of the broker
type Server interface {
	// Start causes the api server to start serving HTTP requests. It will block
	// until an error occurs and will return that error.
	Start(context.Context) error
}

type server struct {
	port            int
	store           storage.Store
	asyncEngine     async.Engine
	router          *mux.Router
	modules         map[string]service.Module
	catalog         service.Catalog
	catalogResponse []byte
	// This allows tests to inject an alternative implementation of this function
	listenAndServe func(string, http.Handler) error
}

// NewServer returns an HTTP router
func NewServer(
	port int,
	store storage.Store,
	asyncEngine async.Engine,
	modules []service.Module,
) (Server, error) {
	s := &server{
		port:           port,
		store:          store,
		modules:        make(map[string]service.Module),
		asyncEngine:    asyncEngine,
		listenAndServe: listenAndServe,
	}
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc(
		"/v2/catalog",
		s.getCatalog,
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		s.provision,
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/last_operation",
		s.poll,
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		s.bind,
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		s.unbind,
	).Methods(http.MethodDelete)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		s.deprovision,
	).Methods(http.MethodDelete)
	// TODO: Delete this later-- this is just to aid in hacking
	router.HandleFunc(
		"/v2/test",
		s.test,
	).Methods(http.MethodGet)
	s.router = router

	services := []service.Service{}
	for _, module := range modules {
		catalog, err := module.GetCatalog()
		if err != nil {
			return nil, err
		}
		for _, svc := range catalog.GetServices() {
			existingModule, ok := s.modules[svc.GetID()]
			if ok {
				// This means we have more than one module claiming to provide services
				// with an ID in common. This is a SERIOUS problem.
				return nil, fmt.Errorf(
					"module %s and module %s BOTH provide a service with the id %s",
					existingModule.GetName(),
					module.GetName(),
					svc.GetID(),
				)
			}
			s.modules[svc.GetID()] = module
		}
		services = append(services, catalog.GetServices()...)
	}
	s.catalog = service.NewCatalog(services)
	catalogJSONStr, err := s.catalog.ToJSONString()
	if err != nil {
		return nil, err
	}
	s.catalogResponse = []byte(catalogJSONStr)

	return s, nil
}

func (s *server) Start(ctx context.Context) error {
	errChan := make(chan error)
	defer close(errChan)
	go func() {
		log.Printf("Listening on http://0.0.0.0:%d", s.port)
		// Start listening. This blocks until it errors or is interrupted.
		errChan <- s.listenAndServe(fmt.Sprintf(":%d", s.port), s.router)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

var listenAndServe = func(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
