---
status: completed
summary: Bumped go-git/v5 to v5.19.0 (CVE-2026-45022), updated Go to 1.26.3, removed 6 stale unused ignore entries from .osv-scanner.toml, and updated CHANGELOG.md with an Unreleased section.
container: mqtt-kafka-connector-004-update-go-git-and-suppress-docker
dark-factory-version: v0.156.1-1-g04f3863-dirty
created: "2026-05-12T13:00:00Z"
queued: "2026-05-12T17:14:40Z"
started: "2026-05-12T19:42:00Z"
completed: "2026-05-12T19:46:34Z"
---

<summary>
- Bumps `github.com/go-git/go-git/v5` from v5.18.0 to v5.19.0 (CVE-2026-45022 High)
- Bumps Go from 1.26.2 to 1.26.3 (stdlib CVEs reported by osv-scanner)
- Resolves Dependabot go-git advisory for bborbe/mqtt-kafka-connector
- For `github.com/docker/docker` (no upstream fix yet — latest is v28.5.2, advisory wants >= v29.3.1): keeps the existing advisory IDs in `.trivyignore` and `.osv-scanner.toml`; appends any NEW IDs that surface during scan
- Removes stale "unused ignore" entries from `.osv-scanner.toml` (osv-scanner errors on these)
- `make precommit` exits 0 after the change
- CHANGELOG `## Unreleased` documents the bumps and ignore-list cleanup
</summary>

<objective>
Patch CVE-2026-45022 in `github.com/go-git/go-git/v5` by upgrading to v5.19.0, patch stdlib CVEs via Go 1.26.3, and unblock `make precommit` for the docker/docker advisory by adding the unresolvable IDs to the project's security-scanner ignore lists.
</objective>

<context>
Read `CLAUDE.md` for project conventions.

Current `go.mod` (both indirect):
- `github.com/go-git/go-git/v5 v5.18.0` — has fix: bump to v5.19.0
- `github.com/docker/docker v28.5.2+incompatible` — no upstream fix (latest available is v28.5.2; advisory wants >= v29.3.1)

Current Go: 1.26.2 in `go.mod`. This repo has NO Dockerfile (verified 2026-05-12) — only update `go.mod`.

Existing ignore-file patterns to mirror:

`.trivyignore`:
```
# github.com/docker/docker indirect dep, no fix available via Go modules
CVE-xxxx-xxxxx
```

`.osv-scanner.toml`:
```toml
[[IgnoredVulns]]
id = "GHSA-xxxx-xxxx-xxxx"
reason = "github.com/docker/docker indirect dep, no fix available"
```
</context>

<requirements>
1. Bump go-git/v5:
   ```bash
   go get github.com/go-git/go-git/v5@v5.19.0
   go mod tidy
   ```

2. Bump Go version 1.26.2 → 1.26.3 (patches stdlib CVEs reported by osv-scanner):
   - Edit `go.mod`: change `go 1.26.2` to `go 1.26.3`
   - No Dockerfile exists in this repo — skip Dockerfile update
   - Run `go mod tidy` again

3. Remove stale "unused ignore" entries from `.osv-scanner.toml`. After running `make precommit`, osv-scanner reports which `[[IgnoredVulns]]` IDs are unused. Delete those blocks. Verify against the actual scanner output before deleting.

4. Run `make precommit`. If it fails because trivy or osv-scanner reports NEW advisory IDs (CVE-* or GHSA-*) on `github.com/docker/docker` that are NOT yet in the ignore files, add them:
   - Append CVE-IDs to `.trivyignore` under the existing `# github.com/docker/docker indirect dep, no fix available via Go modules` block (create the block if it doesn't exist).
   - Append GHSA-IDs as new `[[IgnoredVulns]]` blocks in `.osv-scanner.toml` with `reason = "github.com/docker/docker indirect dep, no fix available"`.
   - Re-run `make precommit` until it exits 0.

5. Do NOT add ignore entries for any advisory that has a fix available (those must be patched, not suppressed). Only docker/docker is allowed to be suppressed in this prompt.

6. Update `CHANGELOG.md` under `## Unreleased`:
   ```
   - security: bump github.com/go-git/go-git/v5 to v5.19.0 (CVE-2026-45022)
   - security: bump Go to 1.26.3 (stdlib CVEs)
   - chore: remove stale unused ignore entries from .osv-scanner.toml
   - security: suppress docker/docker advisories <list-IDs> in .trivyignore/.osv-scanner.toml (no upstream fix; latest is v28.5.2, advisory wants >= v29.3.1)
   ```
   Replace `<list-IDs>` with the IDs you actually added (or write "no new IDs" if the existing ignore lists already covered everything).

7. Verify:
   - `go list -m github.com/go-git/go-git/v5` reports `v5.19.0`
   - `grep '^go ' go.mod` shows `go 1.26.3`
   - `make precommit` exits 0
</requirements>

<constraints>
- Only edit: `go.mod`, `go.sum`, `.trivyignore`, `.osv-scanner.toml`, `CHANGELOG.md`
- Do NOT bump unrelated deps
- Do NOT add docker/docker as a direct dep (must stay `// indirect`)
- Do NOT add a `replace` or `exclude` directive
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
go list -m github.com/go-git/go-git/v5     # must print v5.19.0
grep '^go ' go.mod                          # must print "go 1.26.3"
make precommit                              # must exit 0
```
</verification>
