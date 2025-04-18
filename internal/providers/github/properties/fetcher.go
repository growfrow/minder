// SPDX-FileCopyrightText: Copyright 2024 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

// Package properties provides utility functions for fetching and managing properties
package properties

import (
	"context"

	go_github "github.com/google/go-github/v63/github"

	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	"github.com/mindersec/minder/pkg/entities/properties"
)

// GhPropertyWrapper is a function that fetches a property from the GitHub API
type GhPropertyWrapper func(
	ctx context.Context, ghCli *go_github.Client, isOrg bool, lookupProperties *properties.Properties,
) (map[string]any, error)

// GhPropertyFetcher is an interface for fetching properties from the GitHub API
type GhPropertyFetcher interface {
	ghPropertyFetcherBase
	GetName(props *properties.Properties) (string, error)
}

type ghPropertyFetcherBase interface {
	WrapperForProperty(propertyKey string) GhPropertyWrapper
	AllPropertyWrappers() []GhPropertyWrapper
	OperationalProperties() []string
}

// GhPropertyFetcherFactory is an interface for creating GhPropertyFetcher instances
type GhPropertyFetcherFactory interface {
	EntityPropertyFetcher(entType minderv1.Entity) GhPropertyFetcher
}

type ghEntityFetcher struct{}

// NewPropertyFetcherFactory creates a new GhPropertyFetcherFactory
func NewPropertyFetcherFactory() GhPropertyFetcherFactory {
	return ghEntityFetcher{}
}

// EntityPropertyFetcher returns a GhPropertyFetcher for the given entity type
func (ghEntityFetcher) EntityPropertyFetcher(entType minderv1.Entity) GhPropertyFetcher {
	// nolint:exhaustive
	switch entType {
	case minderv1.Entity_ENTITY_PULL_REQUESTS:
		return NewPullRequestFetcher()
	case minderv1.Entity_ENTITY_REPOSITORIES:
		return NewRepositoryFetcher()
	case minderv1.Entity_ENTITY_ARTIFACTS:
		return NewArtifactFetcher()
	case minderv1.Entity_ENTITY_RELEASE:
		return NewReleaseFetcher()
	}

	return nil
}

type propertyOrigin struct {
	keys    []string
	wrapper GhPropertyWrapper
}

type propertyFetcherBase struct {
	propertyOrigins       []propertyOrigin
	operationalProperties []string
}

var _ ghPropertyFetcherBase = propertyFetcherBase{}

// WrapperForProperty returns the property wrapper for the given property key
func (pfb propertyFetcherBase) WrapperForProperty(propertyKey string) GhPropertyWrapper {
	for _, po := range pfb.propertyOrigins {
		for _, k := range po.keys {
			if k == propertyKey {
				return po.wrapper
			}
		}
	}

	return nil
}

// AllPropertyWrappers returns all property wrappers for the repository
func (pfb propertyFetcherBase) AllPropertyWrappers() []GhPropertyWrapper {
	wrappers := make([]GhPropertyWrapper, 0, len(pfb.propertyOrigins))
	for _, po := range pfb.propertyOrigins {
		wrappers = append(wrappers, po.wrapper)
	}
	return wrappers
}

// OperationalProperties returns the operational properties for the repository
func (pfb propertyFetcherBase) OperationalProperties() []string {
	return pfb.operationalProperties
}
