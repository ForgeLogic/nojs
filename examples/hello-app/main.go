//go:build js || wasm

package main

import (
	"helloapp/components"

	router "github.com/ForgeLogic/nojs-router"
	"github.com/ForgeLogic/nojs/runtime"
)

// Each component type needs a unique ID for the router's diffing algorithm.
const helloAppTypeID uint32 = 100

func main() {
	// 1. Create the router engine.
	routerEngine := router.NewEngine(nil)

	// 2. Create the renderer — it mounts the app into the <div id="app"> element.
	renderer := runtime.NewRenderer(routerEngine, "#app")
	routerEngine.SetRenderer(renderer)

	// 3. Register routes: map URL paths to component factories.
	routerEngine.RegisterRoutes([]router.Route{
		{
			Path: "/",
			Chain: []router.ComponentMetadata{
				{
					Factory: func(p map[string]string) runtime.Component {
						return &components.HelloApp{}
					},
					TypeID: helloAppTypeID,
				},
			},
		},
	})

	// 4. Create an AppShell (no persistent layout — renders the component directly).
	appShell := router.NewAppShell(nil)
	renderer.SetCurrentComponent(appShell, "app-shell")
	renderer.ReRender()

	// 5. Start the router. The callback is called on every navigation.
	err := routerEngine.Start(func(chain []runtime.Component, key string) {
		appShell.SetPage(chain, key)
	})
	if err != nil {
		panic(err)
	}

	// Keep the Go program alive.
	select {}
}
