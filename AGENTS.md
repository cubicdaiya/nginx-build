# Repository Guidelines

## Project Structure & Module Organization
CLI entrypoints sit in `nginx-build.go` and the `command/` package. Build orchestration and artifact helpers live in `builder/`, `configure/`, and `util/`. Third-party integration logic is grouped under `module3rd/` and `openresty/`. Example configs that drive CI and local builds reside in `config/`. Tests are colocated with the code they cover (for example `download_test.go`, `configure/configure_test.go`). Generated assets, such as demo GIFs, stay in `images/`, while temporary build outputs default to `work/` and should be git-ignored.

## Build, Test, and Development Commands
Use `make nginx-build` to compile the binary with the current commit hash embedded. `./nginx-build -d work` performs a minimal fetch-and-build cycle in a sandbox directory. Run `make build-example` to validate the sample module + configure combo stored in `config/`. Execute `make check` or `go test ./...` before every PR; set `GOCACHE=$(pwd)/.cache` if the environment restricts the default Go cache. Apply formatting with `make fmt`, and reset local artifacts via `make clean`.

## Coding Style & Naming Conventions
All Go code must be `gofmt`-clean (tabs for indentation). Favor descriptive package names (`builder`, `module3rd`) and CamelCase identifiers; export only what consumers need. Command-line flags follow dash-case (`-customssl`, `-workdir`) and should be declared in `option.go` with matching help text. Keep module paths stable so CI workflows that rely on directory globs remain accurate.

## Testing Guidelines
Prefer table-driven Go tests living next to their targets (`*_test.go`). Name tests after the feature under exercise (`TestGeneratePreservesTrailingComments`). Run `go test ./...` for unit coverage and rerun targeted integ tests such as `./nginx-build -d work -clear -openssl` when touching dependency download logic. Document long-running or flaky scenarios directly in the PR body.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`fix: handle trailing comments`, `chore(ci): tighten go workflow paths`). Combine related edits; avoid mixing refactors with feature changes. Every PR should describe the behavior change, list verification commands, and link issues when applicable. If CI is intentionally skipped, justify it explicitly. Attach logs or screenshots for build output changes that reviewers cannot reproduce quickly.

## CI & Security Notes
GitHub Actions only trigger when relevant files change; confirm your edits touch the directories listed in `.github/workflows/`. For security-sensitive updates, run `make check` plus `go test ./...` and monitor `govulncheck` results. Do not commit credentials or generated binaries (`nginx-build`).
