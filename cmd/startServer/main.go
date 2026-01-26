package main

import (
	fxApp "github.com/lkgiovani/go-boilerplate/infra/fx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxApp.AppModule,
	).Run()
}
