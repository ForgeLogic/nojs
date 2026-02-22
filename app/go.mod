module github.com/ForgeLogic/app

go 1.25.1

replace github.com/ForgeLogic/nojs => ../nojs

replace github.com/ForgeLogic/nojs-router => ../router

require (
	github.com/ForgeLogic/nojs v0.0.0-00010101000000-000000000000
	github.com/ForgeLogic/nojs-router v0.0.0-00010101000000-000000000000
)
