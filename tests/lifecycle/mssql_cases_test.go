// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	ss "github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/sqldb"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
)

func getMssqlCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]serviceLifecycleTestCase, error) {
	msSQLManager, err := ss.NewManager()
	if err != nil {
		return nil, err
	}

	return []serviceLifecycleTestCase{
		{ // new server scenario
			module:      sqldb.New(armDeployer, msSQLManager),
			description: "new server and database",
			serviceID:   "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:      "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &sqldb.ProvisioningParameters{
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
			},
			bindingParameters: &sqldb.BindingParameters{},
			testCredentials:   testMsSQLCreds(),
		},
	}, nil
}

func testMsSQLCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {
		cdts, ok := credentials.(*sqldb.Credentials)
		if !ok {
			return fmt.Errorf("error casting credentials as *mssql.Credentials")
		}

		query := url.Values{}
		query.Add("database", cdts.Database)
		query.Add("encrypt", "true")
		query.Add("TrustServerCertificate", "true")

		u := &url.URL{
			Scheme: "sqlserver",
			User: url.UserPassword(
				cdts.Username,
				cdts.Password,
			),
			Host:     fmt.Sprintf("%s:%d", cdts.Host, cdts.Port),
			RawQuery: query.Encode(),
		}

		db, err := sql.Open("mssql", u.String())
		if err != nil {
			return fmt.Errorf("error validating the database arguments: %s", err)
		}

		if err = db.Ping(); err != nil {
			return fmt.Errorf("error connecting to the database: %s", err)
		}
		defer db.Close() // nolint: errcheck

		rows, err := db.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='CONTROL'") // nolint: lll
		if err != nil {
			return fmt.Errorf(
				`error querying SELECT from table fn_my_permissions: %s`,
				err,
			)
		}
		defer rows.Close() // nolint: errcheck
		if !rows.Next() {
			return fmt.Errorf(
				`error user doesn't have permission 'CONTROL'`,
			)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf(
				`error iterating rows`,
			)
		}

		return nil
	}
}
