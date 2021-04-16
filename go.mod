module github.com/dodo-cli/dodo-build

go 1.16

replace (
	github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)

require (
	github.com/containerd/console v1.0.1
	github.com/docker/docker v20.10.2+incompatible
	github.com/dodo-cli/dodo-core v0.0.0-20210318162438-a68f0193ede5
	github.com/hashicorp/go-hclog v0.15.0
	github.com/moby/buildkit v0.8.0-rc3
	github.com/oclaussen/go-gimme/configfiles v0.0.0-20200205175519-d9560e60c720
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/tonistiigi/fsutil v0.0.0-20201103201449-0834f99b7b85
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
)
