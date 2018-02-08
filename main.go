package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var reGH = regexp.MustCompile("https://github.com/(.*)/(.*).git")

func main() {
	flag.Bool("help", false, "Show usage.")

	hub := flag.Bool("hub", false, "Interact with GitHub API. Must be set for -user and -pass to have effect.")
	user := flag.String("user", "", "Auth username. Default 'git credential' username for -url.")
	pass := flag.String("pass", "", "Auth password. Default 'git credential' password for -url.")

	bref := flag.String("bref", "", "Base branch. Default master.")
	bsha := flag.String("bsha", "", "Base SHA. Default SHA of master/HEAD.")
	ref := flag.String("ref", "", "Branch to check out. Default current branch.")
	sha := flag.String("sha", "", "SHA to reset to. Default SHA of current branch.")

	url := flag.String("url", "", "Canonical repo URL. Default https://github.com/owner/repo.git remote.")

	flag.Usage = func() {
		fmt.Printf("usage: %s [options] [<directory>]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// find path to local .git directory
	dir := "."
	if len(flag.Args()) == 1 {
		dir = flag.Args()[0]
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	ok, err := existsLocal(dir)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Printf("ERROR: %s does not exist\n", dir)
		os.Exit(1)
	}

	ok, err = existsLocal(dir + "/.git")
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Printf("ERROR: %s does not appear to be a git repo\n", dir)
		os.Exit(1)
	}

	// build up Docker args
	args := []string{"run",
		"-e", "USER",
		"-e", "PASS",
		"-e", "BREF",
		"-e", "BSHA",
		"-e", "REF",
		"-e", "SHA",
		"-e", "URL",
		"-v", fmt.Sprintf("%s:%s", dir, "/tmp/repo"),
	}

	// infer empty args
	if *ref == "" {
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			panic(err)
		}
		*ref = strings.TrimSpace(string(out))
	}

	if *sha == "" {
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			panic(err)
		}
		*sha = strings.TrimSpace(string(out))
	}

	if *bref == "" {
		*bref = "master"
	}

	if *bsha == "" {
		out, err := exec.Command("git", "rev-parse", *bref).Output()
		if err != nil {
			panic(err)
		}
		*bsha = strings.TrimSpace(string(out))
	}

	// infer a github.com repo URL
	if *url == "" {
		out, err := exec.Command("git", "remote", "-v").Output()
		if err != nil {
			panic(err)
		}
		if m := reGH.FindString(string(out)); m != "" {
			*url = m
		}
	}
	if *url == "" {
		fmt.Printf("ERROR: No https://github.com/owner/repo.git remote found\n")
		os.Exit(1)
	}

	// infer credentials
	if *hub && (*pass == "" || *user == "") {
		u, err := neturl.Parse(*url)
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("git", "credential", "fill")
		cmd.Stdin = strings.NewReader(fmt.Sprintf("protocol=%s\nhost=%s\n", u.Scheme, u.Host))
		out, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		info := map[string]string{}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), "=")
			info[parts[0]] = parts[1]
		}

		if *pass == "" {
			*pass = info["password"]
		}

		if *user == "" {
			*user = info["username"]
		}
	}

	dpass := "<secret>"
	duser := *user

	if !*hub {
		// discard arguments
		dpass = "<disabled>"
		duser = "<disabled>"
		*pass = ""
		*user = ""
	} else {
		if *pass == "" {
			dpass = "<empty>"
		}
	}

	fmt.Printf("DIR:  %s\nUSER: %s\nPASS: %s\nBREF: %s\nBSHA: %s\nREF:  %s\nSHA:  %s\nURL:  %s\n\n", dir, duser, dpass, *bref, *bsha, *ref, *sha, *url)

	args = append(args, "mixable/bios", "runner")
	cmd := exec.Command("docker", args...)
	cmd.Env = []string{
		fmt.Sprintf("USER=%s", *user),
		fmt.Sprintf("PASS=%s", *pass),
		fmt.Sprintf("BREF=%s", *bref),
		fmt.Sprintf("BSHA=%s", *bsha),
		fmt.Sprintf("REF=%s", *ref),
		fmt.Sprintf("SHA=%s", *sha),
		fmt.Sprintf("URL=%s", *url),
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func existsLocal(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
