package image

import (
	"encoding/json"
	"io"
	"net"
	"os"

	"github.com/containerd/console"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stringid"
	cfgtypes "github.com/dodo/dodo-build/pkg/types"
	"github.com/dodo/dodo-core/pkg/decoder"
	log "github.com/hashicorp/go-hclog"
	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/oclaussen/go-gimme/configfiles"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

func (image *Image) Get() (string, error) {
	if image.config.ForceRebuild || len(image.config.ImageName) == 0 {
		return image.Build()
	}

	imgs, err := image.client.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: filters.NewArgs(filters.Arg("reference", image.config.ImageName)),
		},
	)
	if err != nil || len(imgs) == 0 {
		return image.Build()
	}

	return imgs[0].ID, nil
}

func (image *Image) Build() (string, error) {
	// TODO: refactor here, the dependency on config is uncomfortable
	backdrops := map[string]*cfgtypes.Backdrop{}
	_, err := configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			d := decoder.New(configFile.Path)
			d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
				"backdrops": decoder.Map(cfgtypes.NewBackdrop(), &backdrops),
			})
			return false
		},
	})

	if err != nil {
		log.L().Error("error finding config files", "error", err)
	}

	for _, name := range image.config.Dependencies {
		for _, backdrop := range backdrops {
			if backdrop.Build != nil && backdrop.Build.ImageName == name {
				if image.config.ForceRebuild {
					backdrop.Build.ForceRebuild = true
				}

				dependency, err := NewImage(image.client, image.authConfigs, backdrop.Build)
				if err != nil {
					return "", err
				}

				if _, err := dependency.Get(); err != nil {
					return "", err
				}
			}
		}
	}

	contextData, err := prepareContext(image.config, image.session)
	if err != nil {
		return "", err
	}

	defer contextData.cleanup()

	imageID := ""
	displayCh := make(chan *client.SolveStatus)

	eg, _ := errgroup.WithContext(appcontext.Context())

	eg.Go(func() error {
		return image.session.Run(
			context.TODO(),
			func(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
				return image.client.DialHijack(ctx, "/session", proto, meta)
			},
		)
	})

	if image.config.ForceRebuild {
		eg.Go(func() error {
			cons, err := console.ConsoleFromFile(os.Stderr)
			if err != nil {
				return err
			}

			return progressui.DisplaySolveStatus(context.TODO(), "", cons, os.Stderr, displayCh)
		})
	}

	eg.Go(func() error {
		defer func() {
			close(displayCh)
			image.session.Close()
		}()

		imageID, err = image.runBuild(contextData, displayCh)

		return err
	})

	err = eg.Wait()
	if err != nil {
		return "", err
	}

	if imageID == "" {
		return "", ErrMissingImageID
	}

	return imageID, nil
}

func (image *Image) runBuild(contextData *contextData, displayCh chan *client.SolveStatus) (string, error) {
	args := map[string]*string{}
	for _, arg := range image.config.Arguments {
		args[arg.Key] = &arg.Value
	}

	var tags []string
	if image.config.ImageName != "" {
		tags = append(tags, image.config.ImageName)
	}

	response, err := image.client.ImageBuild(
		context.Background(),
		nil,
		types.ImageBuildOptions{
			Tags:           tags,
			SuppressOutput: false,
			NoCache:        image.config.NoCache,
			Remove:         true,
			ForceRemove:    true,
			PullParent:     image.config.ForcePull,
			Dockerfile:     contextData.dockerfileName,
			BuildArgs:      args,
			AuthConfigs:    image.authConfigs,
			Version:        types.BuilderBuildKit,
			RemoteContext:  contextData.remote,
			SessionID:      image.session.ID(),
			BuildID:        stringid.GenerateRandomID(),
		},
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return handleBuildResult(response.Body, displayCh, image.config.ForceRebuild)
}

func handleBuildResult(response io.Reader, displayCh chan *client.SolveStatus, printOutput bool) (string, error) {
	var imageID string

	decoder := json.NewDecoder(response)

	for {
		var msg jsonmessage.JSONMessage
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				return imageID, nil
			}

			return "", err
		}

		if msg.Error != nil {
			return "", msg.Error
		}

		if msg.Aux != nil {
			if msg.ID == "moby.image.id" {
				var result types.BuildResult
				if err := json.Unmarshal(*msg.Aux, &result); err != nil {
					continue
				}

				imageID = result.ID
			} else if printOutput && msg.ID == "moby.buildkit.trace" {
				var resp controlapi.StatusResponse
				var dt []byte

				if err := json.Unmarshal(*msg.Aux, &dt); err != nil {
					continue
				}

				if err := (&resp).Unmarshal(dt); err != nil {
					continue
				}

				s := client.SolveStatus{}
				for _, v := range resp.Vertexes {
					s.Vertexes = append(s.Vertexes, &client.Vertex{
						Digest:    v.Digest,
						Inputs:    v.Inputs,
						Name:      v.Name,
						Started:   v.Started,
						Completed: v.Completed,
						Error:     v.Error,
						Cached:    v.Cached,
					})
				}
				for _, v := range resp.Statuses {
					s.Statuses = append(s.Statuses, &client.VertexStatus{
						ID:        v.ID,
						Vertex:    v.Vertex,
						Name:      v.Name,
						Total:     v.Total,
						Current:   v.Current,
						Timestamp: v.Timestamp,
						Started:   v.Started,
						Completed: v.Completed,
					})
				}
				for _, v := range resp.Logs {
					s.Logs = append(s.Logs, &client.VertexLog{
						Vertex:    v.Vertex,
						Stream:    int(v.Stream),
						Data:      v.Msg,
						Timestamp: v.Timestamp,
					})
				}

				displayCh <- &s
			}
		}
	}
}
