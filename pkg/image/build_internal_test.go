package image

import (
	"testing"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/moby/buildkit/client"
	"github.com/stretchr/testify/assert"
)

func TestBuildImage(t *testing.T) {
	displayCh := make(chan *client.SolveStatus)
	defer close(displayCh)

	image := fakeImage(t, &api.BuildInfo{
		Context: "./test",
	})
	result, err := image.runBuild(&contextData{
		remote:         "client-session",
		dockerfileName: "Dockerfile",
	}, displayCh)
	assert.Nil(t, err)
	assert.Equal(t, "NewImageID", result)
}

func TestBuildInlineImage(t *testing.T) {
	displayCh := make(chan *client.SolveStatus)
	defer close(displayCh)

	image := fakeImage(t, &api.BuildInfo{
		InlineDockerfile: []string{"FROM scratch"},
	})
	result, err := image.runBuild(&contextData{
		remote:         "client-session",
		dockerfileName: "Dockerfile",
	}, displayCh)
	assert.Nil(t, err)
	assert.Equal(t, "NewImageID", result)
}
