package gluadns

import (
	glua "github.com/yuin/gopher-lua"
)

// Preload adds base64 to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local b64 = require("base64")
func Preload(L *glua.LState) {
	L.PreloadModule("dns", Loader)
}

// Loader is the module loader function.
func Loader(L *glua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}
