package internal

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const dockerImage = "golang:1.26.0-alpine3.23"

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
			`cd /code && go mod init gorun && go mod tidy && go build -o app app.go && ./app`,
		},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", "", err
	}
	defer func() {
		cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCleanup()
		if err := cli.ContainerRemove(cleanupCtx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			log.Printf("container cleanup failed for %s: %v", resp.ID, err)
		}
	}()

	source, err := sourceArchive(code)
	if err != nil {
		return "", "", err
	}
	if err := cli.CopyToContainer(ctx, resp.ID, "/", source, container.CopyToContainerOptions{}); err != nil {
		return "", "", err
	}

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
	defer out.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, out)
	if err != nil {
		return "", "", err
	}

	stderr := cleanStderr(stderrBuf.String())

	return strings.TrimSpace(stdoutBuf.String()), strings.TrimSpace(stderr), nil
}

func sourceArchive(code string) (io.Reader, error) {
	var archive bytes.Buffer
	w := tar.NewWriter(&archive)
	if err := w.WriteHeader(&tar.Header{Name: "code/", Mode: 0755, Typeflag: tar.TypeDir}); err != nil {
		return nil, err
	}
	if err := w.WriteHeader(&tar.Header{Name: "code/app.go", Mode: 0600, Size: int64(len(code))}); err != nil {
		return nil, err
	}
	if _, err := io.WriteString(w, code); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return &archive, nil
}

var excludePrefixes = []string{
	"go: finding module ",
	"go: downloading ",
	"go: found ",
	"go: creating",
}

func cleanStderr(stderr string) string {
	var cleanedLines []string
	for line := range strings.SplitSeq(stderr, "\n") {
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
