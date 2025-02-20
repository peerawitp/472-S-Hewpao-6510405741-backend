package ctx

import "context"

func ProvideContext() context.Context {
	return context.Background()
}
