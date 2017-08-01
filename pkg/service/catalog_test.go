package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCatalog     Catalog
	testCatalogJSON []byte
)

func init() {
	name := "test-name"
	id := "test-id"
	description := "test-description"
	tag := "test-tag"
	bindable := true
	planUpdatable := false
	free := false

	testCatalog = NewCatalog([]Service{
		NewService(
			&ServiceProperties{
				Name:          name,
				ID:            id,
				Description:   description,
				Tags:          []string{tag},
				Bindable:      bindable,
				PlanUpdatable: planUpdatable,
			},
			NewPlan(&PlanProperties{
				ID:          id,
				Name:        name,
				Description: description,
				Free:        free,
			}),
		),
	})

	testCatalogJSONStr := fmt.Sprintf(
		`{
			"services":[
				{
					"name":"%s",
					"id":"%s",
					"description":"%s",
					"tags":["%s"],
					"bindable":%t,
					"plan_updateable":%t,
					"plans":[
						{
							"id":"%s",
							"name":"%s",
							"description":"%s",
							"free":%t
						}
					]
				}
			]
		}`,
		name,
		id,
		description,
		tag,
		bindable,
		planUpdatable,
		id,
		name,
		description,
		free,
	)
	testCatalogJSONStr = strings.Replace(testCatalogJSONStr, " ", "", -1)
	testCatalogJSONStr = strings.Replace(testCatalogJSONStr, "\n", "", -1)
	testCatalogJSONStr = strings.Replace(testCatalogJSONStr, "\t", "", -1)
	testCatalogJSON = []byte(testCatalogJSONStr)
}

func TestNewCatalogFromJSON(t *testing.T) {
	catalog, err := NewCatalogFromJSON(testCatalogJSON)
	assert.Nil(t, err)
	assert.Equal(t, testCatalog, catalog)
}

func TestCatalogToJSON(t *testing.T) {
	json, err := testCatalog.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testCatalogJSON, json)
}

func TestGetNonExistingServiceByID(t *testing.T) {

}

func TestGetExistingServiceByID(t *testing.T) {

}

func TestGetNonExistingPlanByID(t *testing.T) {

}

func TestGetExistingPlanByID(t *testing.T) {

}
