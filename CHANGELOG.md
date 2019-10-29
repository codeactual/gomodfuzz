# Change Log

## v0.1.6

- fix: only display pass summary in --verbose

## v0.1.5

- refactor: move `main` package to `./cmd/gomodfuzz`
- fix: add missing `-mod=vendor` to `make build`
- dep: update first-party dependencies under `internal`

## v0.1.4

- fix: stage deletion of module cache directories using this approach https://github.com/golang/go/blob/go1.12.2/src/cmd/go/internal/modfetch/unzip.go#L161
- dep: update first-party dependencies under `internal`

## v0.1.3

- fix: reenabled stage deletion now that https://github.com/golang/go/issues/30579 landed

## v0.1.2

- fix: unreliable stage deletion

## v0.1.1

- feat: display standard error/output lengths
- dep: update first-party dependencies under `internal`
- docs: fix usage example typos

## v0.1.0

- feat: initial project export from private monorepo
