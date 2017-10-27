// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	ss "github.com/Azure/azure-service-broker/pkg/azure/mssql"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/mssql"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
)

func getMssqlCases(
	armDeployer arm.Deployer,
) ([]moduleLifecycleTestCase, error) {
	msSQLManager, err := ss.NewManager()
	if err != nil {
		return nil, err
	}
	msSQLConfig, err := mssql.GetConfig()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    mssql.New(armDeployer, msSQLManager, msSQLConfig),
			serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:    "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			provisioningParameters: &mssql.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &mssql.BindingParameters{},
			testCredentials:   testMsSQLCreds(),
		},
	}, nil
}

func testMsSQLCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {
		cdts, ok := credentials.(*mssql.Credentials)
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
