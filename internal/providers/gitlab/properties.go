// SPDX-FileCopyrightText: Copyright 2024 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package gitlab

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	"github.com/mindersec/minder/pkg/entities/properties"
)

// Repository Properties
const (
	// RepoPropertyProjectName represents the gitlab project
	RepoPropertyProjectName = "gitlab/project_name"
	// RepoPropertyDefaultBranch represents the gitlab default branch
	RepoPropertyDefaultBranch = "gitlab/default_branch"
	// RepoPropertyNamespace represents the gitlab repo namespace
	RepoPropertyNamespace = "gitlab/namespace"
	// RepoPropertyLicense represents the gitlab repo license
	RepoPropertyLicense = "gitlab/license"
	// RepoPropertyCloneURL represents the gitlab repo clone URL
	RepoPropertyCloneURL = "gitlab/clone_url"
	// RepoPropertyHookID represents the gitlab repo hook ID
	RepoPropertyHookID = "gitlab/hook_id"
	// RepoPropertyHookURL represents the gitlab repo hook URL
	RepoPropertyHookURL = "gitlab/hook_url"
)

// Pull Request Properties
const (
	// PullRequestProjectID represents the gitlab project ID
	PullRequestProjectID = "gitlab/project_id"
	// PullRequestNumber represents the gitlab merge request number
	PullRequestNumber = "gitlab/merge_request_number"
	// PullRequestAuthor represents the gitlab author
	PullRequestAuthor = "gitlab/author"
)

// Release Properties
const (
	// ReleasePropertyProjectID represents the gitlab project ID
	ReleasePropertyProjectID = "gitlab/project_id"
	// ReleasePropertyTag represents the gitlab release tag name.
	// NOTE: This is used for release discovery, not for creating releases.
	ReleasePropertyTag = "gitlab/tag"
	// ReleasePropertyBranch represents the gitlab release branch
	ReleasePropertyBranch = "gitlab/branch"
)

// FetchAllProperties implements the provider interface
func (c *gitlabClient) FetchAllProperties(
	ctx context.Context, getByProps *properties.Properties, entType minderv1.Entity, _ *properties.Properties,
) (*properties.Properties, error) {
	if !c.SupportsEntity(entType) {
		return nil, fmt.Errorf("entity type %s not supported", entType)
	}

	//nolint:exhaustive // We only support two entity types for now.
	switch entType {
	case minderv1.Entity_ENTITY_REPOSITORIES:
		return c.getPropertiesForRepo(ctx, getByProps)
	case minderv1.Entity_ENTITY_PULL_REQUESTS:
		return c.getPropertiesForPullRequest(ctx, getByProps)
	case minderv1.Entity_ENTITY_RELEASE:
		return c.getPropertiesForRelease(ctx, getByProps)
	default:
		return nil, fmt.Errorf("entity type %s not supported", entType)
	}
}

// FetchProperty implements the provider interface
// TODO: Implement this
func (*gitlabClient) FetchProperty(
	_ context.Context, _ *properties.Properties, _ minderv1.Entity, _ string) (*properties.Property, error) {
	return nil, nil
}

// GetEntityName implements the provider interface
func (c *gitlabClient) GetEntityName(entityType minderv1.Entity, props *properties.Properties) (string, error) {
	if props == nil {
		return "", errors.New("properties are nil")
	}

	if !c.SupportsEntity(entityType) {
		return "", fmt.Errorf("entity type %s not supported", entityType)
	}

	//nolint:exhaustive // We only support two entity types for now.
	switch entityType {
	case minderv1.Entity_ENTITY_REPOSITORIES:
		return getRepoNameFromProperties(props)
	case minderv1.Entity_ENTITY_PULL_REQUESTS:
		return getPullRequestNameFromProperties(props)
	case minderv1.Entity_ENTITY_RELEASE:
		return getReleaseNameFromProperties(props)
	default:
		return "", fmt.Errorf("entity type %s not supported", entityType)
	}
}

// PropertiesToProtoMessage implements the ProtoMessageConverter interface
func (c *gitlabClient) PropertiesToProtoMessage(
	entType minderv1.Entity, props *properties.Properties,
) (protoreflect.ProtoMessage, error) {
	if !c.SupportsEntity(entType) {
		return nil, fmt.Errorf("entity type %s is not supported by the gitlab provider", entType)
	}

	//nolint:exhaustive // We only support two entity types for now.
	switch entType {
	case minderv1.Entity_ENTITY_REPOSITORIES:
		return repoV1FromProperties(props)
	case minderv1.Entity_ENTITY_PULL_REQUESTS:
		return pullRequestV1FromProperties(props)
	case minderv1.Entity_ENTITY_RELEASE:
		return releaseEntityV1FromProperties(props)
	default:
		return nil, fmt.Errorf("entity type %s not supported", entType)
	}
}

func getStringProp(props *properties.Properties, key string) (string, error) {
	value, err := props.GetProperty(key).AsString()
	if err != nil {
		return "", fmt.Errorf("property %s not found or not a string", key)
	}

	return value, nil
}
