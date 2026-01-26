package fx

import (
	"go.uber.org/fx"
)

var AppModule = fx.Options(
	configModule,
	infraModule,
	DomainModule,
	RoutesModule,
	ServerModule,
)
