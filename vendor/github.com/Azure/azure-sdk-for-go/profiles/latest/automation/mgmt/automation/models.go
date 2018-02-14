// +build go1.9

// Copyright 2017 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder
// commit ID: 2014fbbf031942474ad27a5a66dffaed5347f3fb

package automation

import original "github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"

type CredentialClient = original.CredentialClient
type JobStreamClient = original.JobStreamClient
type NodeReportsClient = original.NodeReportsClient
type TestJobStreamsClient = original.TestJobStreamsClient
type AccountClient = original.AccountClient
type JobClient = original.JobClient
type RunbookDraftClient = original.RunbookDraftClient
type WebhookClient = original.WebhookClient
type ActivityClient = original.ActivityClient
type FieldsClient = original.FieldsClient
type JobScheduleClient = original.JobScheduleClient
type AccountState = original.AccountState

const (
	Ok          AccountState = original.Ok
	Suspended   AccountState = original.Suspended
	Unavailable AccountState = original.Unavailable
)

type AgentRegistrationKeyName = original.AgentRegistrationKeyName

const (
	Primary   AgentRegistrationKeyName = original.Primary
	Secondary AgentRegistrationKeyName = original.Secondary
)

type ContentSourceType = original.ContentSourceType

const (
	EmbeddedContent ContentSourceType = original.EmbeddedContent
	URI             ContentSourceType = original.URI
)

type DscConfigurationProvisioningState = original.DscConfigurationProvisioningState

const (
	Succeeded DscConfigurationProvisioningState = original.Succeeded
)

type DscConfigurationState = original.DscConfigurationState

const (
	DscConfigurationStateEdit      DscConfigurationState = original.DscConfigurationStateEdit
	DscConfigurationStateNew       DscConfigurationState = original.DscConfigurationStateNew
	DscConfigurationStatePublished DscConfigurationState = original.DscConfigurationStatePublished
)

type HTTPStatusCode = original.HTTPStatusCode

const (
	Accepted                     HTTPStatusCode = original.Accepted
	Ambiguous                    HTTPStatusCode = original.Ambiguous
	BadGateway                   HTTPStatusCode = original.BadGateway
	BadRequest                   HTTPStatusCode = original.BadRequest
	Conflict                     HTTPStatusCode = original.Conflict
	Continue                     HTTPStatusCode = original.Continue
	Created                      HTTPStatusCode = original.Created
	ExpectationFailed            HTTPStatusCode = original.ExpectationFailed
	Forbidden                    HTTPStatusCode = original.Forbidden
	Found                        HTTPStatusCode = original.Found
	GatewayTimeout               HTTPStatusCode = original.GatewayTimeout
	Gone                         HTTPStatusCode = original.Gone
	HTTPVersionNotSupported      HTTPStatusCode = original.HTTPVersionNotSupported
	InternalServerError          HTTPStatusCode = original.InternalServerError
	LengthRequired               HTTPStatusCode = original.LengthRequired
	MethodNotAllowed             HTTPStatusCode = original.MethodNotAllowed
	Moved                        HTTPStatusCode = original.Moved
	MovedPermanently             HTTPStatusCode = original.MovedPermanently
	MultipleChoices              HTTPStatusCode = original.MultipleChoices
	NoContent                    HTTPStatusCode = original.NoContent
	NonAuthoritativeInformation  HTTPStatusCode = original.NonAuthoritativeInformation
	NotAcceptable                HTTPStatusCode = original.NotAcceptable
	NotFound                     HTTPStatusCode = original.NotFound
	NotImplemented               HTTPStatusCode = original.NotImplemented
	NotModified                  HTTPStatusCode = original.NotModified
	OK                           HTTPStatusCode = original.OK
	PartialContent               HTTPStatusCode = original.PartialContent
	PaymentRequired              HTTPStatusCode = original.PaymentRequired
	PreconditionFailed           HTTPStatusCode = original.PreconditionFailed
	ProxyAuthenticationRequired  HTTPStatusCode = original.ProxyAuthenticationRequired
	Redirect                     HTTPStatusCode = original.Redirect
	RedirectKeepVerb             HTTPStatusCode = original.RedirectKeepVerb
	RedirectMethod               HTTPStatusCode = original.RedirectMethod
	RequestedRangeNotSatisfiable HTTPStatusCode = original.RequestedRangeNotSatisfiable
	RequestEntityTooLarge        HTTPStatusCode = original.RequestEntityTooLarge
	RequestTimeout               HTTPStatusCode = original.RequestTimeout
	RequestURITooLong            HTTPStatusCode = original.RequestURITooLong
	ResetContent                 HTTPStatusCode = original.ResetContent
	SeeOther                     HTTPStatusCode = original.SeeOther
	ServiceUnavailable           HTTPStatusCode = original.ServiceUnavailable
	SwitchingProtocols           HTTPStatusCode = original.SwitchingProtocols
	TemporaryRedirect            HTTPStatusCode = original.TemporaryRedirect
	Unauthorized                 HTTPStatusCode = original.Unauthorized
	UnsupportedMediaType         HTTPStatusCode = original.UnsupportedMediaType
	Unused                       HTTPStatusCode = original.Unused
	UpgradeRequired              HTTPStatusCode = original.UpgradeRequired
	UseProxy                     HTTPStatusCode = original.UseProxy
)

type JobStatus = original.JobStatus

const (
	JobStatusActivating   JobStatus = original.JobStatusActivating
	JobStatusBlocked      JobStatus = original.JobStatusBlocked
	JobStatusCompleted    JobStatus = original.JobStatusCompleted
	JobStatusDisconnected JobStatus = original.JobStatusDisconnected
	JobStatusFailed       JobStatus = original.JobStatusFailed
	JobStatusNew          JobStatus = original.JobStatusNew
	JobStatusRemoving     JobStatus = original.JobStatusRemoving
	JobStatusResuming     JobStatus = original.JobStatusResuming
	JobStatusRunning      JobStatus = original.JobStatusRunning
	JobStatusStopped      JobStatus = original.JobStatusStopped
	JobStatusStopping     JobStatus = original.JobStatusStopping
	JobStatusSuspended    JobStatus = original.JobStatusSuspended
	JobStatusSuspending   JobStatus = original.JobStatusSuspending
)

type JobStreamType = original.JobStreamType

const (
	Any      JobStreamType = original.Any
	Debug    JobStreamType = original.Debug
	Error    JobStreamType = original.Error
	Output   JobStreamType = original.Output
	Progress JobStreamType = original.Progress
	Verbose  JobStreamType = original.Verbose
	Warning  JobStreamType = original.Warning
)

type ModuleProvisioningState = original.ModuleProvisioningState

const (
	ModuleProvisioningStateActivitiesStored            ModuleProvisioningState = original.ModuleProvisioningStateActivitiesStored
	ModuleProvisioningStateCancelled                   ModuleProvisioningState = original.ModuleProvisioningStateCancelled
	ModuleProvisioningStateConnectionTypeImported      ModuleProvisioningState = original.ModuleProvisioningStateConnectionTypeImported
	ModuleProvisioningStateContentDownloaded           ModuleProvisioningState = original.ModuleProvisioningStateContentDownloaded
	ModuleProvisioningStateContentRetrieved            ModuleProvisioningState = original.ModuleProvisioningStateContentRetrieved
	ModuleProvisioningStateContentStored               ModuleProvisioningState = original.ModuleProvisioningStateContentStored
	ModuleProvisioningStateContentValidated            ModuleProvisioningState = original.ModuleProvisioningStateContentValidated
	ModuleProvisioningStateCreated                     ModuleProvisioningState = original.ModuleProvisioningStateCreated
	ModuleProvisioningStateCreating                    ModuleProvisioningState = original.ModuleProvisioningStateCreating
	ModuleProvisioningStateFailed                      ModuleProvisioningState = original.ModuleProvisioningStateFailed
	ModuleProvisioningStateModuleDataStored            ModuleProvisioningState = original.ModuleProvisioningStateModuleDataStored
	ModuleProvisioningStateModuleImportRunbookComplete ModuleProvisioningState = original.ModuleProvisioningStateModuleImportRunbookComplete
	ModuleProvisioningStateRunningImportModuleRunbook  ModuleProvisioningState = original.ModuleProvisioningStateRunningImportModuleRunbook
	ModuleProvisioningStateStartingImportModuleRunbook ModuleProvisioningState = original.ModuleProvisioningStateStartingImportModuleRunbook
	ModuleProvisioningStateSucceeded                   ModuleProvisioningState = original.ModuleProvisioningStateSucceeded
	ModuleProvisioningStateUpdating                    ModuleProvisioningState = original.ModuleProvisioningStateUpdating
)

type RunbookProvisioningState = original.RunbookProvisioningState

const (
	RunbookProvisioningStateSucceeded RunbookProvisioningState = original.RunbookProvisioningStateSucceeded
)

type RunbookState = original.RunbookState

const (
	RunbookStateEdit      RunbookState = original.RunbookStateEdit
	RunbookStateNew       RunbookState = original.RunbookStateNew
	RunbookStatePublished RunbookState = original.RunbookStatePublished
)

type RunbookTypeEnum = original.RunbookTypeEnum

const (
	Graph                   RunbookTypeEnum = original.Graph
	GraphPowerShell         RunbookTypeEnum = original.GraphPowerShell
	GraphPowerShellWorkflow RunbookTypeEnum = original.GraphPowerShellWorkflow
	PowerShell              RunbookTypeEnum = original.PowerShell
	PowerShellWorkflow      RunbookTypeEnum = original.PowerShellWorkflow
	Script                  RunbookTypeEnum = original.Script
)

type ScheduleDay = original.ScheduleDay

const (
	Friday    ScheduleDay = original.Friday
	Monday    ScheduleDay = original.Monday
	Saturday  ScheduleDay = original.Saturday
	Sunday    ScheduleDay = original.Sunday
	Thursday  ScheduleDay = original.Thursday
	Tuesday   ScheduleDay = original.Tuesday
	Wednesday ScheduleDay = original.Wednesday
)

type ScheduleFrequency = original.ScheduleFrequency

const (
	Day     ScheduleFrequency = original.Day
	Hour    ScheduleFrequency = original.Hour
	Month   ScheduleFrequency = original.Month
	OneTime ScheduleFrequency = original.OneTime
	Week    ScheduleFrequency = original.Week
)

type SkuNameEnum = original.SkuNameEnum

const (
	Basic SkuNameEnum = original.Basic
	Free  SkuNameEnum = original.Free
)

type Account = original.Account
type AccountCreateOrUpdateParameters = original.AccountCreateOrUpdateParameters
type AccountCreateOrUpdateProperties = original.AccountCreateOrUpdateProperties
type AccountListResult = original.AccountListResult
type AccountListResultIterator = original.AccountListResultIterator
type AccountListResultPage = original.AccountListResultPage
type AccountProperties = original.AccountProperties
type AccountUpdateParameters = original.AccountUpdateParameters
type AccountUpdateProperties = original.AccountUpdateProperties
type Activity = original.Activity
type ActivityListResult = original.ActivityListResult
type ActivityListResultIterator = original.ActivityListResultIterator
type ActivityListResultPage = original.ActivityListResultPage
type ActivityOutputType = original.ActivityOutputType
type ActivityParameter = original.ActivityParameter
type ActivityParameterSet = original.ActivityParameterSet
type ActivityProperties = original.ActivityProperties
type AdvancedSchedule = original.AdvancedSchedule
type AdvancedScheduleMonthlyOccurrence = original.AdvancedScheduleMonthlyOccurrence
type AgentRegistration = original.AgentRegistration
type AgentRegistrationKeys = original.AgentRegistrationKeys
type AgentRegistrationRegenerateKeyParameter = original.AgentRegistrationRegenerateKeyParameter
type Certificate = original.Certificate
type CertificateCreateOrUpdateParameters = original.CertificateCreateOrUpdateParameters
type CertificateCreateOrUpdateProperties = original.CertificateCreateOrUpdateProperties
type CertificateListResult = original.CertificateListResult
type CertificateListResultIterator = original.CertificateListResultIterator
type CertificateListResultPage = original.CertificateListResultPage
type CertificateProperties = original.CertificateProperties
type CertificateUpdateParameters = original.CertificateUpdateParameters
type CertificateUpdateProperties = original.CertificateUpdateProperties
type Connection = original.Connection
type ConnectionCreateOrUpdateParameters = original.ConnectionCreateOrUpdateParameters
type ConnectionCreateOrUpdateProperties = original.ConnectionCreateOrUpdateProperties
type ConnectionListResult = original.ConnectionListResult
type ConnectionListResultIterator = original.ConnectionListResultIterator
type ConnectionListResultPage = original.ConnectionListResultPage
type ConnectionProperties = original.ConnectionProperties
type ConnectionType = original.ConnectionType
type ConnectionTypeAssociationProperty = original.ConnectionTypeAssociationProperty
type ConnectionTypeCreateOrUpdateParameters = original.ConnectionTypeCreateOrUpdateParameters
type ConnectionTypeCreateOrUpdateProperties = original.ConnectionTypeCreateOrUpdateProperties
type ConnectionTypeListResult = original.ConnectionTypeListResult
type ConnectionTypeListResultIterator = original.ConnectionTypeListResultIterator
type ConnectionTypeListResultPage = original.ConnectionTypeListResultPage
type ConnectionTypeProperties = original.ConnectionTypeProperties
type ConnectionUpdateParameters = original.ConnectionUpdateParameters
type ConnectionUpdateProperties = original.ConnectionUpdateProperties
type ContentHash = original.ContentHash
type ContentLink = original.ContentLink
type ContentSource = original.ContentSource
type Credential = original.Credential
type CredentialCreateOrUpdateParameters = original.CredentialCreateOrUpdateParameters
type CredentialCreateOrUpdateProperties = original.CredentialCreateOrUpdateProperties
type CredentialListResult = original.CredentialListResult
type CredentialListResultIterator = original.CredentialListResultIterator
type CredentialListResultPage = original.CredentialListResultPage
type CredentialProperties = original.CredentialProperties
type CredentialUpdateParameters = original.CredentialUpdateParameters
type CredentialUpdateProperties = original.CredentialUpdateProperties
type DscCompilationJob = original.DscCompilationJob
type DscCompilationJobCreateParameters = original.DscCompilationJobCreateParameters
type DscCompilationJobCreateProperties = original.DscCompilationJobCreateProperties
type DscCompilationJobListResult = original.DscCompilationJobListResult
type DscCompilationJobListResultIterator = original.DscCompilationJobListResultIterator
type DscCompilationJobListResultPage = original.DscCompilationJobListResultPage
type DscCompilationJobProperties = original.DscCompilationJobProperties
type DscConfiguration = original.DscConfiguration
type DscConfigurationAssociationProperty = original.DscConfigurationAssociationProperty
type DscConfigurationCreateOrUpdateParameters = original.DscConfigurationCreateOrUpdateParameters
type DscConfigurationCreateOrUpdateProperties = original.DscConfigurationCreateOrUpdateProperties
type DscConfigurationListResult = original.DscConfigurationListResult
type DscConfigurationListResultIterator = original.DscConfigurationListResultIterator
type DscConfigurationListResultPage = original.DscConfigurationListResultPage
type DscConfigurationParameter = original.DscConfigurationParameter
type DscConfigurationProperties = original.DscConfigurationProperties
type DscMetaConfiguration = original.DscMetaConfiguration
type DscNode = original.DscNode
type DscNodeConfiguration = original.DscNodeConfiguration
type DscNodeConfigurationAssociationProperty = original.DscNodeConfigurationAssociationProperty
type DscNodeConfigurationCreateOrUpdateParameters = original.DscNodeConfigurationCreateOrUpdateParameters
type DscNodeConfigurationListResult = original.DscNodeConfigurationListResult
type DscNodeConfigurationListResultIterator = original.DscNodeConfigurationListResultIterator
type DscNodeConfigurationListResultPage = original.DscNodeConfigurationListResultPage
type DscNodeExtensionHandlerAssociationProperty = original.DscNodeExtensionHandlerAssociationProperty
type DscNodeListResult = original.DscNodeListResult
type DscNodeListResultIterator = original.DscNodeListResultIterator
type DscNodeListResultPage = original.DscNodeListResultPage
type DscNodeReport = original.DscNodeReport
type DscNodeReportListResult = original.DscNodeReportListResult
type DscNodeReportListResultIterator = original.DscNodeReportListResultIterator
type DscNodeReportListResultPage = original.DscNodeReportListResultPage
type DscNodeUpdateParameters = original.DscNodeUpdateParameters
type DscReportError = original.DscReportError
type DscReportResource = original.DscReportResource
type DscReportResourceNavigation = original.DscReportResourceNavigation
type ErrorResponse = original.ErrorResponse
type FieldDefinition = original.FieldDefinition
type HybridRunbookWorker = original.HybridRunbookWorker
type HybridRunbookWorkerGroup = original.HybridRunbookWorkerGroup
type HybridRunbookWorkerGroupsListResult = original.HybridRunbookWorkerGroupsListResult
type HybridRunbookWorkerGroupsListResultIterator = original.HybridRunbookWorkerGroupsListResultIterator
type HybridRunbookWorkerGroupsListResultPage = original.HybridRunbookWorkerGroupsListResultPage
type HybridRunbookWorkerGroupUpdateParameters = original.HybridRunbookWorkerGroupUpdateParameters
type Job = original.Job
type JobCreateParameters = original.JobCreateParameters
type JobCreateProperties = original.JobCreateProperties
type JobListResult = original.JobListResult
type JobListResultIterator = original.JobListResultIterator
type JobListResultPage = original.JobListResultPage
type JobProperties = original.JobProperties
type JobSchedule = original.JobSchedule
type JobScheduleCreateParameters = original.JobScheduleCreateParameters
type JobScheduleCreateProperties = original.JobScheduleCreateProperties
type JobScheduleListResult = original.JobScheduleListResult
type JobScheduleListResultIterator = original.JobScheduleListResultIterator
type JobScheduleListResultPage = original.JobScheduleListResultPage
type JobScheduleProperties = original.JobScheduleProperties
type JobStream = original.JobStream
type JobStreamListResult = original.JobStreamListResult
type JobStreamListResultIterator = original.JobStreamListResultIterator
type JobStreamListResultPage = original.JobStreamListResultPage
type JobStreamProperties = original.JobStreamProperties
type Module = original.Module
type ModuleCreateOrUpdateParameters = original.ModuleCreateOrUpdateParameters
type ModuleCreateOrUpdateProperties = original.ModuleCreateOrUpdateProperties
type ModuleErrorInfo = original.ModuleErrorInfo
type ModuleListResult = original.ModuleListResult
type ModuleListResultIterator = original.ModuleListResultIterator
type ModuleListResultPage = original.ModuleListResultPage
type ModuleProperties = original.ModuleProperties
type ModuleUpdateParameters = original.ModuleUpdateParameters
type ModuleUpdateProperties = original.ModuleUpdateProperties
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type ReadCloser = original.ReadCloser
type Resource = original.Resource
type RunAsCredentialAssociationProperty = original.RunAsCredentialAssociationProperty
type Runbook = original.Runbook
type RunbookAssociationProperty = original.RunbookAssociationProperty
type RunbookCreateOrUpdateDraftParameters = original.RunbookCreateOrUpdateDraftParameters
type RunbookCreateOrUpdateDraftProperties = original.RunbookCreateOrUpdateDraftProperties
type RunbookCreateOrUpdateParameters = original.RunbookCreateOrUpdateParameters
type RunbookCreateOrUpdateProperties = original.RunbookCreateOrUpdateProperties
type RunbookDraft = original.RunbookDraft
type RunbookDraftCreateOrUpdateFuture = original.RunbookDraftCreateOrUpdateFuture
type RunbookDraftPublishFuture = original.RunbookDraftPublishFuture
type RunbookDraftUndoEditResult = original.RunbookDraftUndoEditResult
type RunbookListResult = original.RunbookListResult
type RunbookListResultIterator = original.RunbookListResultIterator
type RunbookListResultPage = original.RunbookListResultPage
type RunbookParameter = original.RunbookParameter
type RunbookProperties = original.RunbookProperties
type RunbookUpdateParameters = original.RunbookUpdateParameters
type RunbookUpdateProperties = original.RunbookUpdateProperties
type Schedule = original.Schedule
type ScheduleAssociationProperty = original.ScheduleAssociationProperty
type ScheduleCreateOrUpdateParameters = original.ScheduleCreateOrUpdateParameters
type ScheduleCreateOrUpdateProperties = original.ScheduleCreateOrUpdateProperties
type ScheduleListResult = original.ScheduleListResult
type ScheduleListResultIterator = original.ScheduleListResultIterator
type ScheduleListResultPage = original.ScheduleListResultPage
type ScheduleProperties = original.ScheduleProperties
type ScheduleUpdateParameters = original.ScheduleUpdateParameters
type ScheduleUpdateProperties = original.ScheduleUpdateProperties
type Sku = original.Sku
type Statistics = original.Statistics
type StatisticsListResult = original.StatisticsListResult
type String = original.String
type SubResource = original.SubResource
type TestJob = original.TestJob
type TestJobCreateParameters = original.TestJobCreateParameters
type TypeField = original.TypeField
type TypeFieldListResult = original.TypeFieldListResult
type Usage = original.Usage
type UsageCounterName = original.UsageCounterName
type UsageListResult = original.UsageListResult
type Variable = original.Variable
type VariableCreateOrUpdateParameters = original.VariableCreateOrUpdateParameters
type VariableCreateOrUpdateProperties = original.VariableCreateOrUpdateProperties
type VariableListResult = original.VariableListResult
type VariableListResultIterator = original.VariableListResultIterator
type VariableListResultPage = original.VariableListResultPage
type VariableProperties = original.VariableProperties
type VariableUpdateParameters = original.VariableUpdateParameters
type VariableUpdateProperties = original.VariableUpdateProperties
type Webhook = original.Webhook
type WebhookCreateOrUpdateParameters = original.WebhookCreateOrUpdateParameters
type WebhookCreateOrUpdateProperties = original.WebhookCreateOrUpdateProperties
type WebhookListResult = original.WebhookListResult
type WebhookListResultIterator = original.WebhookListResultIterator
type WebhookListResultPage = original.WebhookListResultPage
type WebhookProperties = original.WebhookProperties
type WebhookUpdateParameters = original.WebhookUpdateParameters
type WebhookUpdateProperties = original.WebhookUpdateProperties
type TestJobsClient = original.TestJobsClient
type CertificateClient = original.CertificateClient
type DscCompilationJobClient = original.DscCompilationJobClient
type DscNodeClient = original.DscNodeClient
type HybridRunbookWorkerGroupClient = original.HybridRunbookWorkerGroupClient
type UsagesClient = original.UsagesClient
type AgentRegistrationInformationClient = original.AgentRegistrationInformationClient
type ConnectionClient = original.ConnectionClient
type DscConfigurationClient = original.DscConfigurationClient
type ModuleClient = original.ModuleClient
type ScheduleClient = original.ScheduleClient
type DscNodeConfigurationClient = original.DscNodeConfigurationClient
type ObjectDataTypesClient = original.ObjectDataTypesClient
type RunbookClient = original.RunbookClient

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type OperationsClient = original.OperationsClient
type VariableClient = original.VariableClient
type ConnectionTypeClient = original.ConnectionTypeClient
type StatisticsClient = original.StatisticsClient

func NewTestJobsClient(subscriptionID string, resourceGroupName string) TestJobsClient {
	return original.NewTestJobsClient(subscriptionID, resourceGroupName)
}
func NewTestJobsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) TestJobsClient {
	return original.NewTestJobsClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewActivityClient(subscriptionID string, resourceGroupName string) ActivityClient {
	return original.NewActivityClient(subscriptionID, resourceGroupName)
}
func NewActivityClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ActivityClient {
	return original.NewActivityClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewFieldsClient(subscriptionID string, resourceGroupName string) FieldsClient {
	return original.NewFieldsClient(subscriptionID, resourceGroupName)
}
func NewFieldsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) FieldsClient {
	return original.NewFieldsClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewJobScheduleClient(subscriptionID string, resourceGroupName string) JobScheduleClient {
	return original.NewJobScheduleClient(subscriptionID, resourceGroupName)
}
func NewJobScheduleClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) JobScheduleClient {
	return original.NewJobScheduleClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewHybridRunbookWorkerGroupClient(subscriptionID string, resourceGroupName string) HybridRunbookWorkerGroupClient {
	return original.NewHybridRunbookWorkerGroupClient(subscriptionID, resourceGroupName)
}
func NewHybridRunbookWorkerGroupClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) HybridRunbookWorkerGroupClient {
	return original.NewHybridRunbookWorkerGroupClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewUsagesClient(subscriptionID string, resourceGroupName string) UsagesClient {
	return original.NewUsagesClient(subscriptionID, resourceGroupName)
}
func NewUsagesClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) UsagesClient {
	return original.NewUsagesClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewCertificateClient(subscriptionID string, resourceGroupName string) CertificateClient {
	return original.NewCertificateClient(subscriptionID, resourceGroupName)
}
func NewCertificateClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) CertificateClient {
	return original.NewCertificateClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewDscCompilationJobClient(subscriptionID string, resourceGroupName string) DscCompilationJobClient {
	return original.NewDscCompilationJobClient(subscriptionID, resourceGroupName)
}
func NewDscCompilationJobClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) DscCompilationJobClient {
	return original.NewDscCompilationJobClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewDscNodeClient(subscriptionID string, resourceGroupName string) DscNodeClient {
	return original.NewDscNodeClient(subscriptionID, resourceGroupName)
}
func NewDscNodeClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) DscNodeClient {
	return original.NewDscNodeClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewModuleClient(subscriptionID string, resourceGroupName string) ModuleClient {
	return original.NewModuleClient(subscriptionID, resourceGroupName)
}
func NewModuleClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ModuleClient {
	return original.NewModuleClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewScheduleClient(subscriptionID string, resourceGroupName string) ScheduleClient {
	return original.NewScheduleClient(subscriptionID, resourceGroupName)
}
func NewScheduleClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ScheduleClient {
	return original.NewScheduleClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewAgentRegistrationInformationClient(subscriptionID string, resourceGroupName string) AgentRegistrationInformationClient {
	return original.NewAgentRegistrationInformationClient(subscriptionID, resourceGroupName)
}
func NewAgentRegistrationInformationClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) AgentRegistrationInformationClient {
	return original.NewAgentRegistrationInformationClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewConnectionClient(subscriptionID string, resourceGroupName string) ConnectionClient {
	return original.NewConnectionClient(subscriptionID, resourceGroupName)
}
func NewConnectionClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ConnectionClient {
	return original.NewConnectionClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewDscConfigurationClient(subscriptionID string, resourceGroupName string) DscConfigurationClient {
	return original.NewDscConfigurationClient(subscriptionID, resourceGroupName)
}
func NewDscConfigurationClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) DscConfigurationClient {
	return original.NewDscConfigurationClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewDscNodeConfigurationClient(subscriptionID string, resourceGroupName string) DscNodeConfigurationClient {
	return original.NewDscNodeConfigurationClient(subscriptionID, resourceGroupName)
}
func NewDscNodeConfigurationClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) DscNodeConfigurationClient {
	return original.NewDscNodeConfigurationClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewObjectDataTypesClient(subscriptionID string, resourceGroupName string) ObjectDataTypesClient {
	return original.NewObjectDataTypesClient(subscriptionID, resourceGroupName)
}
func NewObjectDataTypesClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ObjectDataTypesClient {
	return original.NewObjectDataTypesClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewRunbookClient(subscriptionID string, resourceGroupName string) RunbookClient {
	return original.NewRunbookClient(subscriptionID, resourceGroupName)
}
func NewRunbookClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) RunbookClient {
	return original.NewRunbookClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}
func New(subscriptionID string, resourceGroupName string) BaseClient {
	return original.New(subscriptionID, resourceGroupName)
}
func NewWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewOperationsClient(subscriptionID string, resourceGroupName string) OperationsClient {
	return original.NewOperationsClient(subscriptionID, resourceGroupName)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewVariableClient(subscriptionID string, resourceGroupName string) VariableClient {
	return original.NewVariableClient(subscriptionID, resourceGroupName)
}
func NewVariableClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) VariableClient {
	return original.NewVariableClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewConnectionTypeClient(subscriptionID string, resourceGroupName string) ConnectionTypeClient {
	return original.NewConnectionTypeClient(subscriptionID, resourceGroupName)
}
func NewConnectionTypeClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ConnectionTypeClient {
	return original.NewConnectionTypeClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewStatisticsClient(subscriptionID string, resourceGroupName string) StatisticsClient {
	return original.NewStatisticsClient(subscriptionID, resourceGroupName)
}
func NewStatisticsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) StatisticsClient {
	return original.NewStatisticsClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewTestJobStreamsClient(subscriptionID string, resourceGroupName string) TestJobStreamsClient {
	return original.NewTestJobStreamsClient(subscriptionID, resourceGroupName)
}
func NewTestJobStreamsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) TestJobStreamsClient {
	return original.NewTestJobStreamsClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewCredentialClient(subscriptionID string, resourceGroupName string) CredentialClient {
	return original.NewCredentialClient(subscriptionID, resourceGroupName)
}
func NewCredentialClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) CredentialClient {
	return original.NewCredentialClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewJobStreamClient(subscriptionID string, resourceGroupName string) JobStreamClient {
	return original.NewJobStreamClient(subscriptionID, resourceGroupName)
}
func NewJobStreamClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) JobStreamClient {
	return original.NewJobStreamClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewNodeReportsClient(subscriptionID string, resourceGroupName string) NodeReportsClient {
	return original.NewNodeReportsClient(subscriptionID, resourceGroupName)
}
func NewNodeReportsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) NodeReportsClient {
	return original.NewNodeReportsClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewWebhookClient(subscriptionID string, resourceGroupName string) WebhookClient {
	return original.NewWebhookClient(subscriptionID, resourceGroupName)
}
func NewWebhookClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) WebhookClient {
	return original.NewWebhookClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewAccountClient(subscriptionID string, resourceGroupName string) AccountClient {
	return original.NewAccountClient(subscriptionID, resourceGroupName)
}
func NewAccountClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) AccountClient {
	return original.NewAccountClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewJobClient(subscriptionID string, resourceGroupName string) JobClient {
	return original.NewJobClient(subscriptionID, resourceGroupName)
}
func NewJobClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) JobClient {
	return original.NewJobClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
func NewRunbookDraftClient(subscriptionID string, resourceGroupName string) RunbookDraftClient {
	return original.NewRunbookDraftClient(subscriptionID, resourceGroupName)
}
func NewRunbookDraftClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) RunbookDraftClient {
	return original.NewRunbookDraftClientWithBaseURI(baseURI, subscriptionID, resourceGroupName)
}
