package filter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type failingFilter struct {
	called bool
}

func (f *failingFilter) GetHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.called = true
		w.WriteHeader(http.StatusBadRequest)
	}
}

type passingFilter struct {
	called bool
}

func (p *passingFilter) GetHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.called = true
		handler(w, r)
	}
}

func TestFilterChainWithNoFilters(t *testing.T) {
	var handlerCalled = false
	f1 := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}
	filterChain := NewChain()
	filtered := filterChain.GetHandler(f1)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	filtered(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, handlerCalled)
}

func TestFilterChainWithFailingFilter(t *testing.T) {
	failing := &failingFilter{}
	var handlerCalled = false
	f1 := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}
	filterChain := NewChain(failing)
	filtered := filterChain.GetHandler(f1)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	filtered(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.False(t, handlerCalled)
	assert.True(t, failing.called)
}

func TestFilterChainWithPassingFilter(t *testing.T) {
	passing := &passingFilter{}
	var handlerCalled = false
	f1 := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}
	filterChain := NewChain(passing)
	filtered := filterChain.GetHandler(f1)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	filtered(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, handlerCalled)
	assert.True(t, passing.called)
}

func TestFilterChainWithMixedFilters(t *testing.T) {
	passing := &passingFilter{}
	failing := &failingFilter{}
	var handlerCalled = false
	f1 := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}
	filterChain := NewChain(passing, failing)
	filtered := filterChain.GetHandler(f1)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	filtered(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.False(t, handlerCalled)
	assert.True(t, passing.called)
	assert.True(t, failing.called)
}
