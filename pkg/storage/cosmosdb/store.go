package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

// Store is a (./pkg/storage).Store interface implementation to store
// data in Microsoft Azure CosmosDB.
//
// See https://docs.microsoft.com/en-us/azure/cosmos-db/introduction)
// for more details on CosmosDB
type Store struct {
	db           *mgo.Database
	instCollName string
	bindCollName string
}

// NewStore creates a new Store implementation that is backed by CosmosDB's
// MongoDB API
//
// The given session should be configured to talk to CosmosDB, and the given
// instDBName and bindDBName will be used as the names for the Mongo
// collections in which to store instances and bindings, respectively
func NewStore(db *mgo.Database, instCollName, bindCollName string) *Store {
	return &Store{
		db:           db,
		instCollName: instCollName,
		bindCollName: bindCollName,
	}
}

// WriteInstance persists the given instance to the underlying storage
func (s *Store) WriteInstance(instance service.Instance) error {
	coll := s.instColl()
	return coll.Insert(instance)
}

// GetInstance retrieves a persisted instance from the underlying storage by
// instance id
func (s *Store) GetInstance(instanceID string) (service.Instance, bool, error) {
	res := new(service.Instance)
	coll := s.instColl()
	findErr := coll.Find(bson.M{"instanceId": instanceID}).One(res)
	if findErr != nil {
		return service.Instance{}, false, findErr
	}
	return *res, true, nil
}

// DeleteInstance deletes a persisted instance from the underlying storage by
// instance id
func (s *Store) DeleteInstance(instanceID string) (bool, error) {
	coll := s.instColl()
	if err := coll.Remove(bson.M{"instanceId": instanceID}); err != nil {
		return false, err
	}
	return true, nil
}

// WriteBinding persists the given binding to the underlying storage
func (s *Store) WriteBinding(binding service.Binding) error {
	coll := s.bindColl()
	return coll.Insert(binding)
}

// GetBinding retrieves a persisted instance from the underlying storage by
// binding id
func (s *Store) GetBinding(bindingID string) (service.Binding, bool, error) {
	res := new(service.Binding)
	coll := s.bindColl()
	findErr := coll.Find(bson.M{"bindingId": bindingID}).One(res)
	if findErr != nil {
		return service.Binding{}, false, findErr
	}
	return *res, true, nil
}

// DeleteBinding deletes a persisted binding from the underlying storage by
// binding id
func (s *Store) DeleteBinding(bindingID string) (bool, error) {
	coll := s.bindColl()
	if err := coll.Remove(bson.M{"bindingId": bindingID}); err != nil {
		return false, err
	}
	return true, nil
}

// TestConnection tests the connection to the underlying database (if there
// is one)
func (s *Store) TestConnection() error {
	return s.db.Session.Ping()
}

func (s *Store) instColl() *mgo.Collection {
	return s.db.C(s.instCollName)
}

func (s *Store) bindColl() *mgo.Collection {
	return s.db.C(s.bindCollName)
}
