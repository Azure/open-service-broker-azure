package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/file"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/krancour/async"
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
	// Run causes the api server to start serving HTTP requests. It will block
	// until an error occurs and will return that error.
	Run(context.Context) error
}

type server struct {
	apiServerConfig Config
	store           storage.Store
	asyncEngine     async.Engine
	filterChain     filter.Filter
	router          *mux.Router
	catalog         service.Catalog
	catalogResponse []byte
	// This allows tests to inject an alternative implementation of this function
	listenAndServe func(context.Context) error
}

// NewServer returns an HTTP router
func NewServer(
	apiServerConfig Config,
	store storage.Store,
	asyncEngine async.Engine,
	filterChain filter.Filter,
	catalog service.Catalog,
) (Server, error) {
	s := &server{
		apiServerConfig: apiServerConfig,
		store:           store,
		asyncEngine:     asyncEngine,
		filterChain:     filterChain,
		catalog:         catalog,
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

	catalogJSON, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}
	s.catalogResponse = catalogJSON

	s.listenAndServe = s.defaultListenAndServe

	return s, nil
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	go func() {
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
		Addr:    fmt.Sprintf(":%d", s.apiServerConfig.Port),
		Handler: s.router,
	}
	if s.apiServerConfig.TLSCertPath != "" &&
		s.apiServerConfig.TLSKeyPath != "" &&
		file.Exists(s.apiServerConfig.TLSCertPath) &&
		file.Exists(s.apiServerConfig.TLSKeyPath) {
		log.WithField(
			"address",
			fmt.Sprintf("https://0.0.0.0:%d", s.apiServerConfig.Port),
		).Info("API server is listening with TLS enabled")
		go func() {
			select {
			case errChan <- svr.ListenAndServeTLS(
				s.apiServerConfig.TLSCertPath,
				s.apiServerConfig.TLSKeyPath,
			):
			case <-ctx.Done():
			}
		}()
	} else {
		log.WithField(
			"address",
			fmt.Sprintf("http://0.0.0.0:%d", s.apiServerConfig.Port),
		).Warn("API server is listening with TLS disabled")
		go func() {
			select {
			case errChan <- svr.ListenAndServe():
			case <-ctx.Done():
			}
		}()
	}
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
