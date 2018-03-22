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
	props := map[string]*ParameterSchema{}
	props["location"] = &ParameterSchema{
		Type: "string",
		Description: "The Azure region in which to provision " +
			"applicable resources.",
	}
	props["resourceGroup"] = &ParameterSchema{
		Type: "string",
		Description: "The (new or existing) resource group with " +
			"which to associate new resources.",
	}
	props["tags"] = &ParameterSchema{
		Type: "object",
		Description: "Tags to be applied to new resources, specified " +
			"as key/value pairs.",
		Additional: &ParameterSchema{
			Type: "string",
		},
	}

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
			nil,
			NewPlan(&PlanProperties{
				ID:          id,
				Name:        name,
				Description: description,
				Free:        free,
				ProvisionParamsSchema: props,
			}),
		),
	})

	testCatalogTemplate := `{
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
						"free":%t,
						"schemas": {
							"service_instance": {
								"create": {
									"parameters": {
										"$schema": "http://json-schema.org/draft-04/schema#",
										"type": "object",
										"properties": {
											"location": {
												"type": "string",
												"description": "%s"
											},
											"resourceGroup": {
												"type": "string",
												"description": "%s"
											},
											"tags": {
												"type": "object",
												"description": "%s",
												"additional": {
													"type" : "string"
												}
											}
										}
									}
								}
							}
						}
					}
				]
			}
		]
	}`

	testCatalogTemplate = strings.Replace(testCatalogTemplate, " ", "", -1)
	testCatalogTemplate = strings.Replace(testCatalogTemplate, "\n", "", -1)
	testCatalogTemplate = strings.Replace(testCatalogTemplate, "\t", "", -1)

	testCatalogJSONStr := fmt.Sprintf(
		testCatalogTemplate,
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
		"The Azure region in which to provision applicable resources.",
		"The (new or existing) resource group with which to associate new resources.",
		"Tags to be applied to new resources, specified as key/value pairs.",
	)

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
