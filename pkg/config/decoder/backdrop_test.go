package decoder

import (
	"testing"

	"github.com/dodo/dodo-build/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const fullExample = `
image:
  name: testimage
  context: .
  dockerfile: Dockerfile
  steps:
    - RUN hello
    - RUN world
  args:
    - FOO=BAR
  no_cache: true
  force_rebuild: true
  force_pull: true
`

func TestFullExample(t *testing.T) {
	config := getExampleConfig(t, fullExample)
	assert.NotNil(t, config.Build)
	assert.Equal(t, "testimage", config.Build.ImageName)
	assert.Equal(t, ".", config.Build.Context)
	assert.Equal(t, "Dockerfile", config.Build.Dockerfile)
	assert.Equal(t, []string{"RUN hello", "RUN world"}, config.Build.InlineDockerfile)
	assert.Equal(t, 1, len(config.Build.Arguments))
	assert.Equal(t, "FOO", config.Build.Arguments[0].Key)
	assert.Equal(t, "BAR", config.Build.Arguments[0].Value)
	assert.True(t, config.Build.NoCache)
	assert.True(t, config.Build.ForceRebuild)
	assert.True(t, config.Build.ForcePull)
}

func getExampleConfig(t *testing.T, yamlConfig string) types.Backdrop {
	var mapType map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlConfig), &mapType)
	assert.Nil(t, err)
	decoder := NewDecoder("example")
	config, err := decoder.DecodeBackdrop("example", mapType)
	assert.Nil(t, err)
	return config
}
