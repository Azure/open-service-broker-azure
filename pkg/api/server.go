package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"context"

	"github.com/Azure/azure-service-broker/pkg/async"
	"github.com/Azure/azure-service-broker/pkg/crypto"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
	"github.com/gorilla/mux"
)

type errHTTPServerStopped struct {
	err error
}

func (e *errHTTPServerStopped) Error() string {
	if e.err == nil {
		return "http server stopped"
	}
	return fmt.Sprintf("http server stopped: %s", e.err)
}

// Server is an interface for components that respond to HTTP requests on behalf
// of the broker
type Server interface {
	// Start causes the api server to start serving HTTP requests. It will block
	// until an error occurs and will return that error.
	Start(context.Context) error
}

type server struct {
	port        int
	store       storage.Store
	asyncEngine async.Engine
	codec       crypto.Codec
	router      *mux.Router
	// Modules indexed by service
	modules map[string]service.Module
	// Provisioners indexed by service
	provisioners map[string]service.Provisioner
	// Deprovisioners indexed by service
	deprovisioners  map[string]service.Deprovisioner
	catalog         service.Catalog
	catalogResponse []byte
	// This allows tests to inject an alternative implementation of this function
	listenAndServe func(context.Context) error
}

// NewServer returns an HTTP router
func NewServer(
	port int,
	store storage.Store,
	asyncEngine async.Engine,
	codec crypto.Codec,
	modules map[string]service.Module,
	provisioners map[string]service.Provisioner,
	deprovisioners map[string]service.Deprovisioner,
) (Server, error) {
	s := &server{
		port:           port,
		store:          store,
		asyncEngine:    asyncEngine,
		codec:          codec,
		modules:        modules,
		provisioners:   provisioners,
		deprovisioners: deprovisioners,
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
	s.router = router

	services := []service.Service{}
	for _, module := range modules {
		catalog, err := module.GetCatalog()
		if err != nil {
			return nil, err
		}
		services = append(services, catalog.GetServices()...)
	}
	s.catalog = service.NewCatalog(services)
	catalogJSONStr, err := s.catalog.ToJSONString()
	if err != nil {
		return nil, err
	}
	s.catalogResponse = []byte(catalogJSONStr)

	s.listenAndServe = s.defaultListenAndServe

	return s, nil
}

func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	go func() {
		log.Printf("Listening on http://0.0.0.0:%d", s.port)
		err := s.listenAndServe(ctx)
		hss := &errHTTPServerStopped{err: err}
		select {
		case errChan <- hss:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (s *server) defaultListenAndServe(ctx context.Context) error {
	errChan := make(chan error)
	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}
	go func() {
		err := svr.ListenAndServe()
		select {
		case errChan <- err:
		case <-ctx.Done():
		}
	}()
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Second*5,
		)
		defer cancel()
		svr.Shutdown(shutdownCtx)
		return ctx.Err()
	}
}
