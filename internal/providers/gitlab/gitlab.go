// SPDX-FileCopyrightText: Copyright 2024 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

// Package gitlab provides the GitLab OAuth provider implementation
package gitlab

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mindersec/minder/internal/db"
	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	config "github.com/mindersec/minder/pkg/config/server"
	provifv1 "github.com/mindersec/minder/pkg/providers/v1"
)

// Class is the string that represents the GitLab provider class
const Class = "gitlab"

// Implements is the list of provider types that the DockerHub provider implements
var Implements = []db.ProviderType{
	db.ProviderTypeGit,
	db.ProviderTypeRest,
	db.ProviderTypeRepoLister,
}

// AuthorizationFlows is the list of authorization flows that the DockerHub provider supports
var AuthorizationFlows = []db.AuthorizationFlow{
	db.AuthorizationFlowUserInput,
	db.AuthorizationFlowOauth2AuthorizationCodeFlow,
}

// Ensure that the GitLab provider implements the right interfaces
var _ provifv1.Git = (*gitlabClient)(nil)
var _ provifv1.REST = (*gitlabClient)(nil)
var _ provifv1.RepoLister = (*gitlabClient)(nil)

type gitlabClient struct {
	cred       provifv1.GitLabCredential
	cli        *http.Client
	glcfg      *minderv1.GitLabProviderConfig
	webhookURL string
	gitConfig  config.GitConfig

	// secret for the webhook. This is stored in the
	// structure to allow efficient fetching.
	currentWebhookSecret string
}

// New creates a new GitLab provider
// Note that the webhook URL should already contain the provider class in the path
func New(
	cred provifv1.GitLabCredential,
	cfg *minderv1.GitLabProviderConfig,
	webhookURL string,
	currentWebhookSecret string,
) (*gitlabClient, error) {
	// TODO: We need a context here.
	cli := oauth2.NewClient(context.Background(), cred.GetAsOAuth2TokenSource())

	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://gitlab.com/api/v4/"
	}

	if webhookURL == "" {
		return nil, errors.New("webhook URL is required")
	}

	return &gitlabClient{
		cred:                 cred,
		cli:                  cli,
		glcfg:                cfg,
		webhookURL:           webhookURL,
		currentWebhookSecret: currentWebhookSecret,
		// TODO: Add git config
	}, nil
}

type glConfigWrapper struct {
	GitLab *minderv1.GitLabProviderConfig `json:"gitlab" yaml:"gitlab" mapstructure:"gitlab" validate:"required"`
}

// ParseV1Config parses the raw configuration into a GitLabProviderConfig
//
// TODO: This should be moved to a common location
func ParseV1Config(rawCfg json.RawMessage) (*minderv1.GitLabProviderConfig, error) {
	var cfg glConfigWrapper
	if err := json.Unmarshal(rawCfg, &cfg); err != nil {
		return nil, err
	}

	if cfg.GitLab == nil {
		// Return a default but working config
		return &minderv1.GitLabProviderConfig{}, nil
	}

	return cfg.GitLab, nil
}

// MarshalV1Config marshals and validates the given config
// so it can safely be stored in the database
func MarshalV1Config(rawCfg json.RawMessage) (json.RawMessage, error) {
	var w glConfigWrapper
	if err := json.Unmarshal(rawCfg, &w); err != nil {
		return nil, err
	}

	// TODO: Add validation
	// err := w.GitLab.Validate()
	// if err != nil {
	// 	return nil, fmt.Errorf("error validating gitlab config: %w", err)
	// }

	return json.Marshal(w)
}

// CanImplement returns true if the provider can implement the given trait
func (*gitlabClient) CanImplement(trait minderv1.ProviderType) bool {
	return trait == minderv1.ProviderType_PROVIDER_TYPE_GIT ||
		trait == minderv1.ProviderType_PROVIDER_TYPE_REST
}

func (c *gitlabClient) GetCredential() provifv1.GitLabCredential {
	return c.cred
}

// SupportsEntity implements the Provider interface
func (*gitlabClient) SupportsEntity(entType minderv1.Entity) bool {
	return entType == minderv1.Entity_ENTITY_REPOSITORIES ||
		entType == minderv1.Entity_ENTITY_PULL_REQUESTS ||
		entType == minderv1.Entity_ENTITY_RELEASE
}
