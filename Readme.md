# bios

bios is a utility for running a `bios.sh` script on your laptop. It uses the [mixable/bios Docker image](https://hub.docker.com/r/mixable/bios/) to run in an environment very similar to the [BIOS GitHub App](https://www.mixable.net/docs/bios/).


```console
$ go install github.com/nzoschke/bios

$ bios -h
usage: bios [options] [<directory>]
  -help
    	Show usage.

  -hub
    	Interact with GitHub API. Must be set for -user and -pass to have effect.
  -pass string
    	Auth password. Default 'git credential' password for -url.
  -user string
    	Auth username. Default 'git credential' username for -url.

  -ref string
    	Branch to check out. Default current branch.
  -sha string
    	SHA to reset to. Default current SHA.
  -url string
    	Canonical repo URL. Default https://github.com/owner/repo.git remote.

$ bios
DIR:  /Users/noah/go/src/github.com/nzoschke/bios
USER: <disabled>
PASS: <disabled>
REF:  local
SHA:  8a4c62642ae54da63f9bd7b7124d0e487427bcc5
URL:  https://github.com/nzoschke/bios.git

000:0 $ run -s Cloning git clone file:///tmp/repo/.git --branch local --single-branch .
000:3 $ run -s Linting golint -set_exit_status github.com/nzoschke/bios
000:3 $ run -s Vetting go vet -x github.com/nzoschke/bios
000:4 $ run -s Building go build -v github.com/nzoschke/bios
000:8 $ run -s Testing go test -v github.com/nzoschke/bios

# Results

Succeeded in 1.1 seconds. 🆗

## Statuses

```diff
+ Cloning
+ Linting
+ Vetting
+ Building
+ Testing
```

## GitHub Interactions

By default `bios` will not interact with the GitHub API to set commit status or comments. Use the `-hub` flag to change this behavior.

Be careful with this setting...

Enabling `-hub` without a `-user` or `-pass` will use the `git credential` helper to locate credentials. These credentials likely have more access than the BIOS GitHub App.

This will also cause an error if the local SHA doesn't exist on GitHub.
