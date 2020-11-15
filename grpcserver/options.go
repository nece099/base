package grpcserver

var (
	defaultOptions = &options{
		unaryRecoveryHandlerFunc:  nil,
		unaryRecoveryHandlerFunc2: nil,
		streamRecoveryHandlerFunc: nil,
	}
)

type options struct {
	unaryRecoveryHandlerFunc  UnaryRecoveryHandlerFunc
	unaryRecoveryHandlerFunc2 UnaryRecoveryHandlerFunc
	streamRecoveryHandlerFunc StreamRecoveryHandlerFunc
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// WithUnaryRecoveryHandler customizes the function for recovering from a panic.
func WithUnaryRecoveryHandler(f UnaryRecoveryHandlerFunc) Option {
	return func(o *options) {
		o.unaryRecoveryHandlerFunc = f
	}
}

func WithUnaryRecoveryHandler2(f UnaryRecoveryHandlerFunc) Option {
	return func(o *options) {
		o.unaryRecoveryHandlerFunc2 = f
	}
}

// WithStreamRecoveryHandler customizes the function for recovering from a panic.
func WithStreamRecoveryHandler(f StreamRecoveryHandlerFunc) Option {
	return func(o *options) {
		o.streamRecoveryHandlerFunc = f
	}
}
