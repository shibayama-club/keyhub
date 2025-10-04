package healthcheck

import "context"

type (
	HealthChecker interface {
		Ping(ctx context.Context) error
		Name() string
	}
	HealthCheckFunc struct {
		name string
		fn   func(ctx context.Context) error
	}
)

func NewHealthCheckFunc(name string, fn func(ctx context.Context) error) HealthCheckFunc {
	return HealthCheckFunc{
		name: name,
		fn:   fn,
	}
}

func (f HealthCheckFunc) Ping(ctx context.Context) error {
	return f.fn(ctx)
}

func (f HealthCheckFunc) Name() string {
	return f.name
}
