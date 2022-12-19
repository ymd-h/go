# go work tidy

This is the POC implementation of `go mod tidy` for multi-module
workspace.

## Usage
```shell
go run github.com/ymd-h/go/worktidy
```


## Expected Result
* `require`s for non-workspace module at all `go.mod`s are properly
  updated as `go mod tidy` do.
* `require`s for workspace module at all `go.mod`s are updated with
  latest version tag



## Limitation
`go.work` file must be placed at the repository root directory, and
the command should be exeuted at the directory.
