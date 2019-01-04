// +build go1.9

// Copyright 2018 Microsoft Corporation
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

package training

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/customvision/training"

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type Classifier = original.Classifier

const (
	Multiclass Classifier = original.Multiclass
	Multilabel Classifier = original.Multilabel
)

type DomainType = original.DomainType

const (
	Classification  DomainType = original.Classification
	ObjectDetection DomainType = original.ObjectDetection
)

type ExportFlavor = original.ExportFlavor

const (
	Linux   ExportFlavor = original.Linux
	Windows ExportFlavor = original.Windows
)

type ExportPlatform = original.ExportPlatform

const (
	CoreML     ExportPlatform = original.CoreML
	DockerFile ExportPlatform = original.DockerFile
	ONNX       ExportPlatform = original.ONNX
	TensorFlow ExportPlatform = original.TensorFlow
)

type ExportStatusModel = original.ExportStatusModel

const (
	Done      ExportStatusModel = original.Done
	Exporting ExportStatusModel = original.Exporting
	Failed    ExportStatusModel = original.Failed
)

type ImageUploadStatus = original.ImageUploadStatus

const (
	ErrorImageFormat       ImageUploadStatus = original.ErrorImageFormat
	ErrorImageSize         ImageUploadStatus = original.ErrorImageSize
	ErrorLimitExceed       ImageUploadStatus = original.ErrorLimitExceed
	ErrorRegionLimitExceed ImageUploadStatus = original.ErrorRegionLimitExceed
	ErrorSource            ImageUploadStatus = original.ErrorSource
	ErrorStorage           ImageUploadStatus = original.ErrorStorage
	ErrorTagLimitExceed    ImageUploadStatus = original.ErrorTagLimitExceed
	ErrorUnknown           ImageUploadStatus = original.ErrorUnknown
	OK                     ImageUploadStatus = original.OK
	OKDuplicate            ImageUploadStatus = original.OKDuplicate
)

type OrderBy = original.OrderBy

const (
	Newest    OrderBy = original.Newest
	Oldest    OrderBy = original.Oldest
	Suggested OrderBy = original.Suggested
)

type BoundingBox = original.BoundingBox
type Domain = original.Domain
type Export = original.Export
type Image = original.Image
type ImageCreateResult = original.ImageCreateResult
type ImageCreateSummary = original.ImageCreateSummary
type ImageFileCreateBatch = original.ImageFileCreateBatch
type ImageFileCreateEntry = original.ImageFileCreateEntry
type ImageIDCreateBatch = original.ImageIDCreateBatch
type ImageIDCreateEntry = original.ImageIDCreateEntry
type ImagePerformance = original.ImagePerformance
type ImagePrediction = original.ImagePrediction
type ImageRegion = original.ImageRegion
type ImageRegionCreateBatch = original.ImageRegionCreateBatch
type ImageRegionCreateEntry = original.ImageRegionCreateEntry
type ImageRegionCreateResult = original.ImageRegionCreateResult
type ImageRegionCreateSummary = original.ImageRegionCreateSummary
type ImageRegionProposal = original.ImageRegionProposal
type ImageTag = original.ImageTag
type ImageTagCreateBatch = original.ImageTagCreateBatch
type ImageTagCreateEntry = original.ImageTagCreateEntry
type ImageTagCreateSummary = original.ImageTagCreateSummary
type ImageURL = original.ImageURL
type ImageURLCreateBatch = original.ImageURLCreateBatch
type ImageURLCreateEntry = original.ImageURLCreateEntry
type Int32 = original.Int32
type Iteration = original.Iteration
type IterationPerformance = original.IterationPerformance
type ListDomain = original.ListDomain
type ListExport = original.ListExport
type ListImage = original.ListImage
type ListImagePerformance = original.ListImagePerformance
type ListIteration = original.ListIteration
type ListProject = original.ListProject
type ListTag = original.ListTag
type Prediction = original.Prediction
type PredictionQueryResult = original.PredictionQueryResult
type PredictionQueryTag = original.PredictionQueryTag
type PredictionQueryToken = original.PredictionQueryToken
type Project = original.Project
type ProjectSettings = original.ProjectSettings
type Region = original.Region
type RegionProposal = original.RegionProposal
type StoredImagePrediction = original.StoredImagePrediction
type Tag = original.Tag
type TagPerformance = original.TagPerformance

func New(aPIKey string) BaseClient {
	return original.New(aPIKey)
}
func NewWithBaseURI(baseURI string, aPIKey string) BaseClient {
	return original.NewWithBaseURI(baseURI, aPIKey)
}
func PossibleClassifierValues() []Classifier {
	return original.PossibleClassifierValues()
}
func PossibleDomainTypeValues() []DomainType {
	return original.PossibleDomainTypeValues()
}
func PossibleExportFlavorValues() []ExportFlavor {
	return original.PossibleExportFlavorValues()
}
func PossibleExportPlatformValues() []ExportPlatform {
	return original.PossibleExportPlatformValues()
}
func PossibleExportStatusModelValues() []ExportStatusModel {
	return original.PossibleExportStatusModelValues()
}
func PossibleImageUploadStatusValues() []ImageUploadStatus {
	return original.PossibleImageUploadStatusValues()
}
func PossibleOrderByValues() []OrderBy {
	return original.PossibleOrderByValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
