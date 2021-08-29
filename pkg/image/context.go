package image

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/urlutil"
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	log "github.com/hashicorp/go-hclog"
	buildkit "github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/filesync"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/pkg/errors"
	fstypes "github.com/tonistiigi/fsutil/types"
)

const clientSession = "client-session"

type contextData struct {
	remote         string
	dockerfileName string
	contextDir     string
}

func (data *contextData) tempdir() (string, error) {
	if len(data.contextDir) == 0 {
		dir, err := ioutil.TempDir("", "dodo-temp-")
		if err != nil {
			return "", fmt.Errorf("could not create temporary directory: %w", err)
		}

		data.contextDir = dir
	}

	return data.contextDir, nil
}

func (data *contextData) cleanup() {
	if data.contextDir != "" {
		os.RemoveAll(data.contextDir)
	}
}

func prepareContext(config *api.BuildInfo, session session) (*contextData, error) {
	log.L().Debug("preparing context")

	data := contextData{
		remote:         clientSession,
		dockerfileName: config.Dockerfile,
	}
	syncedDirs := []filesync.SyncedDir{}

	if config.Context == "" {
		dir, err := data.tempdir()
		if err != nil {
			data.cleanup()

			return nil, err
		}

		syncedDirs = append(syncedDirs, filesync.SyncedDir{Name: "context", Dir: dir})
	} else if _, err := os.Stat(config.Context); err == nil {
		syncedDirs = append(syncedDirs, filesync.SyncedDir{
			Name: "context",
			Dir:  config.Context,
			Map: func(_ string, stat *fstypes.Stat) bool {
				stat.Uid = 0
				stat.Gid = 0

				return true
			},
		})
	} else if urlutil.IsURL(config.Context) {
		data.remote = config.Context
	} else {
		return nil, errors.Errorf("Context directory does not exist: %v", config.Context)
	}

	if len(config.InlineDockerfile) > 0 {
		steps := ""
		for _, step := range config.InlineDockerfile {
			steps = steps + step + "\n"
		}

		dir, err := data.tempdir()
		if err != nil {
			data.cleanup()

			return nil, err
		}

		tempfile := filepath.Join(dir, "Dockerfile")
		if err := writeDockerfile(tempfile, steps); err != nil {
			data.cleanup()

			return nil, err
		}

		data.dockerfileName = filepath.Base(tempfile)
		dockerfileDir := filepath.Dir(tempfile)

		syncedDirs = append(syncedDirs, filesync.SyncedDir{
			Name: "dockerfile",
			Dir:  dockerfileDir,
		})
	} else if config.Dockerfile != "" && data.remote == clientSession {
		data.dockerfileName = filepath.Base(config.Dockerfile)
		dockerfileDir := filepath.Dir(config.Dockerfile)

		syncedDirs = append(syncedDirs, filesync.SyncedDir{
			Name: "dockerfile",
			Dir:  dockerfileDir,
		})
	}

	log.L().Debug(
		"prepared context",
		"remote", data.remote,
		"dockerfileName", data.dockerfileName,
		"contextDir", data.contextDir,
		"config", config,
	)

	if len(syncedDirs) > 0 {
		session.Allow(filesync.NewFSSyncProvider(syncedDirs))
		log.L().Debug("added context directories", "dirs", syncedDirs)
	}

	session.Allow(authprovider.NewDockerAuthProvider(ioutil.Discard))

	if len(config.Secrets) > 0 {
		provider, err := secretsProvider(config)
		if err != nil {
			return nil, err
		}

		session.Allow(provider)
	}

	if len(config.SshAgents) > 0 {
		provider, err := sshAgentProvider(config)
		if err != nil {
			return nil, err
		}

		session.Allow(provider)
	}

	return &data, nil
}

func writeDockerfile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", path, err)
	}
	defer file.Close()

	rc := ioutil.NopCloser(bytes.NewReader([]byte(content)))

	if _, err := io.Copy(file, rc); err != nil {
		return fmt.Errorf("could not write dockerfile: %w", err)
	}

	if err := rc.Close(); err != nil {
		return fmt.Errorf("coud not close dockerfile stream: %w", err)
	}

	return nil
}

func secretsProvider(config *api.BuildInfo) (buildkit.Attachable, error) {
	sources := make([]secretsprovider.Source, 0, len(config.Secrets))

	for _, secret := range config.Secrets {
		source := secretsprovider.Source{
			ID:       secret.Id,
			FilePath: secret.Path,
		}
		sources = append(sources, source)
	}

	store, err := secretsprovider.NewStore(sources)
	if err != nil {
		return nil, fmt.Errorf("could not initialezie secrets store: %w", err)
	}

	return secretsprovider.NewSecretProvider(store), nil
}

func sshAgentProvider(config *api.BuildInfo) (buildkit.Attachable, error) {
	configs := make([]sshprovider.AgentConfig, 0, len(config.SshAgents))

	for _, agent := range config.SshAgents {
		config := sshprovider.AgentConfig{
			ID:    agent.Id,
			Paths: []string{agent.IdentityFile},
		}
		configs = append(configs, config)
	}

	return sshprovider.NewSSHAgentProvider(configs)
}
