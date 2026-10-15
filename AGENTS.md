# Agent instructions

This repository contains the shared Go library and Cobra CLI used by Osixia
image projects. It uses the standalone Dagger Go SDK, not a native Dagger module.

- Never commit, push, tag, publish images, or create/update GitHub releases
  without explicit user approval. Modification requests do not authorize these.
- Preserve the existing API and project hooks. Keep image matrices, custom build
  arguments, runtime fixtures, and example tests in each consuming project.
- Read the relevant code first, preserve existing changes, and keep diffs focused.
- Keep code, comments, CLI help, and documentation in English. Run gofmt.
- Go 1.26.1 or newer is required by Dagger 0.21.10. Keep the SDK and engine
  versions aligned. The SDK downloads the matching CLI when it connects.
- Validate with `go test -race ./...`, `go vet ./...`, `go build ./...`, and
  `golangci-lint run` when available.
- Run image tests through a consuming project's `job test` command in disposable
  containers. Never use deploy/publish or the GitHub tag pipeline as validation:
  they can push images and modify releases and assets, including in dry-run paths
  that still make authenticated API reads.
- The generic test job checks entrypoint `--version`; it must not read a
  consuming project's docs, service scripts, or runtime fixtures.
