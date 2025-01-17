// Copyright 2023 Google LLC
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

package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"golang.org/x/oauth2/google"
	alpha "google.golang.org/api/compute/v0.alpha"
	beta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func newCloud(project string) (cloud.Cloud, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return nil, err
	}

	alpha, err := alpha.New(client)
	if err != nil {
		return nil, err
	}
	beta, err := beta.New(client)
	if err != nil {
		return nil, err
	}
	ga, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	svc := &cloud.Service{
		GA:            ga,
		Alpha:         alpha,
		Beta:          beta,
		ProjectRouter: &cloud.SingleProjectRouter{ID: project},
		RateLimiter:   &cloud.NopRateLimiter{},
	}

	theCloud := cloud.NewGCE(svc)
	return theCloud, nil
}

// verifyCluster checks if the cluster description has the expected configuration.
func verifyCluster(config ClusterConfig) error {
	// Verify if the cluster has the correct currentNodeCount.
	params := []string{
		"container",
		"clusters",
		"describe",
		config.Name,
		"--zone", config.Zone,
		"--format", "value(currentNodeCount)",
	}
	out, err := exec.Command("gcloud", params...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot describe cluster %q using gcloud: %w", config.Name, err)
	}

	outString := strings.TrimSpace(string(out))
	gotNumOfNodes, err := strconv.Atoi(outString)
	if err != nil {
		return fmt.Errorf("failed to convert currentNodeCount %q to int: %w", outString, err)
	}
	if gotNumOfNodes != config.NumOfNodes {
		return fmt.Errorf("expect cluster %s to have %d nodes, got %d", config.Name, config.NumOfNodes, gotNumOfNodes)
	}
	return nil
}
