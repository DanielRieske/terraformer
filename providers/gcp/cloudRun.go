// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcp

import (
	"context"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/run/v2"
	"log"
	"strings"
)

var cloudRunServicesAllowEmptyValues = []string{"labels."}

var cloudRunServicesAdditionalFields = map[string]interface{}{}

type CloudRunGenerator struct {
	GCPService
}

func (g CloudRunGenerator) createRunV2ServiceResources(ctx context.Context, runServicesList *run.ProjectsLocationsServicesListCall) []terraformutils.Resource {
	resources := []terraformutils.Resource{}
	if err := runServicesList.Pages(ctx, func(page *run.GoogleCloudRunV2ListServicesResponse) error {
		for _, obj := range page.Services {
			t := strings.Split(obj.Name, "/")
			name := t[len(t)-1]
			resource := terraformutils.NewResource(
				name,
				obj.Name,
				"google_cloud_run_v2_service",
				g.ProviderName,
				map[string]string{
					"name":     name,
					"project":  g.GetArgs()["project"].(string),
					"location": g.GetArgs()["region"].(compute.Region).Name,
				},
				cloudRunServicesAllowEmptyValues,
				cloudRunServicesAdditionalFields,
			)
			resources = append(resources, resource)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}
	return resources
}

func (g CloudRunGenerator) createRunV2JobResources(ctx context.Context, runServicesList *run.ProjectsLocationsJobsListCall) []terraformutils.Resource {
	resources := []terraformutils.Resource{}
	if err := runServicesList.Pages(ctx, func(page *run.GoogleCloudRunV2ListJobsResponse) error {
		for _, obj := range page.Jobs {
			t := strings.Split(obj.Name, "/")
			name := t[len(t)-1]
			resource := terraformutils.NewResource(
				name,
				obj.Name,
				"google_cloud_run_v2_job",
				g.ProviderName,
				map[string]string{
					"name":     name,
					"project":  g.GetArgs()["project"].(string),
					"location": g.GetArgs()["region"].(compute.Region).Name,
				},
				cloudRunServicesAllowEmptyValues,
				cloudRunServicesAdditionalFields,
			)
			resources = append(resources, resource)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}
	return resources
}

func (g *CloudRunGenerator) InitResources() error {
	ctx := context.Background()
	runService, err := run.NewService(ctx)
	if err != nil {
		return err
	}

	runServiceList := runService.Projects.Locations.Services.List("projects/" + g.GetArgs()["project"].(string) + "/locations/" + g.GetArgs()["region"].(compute.Region).Name)
	runServiceResources := g.createRunV2ServiceResources(ctx, runServiceList)

	runJobList := runService.Projects.Locations.Jobs.List("projects/" + g.GetArgs()["project"].(string) + "/locations/" + g.GetArgs()["region"].(compute.Region).Name)
	runJobResources := g.createRunV2JobResources(ctx, runJobList)

	g.Resources = append(g.Resources, runServiceResources...)
	g.Resources = append(g.Resources, runJobResources...)

	return nil
}
