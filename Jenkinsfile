def label = "worker-${UUID.randomUUID().toString()}"

def containers = []

withCredentials([
  [$class: 'StringBinding', credentialsId: 'DOCKER_HUB_USERNAME', variable: 'DOCKER_HUB_USERNAME'],
  [$class: 'StringBinding', credentialsId: 'DOCKER_HUB_PASSWORD', variable: 'DOCKER_HUB_PASSWORD']
]) {
  containers << containerTemplate(
    name: 'docker',
    image: 'docker:18.03.0-dind',
    ttyEnabled: true,
    command: 'cat',
    envVars: [
      envVar(key: 'DOCKER_HUB_USERNAME', value: env.DOCKER_HUB_USERNAME),
      envVar(key: 'DOCKER_HUB_PASSWORD', value: env.DOCKER_HUB_PASSWORD),
      // TODO: This should be set to microsoft/
      // But be careful with that-- my creds CAN push to there
      // envVar(key: 'DOCKER_REPO', value: 'microsoft/')
      envVar(key: 'DOCKER_REPO', value: 'krancour/')
    ]
  )
}

containers << containerTemplate(
  name: 'helm',
  image: 'quay.io/deis/helm-chart-publishing-tools:v0.1.0',
  ttyEnabled: true,
  command: 'cat',
  envVars: [
    envVar(key: 'SKIP_DOCKER', value: 'TRUE')
  ]
)

containers << containerTemplate(
  name: 'pcf-tile',
  image: 'cfplatformeng/tile-generator:v11.0.4',
  ttyEnabled: true,
  command: 'cat',
  envVars: [
    envVar(key: 'SKIP_DOCKER', value: 'TRUE')
  ]
)

def isMaster = env.BRANCH_NAME == 'master'
def isRelease = env.BRANCH_NAME ==~ /v[0-9]+(\.[0-9]+)*(\-.+)?/

if (!isRelease) {
  withCredentials([azureServicePrincipal('AZURE_CREDENTIALS')]) {
    containers << containerTemplate(
      name: 'go', image: 'quay.io/deis/lightweight-docker-go:v0.2.0',
      ttyEnabled: true,
      command: 'cat',
      envVars: [
        envVar(key: 'SKIP_DOCKER', value: 'TRUE'),
        envVar(key: 'STORAGE_REDIS_HOST', value: 'localhost'),
        envVar(key: 'ASYNC_REDIS_HOST', value: 'localhost'),
        envVar(key: 'CGO_ENABLED', value: '0'),
        envVar(key: 'AZURE_TENANT_ID', value: env.AZURE_TENANT_ID),
        envVar(key: 'AZURE_SUBSCRIPTION_ID', value: env.AZURE_SUBSCRIPTION_ID),
        envVar(key: 'AZURE_CLIENT_ID', value: env.AZURE_CLIENT_ID),
        envVar(key: 'AZURE_CLIENT_SECRET', value: env.AZURE_CLIENT_SECRET)
      ]
    )
  }
  containers << containerTemplate(
    name: 'redis',
    image: 'redis:3.2.4'
  )
  containers << containerTemplate(
    name: 'osb-checker', image: 'quay.io/deis/osb-checker:v0.3.0',
    ttyEnabled: true,
    command: 'cat',
    envVars: [
      envVar(key: 'SKIP_DOCKER', value: 'TRUE')
    ]
  )
}

podTemplate(
  label: label,
  containers: containers,
  volumes:[
    hostPathVolume(
      hostPath: '/var/run/docker.sock',
      mountPath: '/var/run/docker.sock'
    )
  ]
) {

  if (!isMaster && !isRelease) {
    timeout(time: 2, unit: 'DAYS') {
      stage('Hold for Approval') {
        input('Do you approve?')
      }
    }
  }

  node(label) {

    checkout scm

    stage('Prepare Pipeline') {
      if (!isRelease) {
        container('go') {
          sh """
          mkdir -p \$GOPATH/src/github.com/Azure
          ln -s \$(pwd) \$GOPATH/src/github.com/Azure/open-service-broker-azure
          """
        }
        container('osb-checker') {
          sh """
          mkdir -p \$GOPATH/src/github.com/Azure
          ln -s \$(pwd)/tests/api-compliance/localhost-config.json /app/config.json
          """
        }
      }
      container('pcf-tile') {
        sh """
        apk update
        apk add make git
        """
      }
      container('docker') {
        sh """
        apk update
        apk add make git
        docker login -u "\$DOCKER_HUB_USERNAME" -p "\$DOCKER_HUB_PASSWORD"
        """
      }
    }

    if (!isRelease) {

      stage('Preliminary Tests') {

        parallel (

          'lint': {
            stage('Lint') {
              container('go') {
                sh """
                cd \$GOPATH/src/github.com/Azure/open-service-broker-azure
                make lint
                """
              }
            }
          },

          'lint-chart': {
            stage('Lint Chart') {
              container('helm') {
                sh """
                make lint-chart
                """
              }
            }
          },

          'verify-vendored-code': {
            stage('Verify Vendored Code') {
              container('go') {
                sh """
                cd \$GOPATH/src/github.com/Azure/open-service-broker-azure
                make verify-vendored-code
                """
              }
            }
          },

          'test-unit': {
            stage('Run Unit Tests') {
              container('go') {
                sh """
                cd \$GOPATH/src/github.com/Azure/open-service-broker-azure
                make test-unit
                """
              }
            }
          },

          'test-api-compliance': {
            stage('Run OSB API Compliance Tests') {
              container('go') {
                sh """
                cd \$GOPATH/src/github.com/Azure/open-service-broker-azure
                go run cmd/compliance-test-broker/compliance-test-broker.go &
                """
              }
              container('osb-checker') {
                sh 'make test-api-compliance'
              }
            }
          },

          'build': {
            stage('Build') {
              container('docker') {
                sh 'make build'
              }
            }
          },

          'generate-pcf-tile': {
            stage('Generate PCF Tile') {
              container('pcf-tile') {
                sh 'make generate-pcf-tile'
              }
            }
          }

        )

      }

      stage('Service Lifecycle Tests') {
        container('go') {
          sh """
          cd \$GOPATH/src/github.com/Azure/open-service-broker-azure
          make test-service-lifecycles
          """
        }
      }

    }

    if (isMaster) {
      stage ('Publish RC Images') {
        container('docker') {
          sh 'make push-rc'
        }
      }
      stage('Publish RC Chart') {
        container('helm') {
          sh "make publish-rc-chart"
        }
      }
    }

    if (isRelease) {
      stage('Generate PCF Tile') {
        container('pcf-tile') {
          sh "REL_VERSION=$env.BRANCH_NAME make generate-pcf-tile"
        }
      }
      stage ('Publish Release Images') {
        container('docker') {
          sh "REL_VERSION=$env.BRANCH_NAME make push-release"
        }
      }
      stage('Publish Chart') {
        container('helm') {
          sh "REL_VERSION=$env.BRANCH_NAME make publish-release-chart"
        }
      }
    }

  }

}
