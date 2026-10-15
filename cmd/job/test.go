package job

import (
	"context"

	"dagger.io/dagger"
	"github.com/spf13/pflag"

	"github.com/osixia/dagger-build/helpers"
	"github.com/osixia/dagger-build/job"
)

type TestFlags struct {
	ExportFlags
}

var testCmdFlags = &TestFlags{}

// Resolve the job when the command runs so main can install a custom test job.
var testCmd = helpers.NewJobCmd("test", "Run build, export and test jobs", "t", func(ctx context.Context, client *dagger.Client, options any) ([]*job.BuildResult, error) {
	return job.Test(ctx, client, options)
}, testCmdFlags)

func init() {
	// flags
	testCmd.Flags().SortFlags = false
	AddTestFlags(testCmd.Flags(), testCmdFlags)
}

func AddTestFlags(fs *pflag.FlagSet, tf *TestFlags) {
	AddExportFlags(fs, &tf.ExportFlags)
}

func (tf *TestFlags) ToJobOptions() any {

	return &job.TestOptions{
		ExportOptions: *tf.ExportFlags.ToJobOptions().(*job.ExportOptions),
	}
}
