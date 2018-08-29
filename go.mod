module github.com/Azure/open-service-broker-azure

require (
	cloud.google.com/go v0.26.0 // indirect
	github.com/Azure/azure-sdk-for-go v19.1.0+incompatible
	github.com/Azure/go-autorest v10.15.2+incompatible
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.15.0+incompatible
	github.com/Sirupsen/logrus v1.0.6
	github.com/aokoli/goutils v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20180824013952-8fac8b954edb
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-redis/redis v6.14.0+incompatible
	github.com/go-sql-driver/mysql v1.4.0
	github.com/google/uuid v1.0.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/huandu/xstrings v1.1.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/krancour/async v1.0.0
	github.com/lib/pq v1.0.0
	github.com/marstr/guid v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.0.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.2.2
	github.com/urfave/cli v1.20.0
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac // indirect
	golang.org/x/sys v0.0.0-20180828065106-d99a578cf41b // indirect
	google.golang.org/appengine v1.1.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

replace github.com/Sirupsen/logrus v1.0.6 => github.com/krancour/logrus v1.0.4-0.20171115205400-9da25b464c10
