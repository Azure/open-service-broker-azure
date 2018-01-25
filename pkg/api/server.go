package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
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
	port            int
	store           storage.Store
	asyncEngine     async.Engine
	filterChain     filter.Filter
	router          *mux.Router
	catalog         service.Catalog
	catalogResponse []byte
	// This allows tests to inject an alternative implementation of this function
	listenAndServe            func(context.Context) error
	defaultAzureLocation      string
	defaultAzureResourceGroup string
}

// NewServer returns an HTTP router
func NewServer(
	port int,
	store storage.Store,
	asyncEngine async.Engine,
	filterChain filter.Filter,
	catalog service.Catalog,
	defaultAzureLocation string,
	defaultAzureResourceGroup string,
) (Server, error) {
	s := &server{
		port:                      port,
		store:                     store,
		asyncEngine:               asyncEngine,
		filterChain:               filterChain,
		catalog:                   catalog,
		defaultAzureLocation:      defaultAzureLocation,
		defaultAzureResourceGroup: defaultAzureResourceGroup,
	}

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc(
		"/v2/catalog",
		filterChain.GetHandler(s.getCatalog),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		filterChain.GetHandler(s.provision),
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		filterChain.GetHandler(s.update),
	).Methods(http.MethodPatch)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/last_operation",
		filterChain.GetHandler(s.poll),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		filterChain.GetHandler(s.bind),
	).Methods(http.MethodPut)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		filterChain.GetHandler(s.unbind),
	).Methods(http.MethodDelete)
	router.HandleFunc(
		"/v2/service_instances/{instance_id}",
		filterChain.GetHandler(s.deprovision),
	).Methods(http.MethodDelete)
	router.HandleFunc(
		"/healthz",
		s.healthCheck, // Filter chain not applied to this reqeust
	).Methods(http.MethodGet)
	s.router = router

	catalogJSON, err := catalog.ToJSON()
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
