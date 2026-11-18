# CLAUDE.md — go-zendesk

A fork of [`nukosuke/go-zendesk`](https://github.com/nukosuke/go-zendesk): a Go
client library for the Zendesk API, published as
`github.com/nukosuke/go-zendesk`. It is a **library**, so there is nothing to
deploy — a release is a tag.

## Layout

| Path | What |
|---|---|
| `zendesk/` | the client: `client.go`, one file per resource, `*_types.go` for the models |
| `zendesk/mock/` | generated with `go generate ./...` (uber-go/mock) — regenerate, never hand-edit |
| `zendesk/*_test.go` | table tests against `httptest` servers, fixtures in the test files |
| `fixture/`, `script/` | request/response fixtures and helper scripts |
| `appveyor.yml` | legacy Windows CI; GitHub Actions is the one that runs |
| `CBPMigration.md` | the offset→cursor pagination guide, still the reference for that area |

## Build and test

CI (`.github/workflows/ci.yml`, `ubuntu-latest`) runs a matrix over Go 1.22 … 1.27:
`go mod download`, then `go test -v -coverprofile=profile.cov ./...`, then uploads
coverage to Coveralls. All matrix legs plus `finalize` are required checks on
`master`.

```sh
go build ./... && go vet ./... && go test ./...
go generate ./...      # only when the mock or a generated file must change
```

Coverage ignores `zendesk/mock/client.go` and `zendesk/*_generated.go` — do not
chase coverage in generated code.

## Releasing

Tag `master` with a semantic version; `go get` picks it up. Nothing is published
anywhere else. Do not rewrite a tag that has been published: someone's `go.sum`
depends on it.

## Conventions

- `master` is protected: branch → PR → every matrix leg green → merge. The last
  merged work is `fix/panic-on-unconfigured-client` (pushed, not yet merged).
- Public API stability matters more than internal tidiness: this fork is a
  drop-in replacement, so signatures stay compatible with upstream unless the
  change is the point.
- Errors are returned, never panicked, and a client that has not been configured
  returns an error rather than dereferencing nil — that is what the recent fix
  branch is about.
- New resource: a `*_types.go` file, one method per endpoint, a mock entry
  (`go generate ./...`), and a table test that asserts both the request and the
  parsed response.
