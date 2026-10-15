package provider

import (
	cmdjob "github.com/osixia/dagger-build/cmd/job"
	"github.com/osixia/dagger-build/helpers"
	"github.com/osixia/dagger-build/job"
	"github.com/osixia/dagger-build/provider"
)

type githubFlags struct {
	cmdjob.PublishFlags

	githubRef   string
	githubToken string

	changelog string
}

var githubCmdFlags = &githubFlags{}

var githubCmd = helpers.NewJobCmd("github --github-ref reference", "Run jobs based on GitHub Actions CI/CD environment", "g", provider.Github, githubCmdFlags)

func init() {
	// flags
	githubCmd.Flags().SortFlags = false

	githubCmd.Flags().StringVarP(&githubCmdFlags.githubRef, "github-ref", "r", "", "github reference")
	githubCmd.Flags().StringVarP(&githubCmdFlags.githubToken, "github-token", "t", "", "github token (required for tags)\n")
	helpers.Must(githubCmd.MarkFlagRequired("github-ref"))

	githubCmd.Flags().StringVarP(&githubCmdFlags.changelog, "changelog", "c", "", "changelog\n")

	cmdjob.AddPublishFlags(githubCmd.Flags(), &githubCmdFlags.PublishFlags)
}

func (gf *githubFlags) ToJobOptions() any {

	return &provider.GithubOptions{
		PublishOptions: *gf.PublishFlags.ToJobOptions().(*job.PublishOptions),

		Ref:   gf.githubRef,
		Token: gf.githubToken,

		Changelog: gf.changelog,
	}
}
