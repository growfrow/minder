// SPDX-FileCopyrightText: Copyright 2024 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package properties

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	go_github "github.com/google/go-github/v63/github"
	"github.com/rs/zerolog"

	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	"github.com/mindersec/minder/pkg/entities/properties"
	v1 "github.com/mindersec/minder/pkg/providers/v1"
)

const (
	// RepoPropertyId represents the github repository ID (numerical)
	RepoPropertyId = "github/repo_id"
	// RepoPropertyName represents the github repository name
	RepoPropertyName = "github/repo_name"
	// RepoPropertyOwner represents the github repository owner
	RepoPropertyOwner = "github/repo_owner"
	// RepoPropertyDeployURL represents the github repository deployment URL
	RepoPropertyDeployURL = "github/deploy_url"
	// RepoPropertyCloneURL represents the github repository clone URL
	RepoPropertyCloneURL = "github/clone_url"
	// RepoPropertyDefaultBranch represents the github repository default branch
	RepoPropertyDefaultBranch = "github/default_branch"
	// RepoPropertyLicense represents the github repository license
	RepoPropertyLicense = "github/license"
	// RepoPropertyPrimaryLanguage represents the github repository language
	RepoPropertyPrimaryLanguage = "github/primary_language"

	// RepoPropertyHookId represents the github repository hook ID
	RepoPropertyHookId = "github/hook_id"
	// RepoPropertyHookUrl represents the github repository hook URL
	RepoPropertyHookUrl = "github/hook_url"
	// RepoPropertyHookName represents the github repository hook name
	RepoPropertyHookName = "github/hook_name"
	// RepoPropertyHookType represents the github repository hook type
	RepoPropertyHookType = "github/hook_type"
	// RepoPropertyHookUiid represents the github repository hook UIID
	RepoPropertyHookUiid = "github/hook_uiid"
)

var repoOperationalProperties = []string{
	RepoPropertyHookId,
	RepoPropertyHookUrl,
}

var repoPropertyDefinitions = []propertyOrigin{
	{
		keys: []string{
			// general entity
			properties.PropertyName,
			properties.PropertyUpstreamID,
			// general repo
			properties.RepoPropertyIsPrivate,
			properties.RepoPropertyIsArchived,
			properties.RepoPropertyIsFork,
			// github-specific
			RepoPropertyId,
			RepoPropertyName,
			RepoPropertyOwner,
			RepoPropertyDeployURL,
			RepoPropertyCloneURL,
			RepoPropertyDefaultBranch,
			RepoPropertyLicense,
			RepoPropertyPrimaryLanguage,
		},
		wrapper: getRepoWrapper,
	},
}

// GitHubRepoToMap converts a github repository to a map
func GitHubRepoToMap(repo *go_github.Repository) map[string]any {
	repoProps := map[string]any{
		// general entity
		properties.PropertyUpstreamID: properties.NumericalValueToUpstreamID(repo.GetID()),
		// general repo
		properties.RepoPropertyIsPrivate:  repo.GetPrivate(),
		properties.RepoPropertyIsArchived: repo.GetArchived(),
		properties.RepoPropertyIsFork:     repo.GetFork(),
		// github-specific
		RepoPropertyId:              repo.GetID(),
		RepoPropertyName:            repo.GetName(),
		RepoPropertyOwner:           repo.GetOwner().GetLogin(),
		RepoPropertyDeployURL:       repo.GetDeploymentsURL(),
		RepoPropertyCloneURL:        repo.GetCloneURL(),
		RepoPropertyDefaultBranch:   repo.GetDefaultBranch(),
		RepoPropertyLicense:         repo.GetLicense().GetSPDXID(),
		RepoPropertyPrimaryLanguage: repo.GetLanguage(),
	}

	repoProps[properties.PropertyName] = fmt.Sprintf("%s/%s", repo.GetOwner().GetLogin(), repo.GetName())

	return repoProps
}

func getRepoWrapper(
	ctx context.Context, ghCli *go_github.Client, isOrg bool, getByProps *properties.Properties,
) (map[string]any, error) {
	_ = isOrg

	name, owner, err := getNameOwnerFromProps(ctx, getByProps)
	if err != nil {
		return nil, fmt.Errorf("error getting name and owner from properties: %w", err)
	}
	zerolog.Ctx(ctx).Debug().Str("name", name).Str("owner", owner).Msg("Fetching repository")

	repo, result, err := ghCli.Repositories.Get(ctx, owner, name)
	if err != nil {
		if result != nil && result.StatusCode == http.StatusNotFound {
			return nil, v1.ErrEntityNotFound
		}
		return nil, err
	}

	return GitHubRepoToMap(repo), nil
}

func getNameOwnerFromProps(ctx context.Context, props *properties.Properties) (string, string, error) {
	repoNameP := props.GetProperty(RepoPropertyName)
	repoOwnerP := props.GetProperty(RepoPropertyOwner)
	if repoNameP != nil && repoOwnerP != nil {
		zerolog.Ctx(ctx).Debug().Msg("returning repo properties directly")
		return repoNameP.GetString(), repoOwnerP.GetString(), nil
	}

	repoNameP = props.GetProperty(properties.PropertyName)
	if repoNameP != nil {
		zerolog.Ctx(ctx).Debug().Msg("parsing the name")
		slice := strings.Split(repoNameP.GetString(), "/")
		if len(slice) != 2 {
			return "", "", errors.New("invalid repo name")
		}

		return slice[1], slice[0], nil
	}

	return "", "", errors.New("missing required properties, either repo-name and repo-owner or name")
}

// RepositoryFetcher is a property fetcher for github repositories
type RepositoryFetcher struct {
	propertyFetcherBase
}

// NewRepositoryFetcher creates a new RepositoryFetcher
func NewRepositoryFetcher() *RepositoryFetcher {
	return &RepositoryFetcher{
		propertyFetcherBase: propertyFetcherBase{
			operationalProperties: repoOperationalProperties,
			propertyOrigins:       repoPropertyDefinitions,
		},
	}
}

// GetName returns the name of the repository
func (*RepositoryFetcher) GetName(props *properties.Properties) (string, error) {
	repoNameP := props.GetProperty(RepoPropertyName)
	repoOwnerP := props.GetProperty(RepoPropertyOwner)

	if repoNameP == nil || repoOwnerP == nil {
		return "", errors.New("missing required properties")
	}

	repoName := repoNameP.GetString()
	if repoName == "" {
		return "", errors.New("missing required repo-name property value")
	}

	repoOwner := repoOwnerP.GetString()
	if repoOwner == "" {
		return "", errors.New("missing required repo-owner property value")
	}

	return fmt.Sprintf("%s/%s", repoOwner, repoName), nil
}

// RepoV1FromProperties creates a minderv1.Repository from a properties.Properties
func RepoV1FromProperties(repoProperties *properties.Properties) (*minderv1.Repository, error) {
	name, err := repoProperties.GetProperty(RepoPropertyName).AsString()
	if err != nil {
		return nil, fmt.Errorf("error fetching name property: %w", err)
	}

	owner, err := repoProperties.GetProperty(RepoPropertyOwner).AsString()
	if err != nil {
		return nil, fmt.Errorf("error fetching owner property: %w", err)
	}

	repoId, err := repoProperties.GetProperty(RepoPropertyId).AsInt64()
	if err != nil {
		return nil, fmt.Errorf("error fetching repo_id property: %w", err)
	}

	isPrivate, err := repoProperties.GetProperty(properties.RepoPropertyIsPrivate).AsBool()
	if err != nil {
		return nil, fmt.Errorf("error fetching is_private property: %w", err)
	}

	isFork, err := repoProperties.GetProperty(properties.RepoPropertyIsFork).AsBool()
	if err != nil {
		return nil, fmt.Errorf("error fetching is_fork property: %w", err)
	}

	pbRepo := &minderv1.Repository{
		Name:          name,
		Owner:         owner,
		RepoId:        repoId,
		HookId:        repoProperties.GetProperty(RepoPropertyHookId).GetInt64(),
		HookUrl:       repoProperties.GetProperty(RepoPropertyHookUrl).GetString(),
		DeployUrl:     repoProperties.GetProperty(RepoPropertyDeployURL).GetString(),
		CloneUrl:      repoProperties.GetProperty(RepoPropertyCloneURL).GetString(),
		HookType:      repoProperties.GetProperty(RepoPropertyHookType).GetString(),
		HookName:      repoProperties.GetProperty(RepoPropertyHookName).GetString(),
		HookUuid:      repoProperties.GetProperty(RepoPropertyHookUiid).GetString(),
		IsPrivate:     isPrivate,
		IsFork:        isFork,
		DefaultBranch: repoProperties.GetProperty(RepoPropertyDefaultBranch).GetString(),
		License:       repoProperties.GetProperty(RepoPropertyLicense).GetString(),
		Properties:    repoProperties.ToProtoStruct(),
	}

	return pbRepo, nil
}
