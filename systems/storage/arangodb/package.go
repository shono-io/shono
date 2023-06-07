package arangodb

import (
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/shono-io/shono/systems"
)

func init() {
	err := service.RegisterCache("arangodb", arangodbCacheConfig(), func(conf *service.ParsedConfig, mgr *service.Resources) (service.Cache, error) {
		return cacheFromConfig(conf)
	})
	if err != nil {
		panic(err)
	}

	err = service.RegisterProcessor("arangodb", arangodbProcConfig(), func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
		return procFromConfig(conf)
	})
	if err != nil {
		panic(err)
	}

	systems.RegisterStorageSystem("arangodb", &system{})
}
