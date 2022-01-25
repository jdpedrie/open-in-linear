# Open in Linear

This tool provides a shortcut to opening a linear issue in the desktop app or
browser.

## Install

```
go install github.com/jdpedrie/open-in-linear@latest
```

## Configuration

The name of the linear workspace is required. It can be provided in one of three
ways.

1. The CLI accepts a `-workspace` flag.
2. The current git repo can provide a `linear.workspace` config value. This
   setting only functions when executed from a git repository, or when the
   `-repo` flag is provided.
3. Set the `LINEAR_WORKSPACE` environment variable.

These three locations are checked in order, and the first to be found is
preferred.

### Setting the git config value

```sh
$ git config --add linear.workspace <my-workspace>
```

## Usage

When run from a git repository where the current branch name contains a linear
issue name:

```sh
$ open-in-linear -workspace <my-workspace>
```

Or to use a git repo in a different directory:

```sh
$ open-in-linear -workspace <my-workspace>
```

Or to open any issue by its name:

$ open-in-linear
```sh
$ open-in-linear -workspace <my-workspace> -issue DEV-1
```
