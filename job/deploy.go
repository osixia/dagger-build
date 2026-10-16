package job

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

type DeployOptions struct {
	TestOptions

	Registry string
	Username string
	Password string

	DigestsFile string

	DryRun bool
}

func DefaultDeploy(ctx context.Context, client *dagger.Client, options any) ([]*BuildResult, error) {

	o := options.(*DeployOptions)

	// imported images were tested natively by the runs that built them,
	// testing them again would need emulation for other platforms
	test := Test
	if o.OCIImportFrom != "" {
		test = func(ctx context.Context, client *dagger.Client, options any) ([]*BuildResult, error) {
			return Export(ctx, client, &options.(*TestOptions).ExportOptions)
		}
	}

	builds, err := test(ctx, client, &o.TestOptions)
	if err != nil {
		return nil, err
	}

	// all tags of an image share its digest: record it once, so it is signed once
	digests := map[string]bool{}

	for _, b := range builds {

		publishOptions := dagger.ContainerPublishOpts{
			PlatformVariants: b.Containers,
			// Some registries may require explicit use of docker mediatypes
			// rather than the default OCI mediatypes
			// MediaTypes: dagger.Dockermediatypes,
		}

		for _, tag := range b.Tags {

			img := fmt.Sprintf("%v:%v", b.Image.Name, tag)

			c := client.Container()

			if o.Registry != "" {

				if o.Registry != "" {
					img = fmt.Sprintf("%v/%v", o.Registry, img)
				}

				if o.Username != "" {
					c = c.WithRegistryAuth(o.Registry, o.Username, client.SetSecret("registry-password", o.Password))
				}

			}

			fmt.Printf("pushing image %v ...\n", img)

			if o.DryRun {
				continue
			}

			digest, err := c.Publish(ctx, img, publishOptions)
			if err != nil {
				return nil, err
			}
			fmt.Printf("image %v pushed with digest %v\n", img, digest)

			_, sha, _ := strings.Cut(digest, "@")
			if o.DigestsFile != "" && !digests[sha] {
				digests[sha] = true

				f, err := os.OpenFile(o.DigestsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return nil, fmt.Errorf("opening cosign output file: %w", err)
				}
				_, writeErr := fmt.Fprintf(f, "%v\n", digest)
				if closeErr := f.Close(); closeErr != nil && writeErr == nil {
					writeErr = closeErr
				}
				if writeErr != nil {
					return nil, fmt.Errorf("writing to cosign output file: %w", writeErr)
				}
			}
		}
	}

	return builds, nil
}
