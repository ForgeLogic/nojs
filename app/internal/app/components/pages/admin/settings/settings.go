package settings

import "github.com/vcrobe/nojs/runtime"

type Settings struct {
	runtime.ComponentBase
}

func (*Settings) OnMount() {
	println("[Settings] Mounted")
}

func (*Settings) OnUnmount() {
	println("[Settings] Unmounted")
}
