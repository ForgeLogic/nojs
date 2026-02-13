//go:build js || wasm
// +build js wasm

package main

import (
	"fmt"

	"github.com/vcrobe/app/internal/app/components/pages"
	"github.com/vcrobe/app/internal/app/components/pages/admin"
	"github.com/vcrobe/app/internal/app/components/pages/admin/layouts"
	"github.com/vcrobe/app/internal/app/components/pages/admin/settings"
	sharedlayouts "github.com/vcrobe/app/internal/app/components/shared/layouts"
	"github.com/vcrobe/app/internal/app/context"
	"github.com/vcrobe/nojs/console"
	"github.com/vcrobe/nojs/router"
	"github.com/vcrobe/nojs/runtime"
	"github.com/vcrobe/nojs/vdom"
)

// TypeID constants for route components
const (
	MainLayout_TypeID   uint32 = 0x8F22A1BC
	AdminLayout_TypeID  uint32 = 0x7E11B2AD
	HomePage_TypeID     uint32 = 0x6C00C9FE
	AboutPage_TypeID    uint32 = 0x5D11D8CF
	AdminPage_TypeID    uint32 = 0x4E22E7B0
	SettingsPage_TypeID uint32 = 0x3F33F681
	BlogPage_TypeID     uint32 = 0x2E44F592
)

// AppShell is a stable root component that holds persistent layouts (app shell)
// and swaps only the BodyContent slot when navigation occurs. This preserves
// layout instances and their internal state across navigations including sublayouts.
type AppShell struct {
	runtime.ComponentBase

	// persistent layout instance (app shell)
	mainLayout *sharedlayouts.MainLayout

	// current chain of component instances (all from router, volatile)
	currentChain []runtime.Component
	currentKey   string
}

func NewAppShell(mainLayout *sharedlayouts.MainLayout) *AppShell {
	return &AppShell{
		mainLayout:   mainLayout,
		currentChain: make([]runtime.Component, 0),
	}
}

// SetPage replaces the volatile chain of component instances and triggers a re-render.
// The chain includes components from the router (from pivot onwards).
// When pivot > 0, the chain doesn't include MainLayout (it's preserved).
// The RenderChild mechanism ensures layouts are reused efficiently,
// and VDOM patching only updates what changed.
func (a *AppShell) SetPage(chain []runtime.Component, key string) {
	console.Log("[AppShell.SetPage] Called with", len(chain), "components, key:", key)
	if len(chain) > 0 {
		console.Log("[AppShell.SetPage] First component type:", fmt.Sprintf("%%T", chain[0]))
	}

	// If the chain doesn't include mainLayout at index 0, prepend it
	// (this happens when pivot > 0 and layouts are preserved)
	if len(chain) == 0 || chain[0] != a.mainLayout {
		console.Log("[AppShell.SetPage] Prepending mainLayout to chain")
		fullChain := make([]runtime.Component, 0, len(chain)+1)
		fullChain = append(fullChain, a.mainLayout)
		fullChain = append(fullChain, chain...)
		a.currentChain = fullChain
	} else {
		a.currentChain = chain
	}
	a.currentKey = key

	// Trigger a re-render of AppShell. RenderChild will reuse mainLayout instance,
	// and VDOM patching will only update the changed slot content.
	console.Log("[AppShell.SetPage] Calling StateHasChanged")
	a.StateHasChanged()
}

// Render composes the persistent MainLayout with the current component chain.
// The chain was set up by the router with all intermediate slots already connected.
// We inject only the root of the chain (first component) into MainLayout.BodyContent.
// The router's chain linking ensures all intermediate slot connections work correctly.
func (a *AppShell) Render(r runtime.Renderer) *vdom.VNode {
	console.Log("[AppShell.Render] Called, chain length:", len(a.currentChain))
	// Ensure renderer is injected into children so StateHasChanged / Navigate work.
	type rendererSetter interface {
		SetRenderer(runtime.Renderer)
	}

	// Ensure MainLayout has renderer
	if a.mainLayout != nil {
		if rs, ok := interface{}(a.mainLayout).(rendererSetter); ok {
			rs.SetRenderer(r)
		}
	}

	// Render the root of the chain into a VNode slice for the slot
	var slotChildren []*vdom.VNode
	if len(a.currentChain) > 0 {
		console.Log("[AppShell.Render] Processing chain")
		// Skip MainLayout if it's at index 0 (AppShell manages it separately)
		chainIndex := 0
		if a.currentChain[0] == a.mainLayout {
			console.Log("[AppShell.Render] First component is mainLayout, skipping")
			chainIndex = 1
		}

		// Render the first non-MainLayout component in the chain
		if chainIndex < len(a.currentChain) {
			rootComponent := a.currentChain[chainIndex]
			console.Log("[AppShell.Render] Rendering component at index", chainIndex, "type:", fmt.Sprintf("%T", rootComponent))

			// Give renderer to root component if possible
			if rs, ok := interface{}(rootComponent).(rendererSetter); ok {
				rs.SetRenderer(r)
			}

			// Use RenderChild to track sublayouts and pages for efficient caching/patching.
			// The key includes the type and pointer so different component types don't collide,
			// but the same preserved instance (e.g., AdminLayout across /admin → /admin/settings)
			// gets reused and efficiently patched.
			slotKey := fmt.Sprintf("slot-root-%T-%p", rootComponent, rootComponent)
			childVNode := r.RenderChild(slotKey, rootComponent)
			if childVNode != nil {
				slotChildren = []*vdom.VNode{childVNode}
			}
		}
	}

	// Inject into layout's BodyContent slot (compiler-generated field)
	// Layouts follow the single-slot convention: BodyContent []*vdom.VNode
	if a.mainLayout != nil {
		// assign slot directly; generated ApplyProps will preserve state on instance reuse
		a.mainLayout.BodyContent = slotChildren
		// Use RenderChild to render mainLayout so the renderer caches it for efficient patching
		// on subsequent navigations. The key "main-layout" identifies this component.
		return r.RenderChild("main-layout", a.mainLayout)
	}

	// Fallback: if no layout, render the first non-MainLayout component from chain
	if len(a.currentChain) > 0 {
		// Skip MainLayout if it's at index 0
		chainIndex := 0
		if a.currentChain[0] == a.mainLayout {
			chainIndex = 1
		}
		if chainIndex < len(a.currentChain) {
			rootComponent := a.currentChain[chainIndex]
			slotKey := fmt.Sprintf("slot-root-%T-%p", rootComponent, rootComponent)
			return r.RenderChild(slotKey, rootComponent)
		}
	}

	// Empty fallback
	return vdom.NewVNode("div", nil, nil, "")
}

func main() {
	// Create shared layout context
	mainLayoutCtx := &context.MainLayoutCtx{
		Title: "My App",
	}

	// Create persistent main layout instance (app shell)
	mainLayout := &sharedlayouts.MainLayout{
		MainLayoutCtx: mainLayoutCtx,
	}

	// Create the router engine first (it will be passed as navigation manager to renderer)
	routerEngine := router.NewEngine(nil)

	// Create the renderer with the engine as the navigation manager
	renderer := runtime.NewRenderer(routerEngine, "#app")

	// Set the renderer on the engine so it can render components
	routerEngine.SetRenderer(renderer)

	// Define all routes with layout chains and TypeIDs
	routerEngine.RegisterRoutes([]router.Route{
		{
			Path: "/",
			Chain: []router.ComponentMetadata{
				{
					Factory: func() runtime.Component { return mainLayout },
					TypeID:  MainLayout_TypeID,
				},
				{
					Factory: func() runtime.Component { return &pages.HomePage{MainLayoutCtx: mainLayoutCtx} },
					TypeID:  HomePage_TypeID,
				},
			},
		},
		{
			Path: "/about",
			Chain: []router.ComponentMetadata{
				{
					Factory: func() runtime.Component { return mainLayout },
					TypeID:  MainLayout_TypeID,
				},
				{
					Factory: func() runtime.Component { return &pages.AboutPage{} },
					TypeID:  AboutPage_TypeID,
				},
			},
		},
		{
			Path: "/blog/{year}",
			Chain: []router.ComponentMetadata{
				{
					Factory: func() runtime.Component { return mainLayout },
					TypeID:  MainLayout_TypeID,
				},
				{
					Factory: func() runtime.Component {
						year := 2026 // Default, would be extracted from URL in real implementation
						return &pages.BlogPage{Year: year}
					},
					TypeID: BlogPage_TypeID,
				},
			},
		},
		{
			Path: "/admin",
			Chain: []router.ComponentMetadata{
				{
					Factory: func() runtime.Component { return mainLayout },
					TypeID:  MainLayout_TypeID,
				},
				{
					Factory: func() runtime.Component { return layouts.NewAdminLayout() },
					TypeID:  AdminLayout_TypeID,
				},
				{
					Factory: func() runtime.Component { return &admin.AdminPage{} },
					TypeID:  AdminPage_TypeID,
				},
			},
		},
		{
			Path: "/admin/settings",
			Chain: []router.ComponentMetadata{
				{
					Factory: func() runtime.Component { return mainLayout },
					TypeID:  MainLayout_TypeID,
				},
				{
					Factory: func() runtime.Component { return layouts.NewAdminLayout() },
					TypeID:  AdminLayout_TypeID,
				},
				{
					Factory: func() runtime.Component { return &settings.Settings{} },
					TypeID:  SettingsPage_TypeID,
				},
			},
		},
	})

	// Create AppShell to wrap the router's page rendering
	appShell := NewAppShell(mainLayout)
	renderer.SetCurrentComponent(appShell, "app-shell")
	renderer.ReRender()

	// Initialize the router with a callback to update AppShell when navigation occurs
	if err := routerEngine.Start(func(chain []runtime.Component, key string) {
		appShell.SetPage(chain, key)
	}); err != nil {
		console.Error("Failed to start router:", err.Error())
		panic(err)
	}

	// Keep the Go program running
	select {}
}
