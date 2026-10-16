package job

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"

	"github.com/osixia/dagger-build/config"
)

// OCI archives let each platform be built and tested natively on its own runner,
// then published together as a single multi-platform image: one archive per
// image and platform, named so the publishing run finds them without a manifest.

// OCIArchiveName returns the archive file name of an image built for a platform.
// It does not depend on the version or the latest flag, so the building and the
// publishing runs agree on it even if they compute tags at different times.
func OCIArchiveName(image *config.Image, options *BuildImageOptions, platform *config.Platform) string {

	variant := "default"
	if len(image.TagPrefixes) > 0 {
		variant = image.TagPrefixes[0]
	}

	name := strings.NewReplacer("/", "-", ":", "-").Replace(image.Name)

	return fmt.Sprintf("%v_%v%v_%v.tar", name, variant, options.TagsSuffix, platform.GoArch)
}

// ValidateOCIArchiveNames checks that no two images share an archive name.
func ValidateOCIArchiveNames(images []*config.Image, options *BuildImageOptions) error {

	names := map[string]*config.Image{}

	for _, image := range images {
		for _, platform := range options.Platforms {
			name := OCIArchiveName(image, options, platform)
			if other, found := names[name]; found && other != image {
				return fmt.Errorf("images %v and %v have the same OCI archive name %v, give them different tag prefixes", other.Name, image.Name, name)
			}
			names[name] = image
		}
	}

	return nil
}

// ExportOCIArchives writes every built container to dir as an OCI archive.
func ExportOCIArchives(ctx context.Context, builds []*BuildResult, options *BuildImageOptions, dir string) error {

	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	for _, b := range builds {
		for i, c := range b.Containers {

			path := filepath.Join(dir, OCIArchiveName(b.Image, options, options.Platforms[i]))

			fmt.Printf("exporting OCI archive %v ...\n", path)
			if _, err := c.AsTarball().Export(ctx, path); err != nil {
				return fmt.Errorf("exporting OCI archive %v: %w", path, err)
			}
		}
	}

	return nil
}

// ImportImage loads the OCI archives of an image, one per platform, instead of building it.
// Archives already carry the labels set at build time.
func ImportImage(ctx context.Context, client *dagger.Client, options *BuildImageOptions, image *config.Image, dir string) (*BuildResult, error) {

	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	b := &BuildResult{
		Image:      image,
		Containers: make([]*dagger.Container, 0, len(options.Platforms)),
		Tags:       BuildTags(image, options, options.TagsSuffix),
	}

	for _, platform := range options.Platforms {

		path := filepath.Join(dir, OCIArchiveName(image, options, platform))
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("OCI archive of %v for %v: %w", image.Name, platform.Name, err)
		}

		fmt.Printf("importing OCI archive %v ...\n", path)
		c, err := client.Container().Import(client.Host().File(path)).Sync(ctx)
		if err != nil {
			return nil, fmt.Errorf("importing OCI archive %v: %w", path, err)
		}

		p, err := c.Platform(ctx)
		if err != nil {
			return nil, err
		}
		if p != platform.Name {
			return nil, fmt.Errorf("OCI archive %v contains a %v image, expected %v", path, p, platform.Name)
		}

		b.Containers = append(b.Containers, c)
	}

	return b, nil
}
