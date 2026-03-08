module github.com/forgelogic/app

go 1.25.1

replace github.com/forgelogic/nojs => ../

replace github.com/forgelogic/nojs-router => ../router

require (
	github.com/forgelogic/nojs v0.0.0-00010101000000-000000000000
	github.com/forgelogic/nojs-router v0.0.0-00010101000000-000000000000
)
