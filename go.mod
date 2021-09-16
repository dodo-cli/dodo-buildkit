module github.com/dodo-cli/dodo-buildkit

go 1.16

replace (
	github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)

require (
	github.com/docker/docker v20.10.2+incompatible
	github.com/dodo-cli/dodo-core v0.1.1-0.20210915164217-1c50006deac2
	github.com/hashicorp/go-hclog v0.15.0
	github.com/jaguilar/vt100 v0.0.0-20150826170717-2703a27b14ea
	github.com/moby/buildkit v0.8.0-rc3
	github.com/morikuni/aec v1.0.0
	github.com/oclaussen/go-gimme/configfiles v0.0.0-20200205175519-d9560e60c720
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/tonistiigi/fsutil v0.0.0-20201103201449-0834f99b7b85
	github.com/tonistiigi/units v0.0.0-20180711220420-6950e57a87ea
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1
)
