package job

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	jobtest "github.com/osixia/dagger-build/job/test"
)

// Test options
// =============================

type TestOptions struct {
	ExportOptions
}

func DefaultTest(ctx context.Context, client *dagger.Client, options any) ([]*BuildResult, error) {

	o := options.(*TestOptions)

	builds, err := Export(ctx, client, &o.ExportOptions)
	if err != nil {
		return nil, err
	}

	for _, b := range builds {

		expectedImageVersion := fmt.Sprintf("%v:%v", b.Image.Name, b.LongestTag())

		for _, i := range b.Containers {

			if err := jobtest.Version(ctx, i, expectedImageVersion); err != nil {
				return nil, err
			}

		}
	}

	return builds, nil
}
