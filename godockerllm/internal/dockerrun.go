package internal

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const dockerImage = "golang:1.25.1-alpine3.22"

func RunCodeInDocker(code string, timeout time.Duration) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", "", err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/"+dockerImage, image.PullOptions{})
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return "", "", err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: dockerImage,
		Cmd: []string{
			"/bin/sh",
			"-c",
			`mkdir -p /code && cd /code && go mod init gorun && echo '` + code + `' > app.go && go mod tidy && go build -o app app.go && ./app`,
		},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", "", err
	}
	defer func() {
		_ = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	}()

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", "", err
		}
	case <-statusCh:
	case <-ctx.Done():
		return "", "", ctx.Err()
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", "", err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, out)
	if err != nil {
		return "", "", err
	}

	stderr := cleanStderr(stderrBuf.String())

	return strings.TrimSpace(stdoutBuf.String()), strings.TrimSpace(stderr), nil
}

var excludePrefixes = []string{
	"go: finding module ",
	"go: downloading ",
	"go: found ",
	"go: creating",
}

func cleanStderr(stderr string) string {
	var cleanedLines []string
	for _, line := range strings.Split(stderr, "\n") {
		exclude := false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(line, prefix) {
				exclude = true
				break
			}
		}
		if !exclude {
			cleanedLines = append(cleanedLines, line)
		}
	}
	return strings.Join(cleanedLines, "\n")
}
