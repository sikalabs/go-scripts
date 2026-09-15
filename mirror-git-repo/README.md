# mirror-git-repo

Mirror a git repo from one URL to another, using the `git` binary
(`git clone --mirror` + `git push --mirror`).

## Install (master)

```
go install github.com/sikalabs/go-scripts/mirror-git-repo@master
```

## Example Usage

```bash
mirror-git-repo \
  --source-url https://gitlab-source.sikademo.com/source-group/source-repo.git \
  --target-url https://gitlab-target.corp.com/target-group/target-repo.git \
  --source-token glpat-source-token \
  --target-token glpat-target-token
```

`--source-token` / `--target-token` are optional and, when set, are injected
into the corresponding HTTPS URL as credentials. They have no effect on
non-HTTP(S) URLs (e.g. `git@host:path.git`) — use SSH keys/agent for those.
