package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-service-broker/pkg/api/authenticator"
	"github.com/Azure/azure-service-broker/pkg/async"
	"github.com/Azure/azure-service-broker/pkg/crypto"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
	log "github.com/Sirupsen/logrus"
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
	port          int
	store         storage.Store
	asyncEngine   async.Engine
	codec         crypto.Codec
	authenticator authenticator.Authenticator
	router        *mux.Router
	// Modules indexed by service
	modules         map[string]service.Module
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
	authenticator authenticator.Authenticator,
	modules map[string]service.Module,
) (Server, error) {
	s := &server{
		port:          port,
		store:         store,
		asyncEngine:   asyncEngine,
		codec:         codec,
		authenticator: authenticator,
		modules:       modules,
	}
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc(
		"/v2/catalog",
		s.authenticator.Authenticate(s.getCatalog),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		s.authenticator.Authenticate(s.provision),
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		s.authenticator.Authenticate(s.update),
	).Methods(http.MethodPatch)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/last_operation",
		s.authenticator.Authenticate(s.poll),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		s.authenticator.Authenticate(s.bind),
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		s.authenticator.Authenticate(s.unbind),
	).Methods(http.MethodDelete)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		s.authenticator.Authenticate(s.deprovision),
	).Methods(http.MethodDelete)
	router.HandleFunc(
		"/healthz",
		s.healthCheck, // No authentication on this request
	).Methods(http.MethodGet)
	s.router = router

	// modules is a map of modules indexed by service id. If a module provides
	// more than one service, we could end up iterating over a given module
	// more than once and adding all its services to the consolidated catalog
	// multiple times. To avoid this, we keep track of modules whose services
	// have already been added to the consolidated catalog in a handledModules
	// map that indexes struct{}{} (no allocation required) by module name. (This
	// is a poor man's set since Go lacks a dedicated set data structure.)
	services := []service.Service{}
	handledModules := map[string]struct{}{}
	var ok bool
	for _, module := range modules {
		if _, ok = handledModules[module.GetName()]; ok {
			continue
		}
		handledModules[module.GetName()] = struct{}{}
		catalog, err := module.GetCatalog()
		if err != nil {
			return nil, err
		}
		services = append(services, catalog.GetServices()...)
	}
	s.catalog = service.NewCatalog(services)
	catalogJSON, err := s.catalog.ToJSON()
	if err != nil {
		return nil, err
	}
	s.catalogResponse = catalogJSON

	s.listenAndServe = s.defaultListenAndServe

	return s, nil
}

func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	go func() {
		log.WithField(
			"address",
			fmt.Sprintf("http://0.0.0.0:%d", s.port),
		).Info("API server is listening")
		select {
		case errChan <- &errHTTPServerStopped{err: s.listenAndServe(ctx)}:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		log.Debug("context canceled; API server shutting down")
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
		select {
		case errChan <- svr.ListenAndServe():
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
		svr.Shutdown(shutdownCtx) // nolint: errcheck
		return ctx.Err()
	}
}
