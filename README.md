# CoPR

Updates branches/PRs in accompanying repositories with output generated from branches/PRs in a main repository.

Example uses:
- For GitOps, manage shadow PRs in the ops repository generated from a main repository.
- For docs, manage shadow PRs in a wiki repository of docs generated from code in other repositories.

At the moment, this only creates/updates branches. In the future this should become more well-integrated with review
flows, e.g. by setting PR statuses.

## Example

We'll use this CoPR configuration for the example:

```yml
outputs:
- repository: github.com/my-org/docs
  generate: ./scripts/gen-docs.sh
  directory: ./dist/docs
- repository: github.com/my-org/config
  generate: ./config/gen.sh
  directory: ./config/output
```

This would be checked into a repository github.com/my-org/product as e.g. `copr.yaml`. We assume the same repository also has:
- a script `scripts/gen-docs.sh` which generates all documentation into the directory `dist/docs`
- a `config/gen.sh` that uses config templates to generate configuration in the directory `config/output`

Running `copr` in this repository will pull the docs and wiki repositories and create branches in each of them. The branch
names will match the current branch name in the source repository. The destination branches will be updated with the
output from running the `generate` command in the source repository.
