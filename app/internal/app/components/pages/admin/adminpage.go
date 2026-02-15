package admin

import "github.com/vcrobe/nojs/runtime"

type AdminPage struct {
	runtime.ComponentBase
}

func (*AdminPage) OnMount() {
	println("[AdminPage] Mounted")
}

func (*AdminPage) OnUnmount() {
	println("[AdminPage] Unmounted")
}
