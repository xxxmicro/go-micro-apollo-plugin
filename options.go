package apollo

import (
	"context"
	"github.com/micro/go-micro/v2/config/source"
)

type addressKey struct{}
type namespaceName struct{}
type clusterKey struct{}
type appIdKey struct{}
type backupConfigPathKey struct{}

func WithNamespace(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, namespaceName{}, c)
	}
}

func WithAddress(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, addressKey{}, c)
	}
}

func WithBackupConfigPath(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, backupConfigPathKey{}, c)
	}
}

func WithAppId(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, appIdKey{}, c)
	}
}

func WithCluster(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, clusterKey{}, c)
	}
}
