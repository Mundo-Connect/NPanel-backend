//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build

package main

import (
	"github.com/npanel-dev/NPanel-backend/internal/biz"
	"github.com/npanel-dev/NPanel-backend/internal/conf"
	"github.com/npanel-dev/NPanel-backend/internal/data"
	"github.com/npanel-dev/NPanel-backend/internal/server"
	"github.com/npanel-dev/NPanel-backend/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application
func wireApp(*conf.Server, *conf.Data, *conf.Application, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
