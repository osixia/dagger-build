package test

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// Version checks that the image entrypoint reports the expected image and tag.
func Version(ctx context.Context, image *dagger.Container, expected string) error {
	version, err := image.WithExec([]string{"--version"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).Stdout(ctx)
	if err != nil {
		return err
	}
	if strings.TrimSuffix(version, "\n") != expected {
		return fmt.Errorf("error: image version is %v expected image version is %v", version, expected)
	}
	return nil
}
