package main

import (
	build "github.com/dodo/dodo-build/pkg/plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
)

func main() {
	build.RegisterPlugin()
	dodo.ServePlugins()
}
