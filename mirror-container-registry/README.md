# mirror-container-registry

## Install (master)

```
go install github.com/sikalabs/go-scripts/mirror-container-registry@master
```

## Example Usage

```bash
mirror-container-registry \
  --source gitlab-registry-source.sikademo.com \
  --target gitlab-registry-target.corp.com \
  --source-group source-group \
  --target-group target-group/target-subgroup \
  --source-token glpat-source-token \
  --target-token glpat-target-token
```
