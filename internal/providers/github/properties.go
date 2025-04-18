// SPDX-FileCopyrightText: Copyright 2024 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/rs/zerolog"

	properties2 "github.com/mindersec/minder/internal/providers/github/properties"
	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	"github.com/mindersec/minder/pkg/entities/properties"
)

// FetchProperty fetches a single property for the given entity
func (c *GitHub) FetchProperty(
	ctx context.Context, getByProps *properties.Properties, entType minderv1.Entity, key string,
) (*properties.Property, error) {
	if c.propertyFetchers == nil {
		return nil, errors.New("property fetchers not initialized")
	}

	fetcher := c.propertyFetchers.EntityPropertyFetcher(entType)
	if fetcher == nil {
		return nil, fmt.Errorf("entity %s not supported", entType)
	}

	wrapper := fetcher.WrapperForProperty(key)
	if wrapper == nil {
		return nil, fmt.Errorf("property %s not supported for entity %s", key, entType)
	}

	props, err := wrapper(ctx, c.client, c.IsOrg(), getByProps)
	if err != nil {
		return nil, fmt.Errorf("error fetching property %s for entity %s: %w", key, entType, err)
	}
	value, ok := props[key]
	if !ok {
		return nil, errors.New("requested property not found in result")
	}
	return properties.NewProperty(value)
}

// FetchAllProperties fetches all properties for the given entity
func (c *GitHub) FetchAllProperties(
	ctx context.Context, getByProps *properties.Properties, entType minderv1.Entity, cachedProps *properties.Properties,
) (*properties.Properties, error) {
	if c.propertyFetchers == nil {
		return nil, errors.New("property fetchers not initialized")
	}

	zerolog.Ctx(ctx).Debug().
		Str("entity", entType.String()).
		Dict("getByProps", getByProps.ToLogDict()).
		Msg("Fetching all properties")

	fetcher := c.propertyFetchers.EntityPropertyFetcher(entType)
	result := make(map[string]any)
	for _, wrapper := range fetcher.AllPropertyWrappers() {
		props, err := wrapper(ctx, c.client, c.IsOrg(), getByProps)
		if err != nil {
			return nil, fmt.Errorf("error fetching properties for entity %s: %w", entType, err)
		}

		for k, v := range props {
			result[k] = v
		}
	}

	upstreamProps := properties.NewProperties(result)

	operational := filterOperational(cachedProps, fetcher)
	return upstreamProps.Merge(operational), nil
}

func filterOperational(cachedProperties *properties.Properties, fetcher properties2.GhPropertyFetcher) *properties.Properties {
	if cachedProperties == nil {
		// Nothing to filter
		return nil
	}

	operational := fetcher.OperationalProperties()
	if len(operational) == 0 {
		return cachedProperties
	}

	filter := func(key string, _ *properties.Property) bool {
		return slices.Contains(operational, key)
	}

	return cachedProperties.FilteredCopy(filter)
}
