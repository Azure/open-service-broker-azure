package storage

import (
	storageSDK "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	generalPurposeV1Manager *generalPurposeV1Manager
	generalPurposeV2Manager *generalPurposeV2Manager
	blobAccountManager      *blobAccountManager
	blobAllInOneManager     *blobAllInOneManager
}

type storageManager struct {
	armDeployer    arm.Deployer
	accountsClient storageSDK.AccountsClient
}

type generalPurposeV1Manager struct {
	storageManager
}

type generalPurposeV2Manager struct {
	storageManager
}

type blobAccountManager struct {
	storageManager
}

type blobAllInOneManager struct {
	storageManager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Storage using "Azure Storage"
func New(
	armDeployer arm.Deployer,
	accountsClient storageSDK.AccountsClient,
) service.Module {
	storageMgr := storageManager{
		armDeployer:    armDeployer,
		accountsClient: accountsClient,
	}
	return &module{
		generalPurposeV1Manager: &generalPurposeV1Manager{storageMgr},
		generalPurposeV2Manager: &generalPurposeV2Manager{storageMgr},
		blobAccountManager:      &blobAccountManager{storageMgr},
		blobAllInOneManager:     &blobAllInOneManager{storageMgr},
	}
}

func (m *module) GetName() string {
	return "storage"
}
