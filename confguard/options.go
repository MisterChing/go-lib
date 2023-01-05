package confguard

import "github.com/go-kratos/kratos/v2/log"

type Option func(*options)

type options struct {
	group    string
	dataID   string
	watchKey string
	guarder  *Guarder
	logger   log.Logger
}

func WithGroup(group string) Option {
	return func(o *options) {
		o.group = group
	}
}

func WithDataID(dataID string) Option {
	return func(o *options) {
		o.dataID = dataID
	}
}

func WithWatchKey(key string) Option {
	return func(o *options) {
		o.watchKey = key
	}
}

func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithGuarder(guarder *Guarder) Option {
	return func(o *options) {
		o.guarder = guarder
	}
}
