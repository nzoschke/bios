# bios

bios is a utility for running a `bios.sh` script on your laptop. It uses the [mixable/bios Docker image](https://hub.docker.com/r/mixable/bios/) to run in an environment very similar to the [BIOS GitHub App](https://www.mixable.net/docs/bios/).


```console
$ go get -u github.com/nzoschke/bios

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

  -bref string
        Base branch. Default master.
  -bsha string
        Base SHA. Default SHA of master/HEAD.
  -ref string
        Branch to check out. Default current branch.
  -sha string
        SHA to reset to. Default SHA of current branch.

  -url string
        Canonical repo URL. Default https://github.com/owner/repo.git remote.

## Check out a repo with `bios.sh` and give it a try!
$ git clone https://github.com/nzoschke/bios.git && cd bios
$ bios
DIR:  /tmp/bios
USER: <disabled>
PASS: <disabled>
BREF: master
BSHA: 203e2f2aec1e3c46a38219ba9075860a5f99a031
REF:  master
SHA:  203e2f2aec1e3c46a38219ba9075860a5f99a031
URL:  https://github.com/nzoschke/bios.git

000:0 $ run -s Cloning   git clone file:///tmp/repo/.git --branch master --single-branch src/github.com/nzoschke/bios
000:4 $ run -s Resetting git reset --hard 203e2f2aec1e3c46a38219ba9075860a5f99a031
000:4 $ run -s Fetching  git fetch origin master
000:6 $ run -s Linting   golint -set_exit_status github.com/nzoschke/bios
000:6 $ run -s Vetting   go vet -x github.com/nzoschke/bios
000:7 $ run -s Building  go build -v github.com/nzoschke/bios
001:0 $ run -s Testing   go test -v github.com/nzoschke/bios

# Results

Succeeded in 1.4 seconds. 🆗

## Statuses

```diff
+ Cloning
+ Resetting
+ Fetching
+ Whitespacing
+ Linting
+ Vetting
+ Building
+ Testing
```

Succeeded (Testing)
```

## GitHub Interactions

By default `bios` will not interact with the GitHub API to set commit status or comments. Use the `-hub` flag to change this behavior.

Be careful with this setting...

Enabling `-hub` without a `-user` or `-pass` will use the `git credential` helper to locate credentials. These credentials likely have more access than the BIOS GitHub App.

This will also cause an error if the local SHA doesn't exist on GitHub.
