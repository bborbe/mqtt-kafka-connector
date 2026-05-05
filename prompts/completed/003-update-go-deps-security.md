---
status: completed
container: mqtt-kafka-connector-003-update-go-deps-security
dark-factory-version: v0.148.4-3-gc45254a
created: "2026-05-05T17:38:43Z"
queued: "2026-05-05T17:38:43Z"
started: "2026-05-05T17:39:31Z"
completed: "2026-05-05T17:49:54Z"
---

<summary>
- Go dependencies updated to latest allowed versions
- `github.com/docker/docker` bumped to v29.3.1 or newer (resolves Dependabot advisory)
- `github.com/go-git/go-git/v5` confirmed at >= v5.18.0 (CVE-2026-41506)
- make precommit passes cleanly
- `## Unreleased` section in CHANGELOG.md lists the bumped versions
</summary>

<objective>
Update Go module dependencies to resolve Dependabot security advisories on `docker/docker` and verify `go-git/v5` is at the patched version.
</objective>

<context>
Read CLAUDE.md for project conventions.
Read `docs/dod.md` for the Definition of Done criteria.

Current state of `go.mod` (both deps are `// indirect`):
- `github.com/docker/docker v28.5.2+incompatible` — vulnerable, advisory: bump to >= v29.3.1
- `github.com/go-git/go-git/v5 v5.18.0` — already at fixed version (CVE-2026-41506); verify only

`updater` is pre-installed in the claude-yolo container.
</context>

<requirements>
1. Run `updater --verbose --yes go` in the **foreground** (do NOT background this command).
2. If `updater` fails on any rename, follow recovery: `grep -r '<stale-identifier>' --exclude-dir=vendor`, fix all occurrences, re-run `make generate`, `make test`. Common rename patterns from prior runs: `*Id` → `*ID`, `*Url` → `*URL`, `HttpClient` → `HTTPClient`.
3. Check `go.mod` for `github.com/docker/docker`. If it is still `< v29.3.1`, run `go get github.com/docker/docker@latest && go mod tidy`. The `+incompatible` suffix in the version string is expected and not an error. The dep may transition from `// indirect` to direct — that is acceptable.
4. Verify `go.mod` shows `github.com/go-git/go-git/v5 >= v5.18.0`. If a regression dropped it below, run `go get github.com/go-git/go-git/v5@latest && go mod tidy`.
5. Run `make precommit` — must pass with exit code 0.
6. Update `CHANGELOG.md`:
   - If a `## Unreleased` heading does not exist, insert one above the most recent version section.
   - Under `## Unreleased`, add one bullet per bumped dep (format: `- Bump github.com/docker/docker to vX.Y.Z (Dependabot advisory)`).
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT run `updater` as a background task — use foreground with `--verbose`
- Existing tests must still pass
- No `exclude` or `replace` directives in go.mod
- Do NOT hand-edit version numbers in `go.mod` — let `updater` / `go get` write them
</constraints>

<verification>
Run `make precommit` — must pass with exit code 0.
Run `go list -m github.com/docker/docker` — version must be >= v29.3.1 (the `+incompatible` suffix is expected).
Run `go list -m github.com/go-git/go-git/v5` — version must be >= v5.18.0.
</verification>
