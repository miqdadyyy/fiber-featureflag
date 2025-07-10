package providers

import "context"

type IFeatureFlagProvider interface {
	SetFeatureFlagStatus(ctx context.Context, key string, value bool) error
	GetFeatureFlagStatus(ctx context.Context, key string) bool
	GetListOfFeatureFlags(ctx context.Context) (map[string]bool, error)
}
