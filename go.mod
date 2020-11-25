module github.com/dodo-cli/dodo-build

go 1.15

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200512144102-f13ba8f2f2fd
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200310163718-4634ce647cf2+incompatible
	github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)

require (
	github.com/containerd/console v1.0.0
	github.com/docker/docker v17.12.0-ce-rc1.0.20200531234253-77e06fda0c94+incompatible
	github.com/dodo-cli/dodo-core v0.0.0-20200821135148-cb332de21be2
	github.com/dodo-cli/dodo-docker v0.0.0-20200819134644-596e4191c197
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-hclog v0.14.1
	github.com/moby/buildkit v0.7.1
	github.com/oclaussen/go-gimme/configfiles v0.0.0-20200205175519-d9560e60c720
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tonistiigi/fsutil v0.0.0-20200512175118-ae3a8d753069
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0
)
