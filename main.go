package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	flag.Bool("help", false, "show usage")
	ref := flag.String("ref", "", "branch to check out")
	sha := flag.String("sha", "", "SHA to set status for")
	url := flag.String("url", "", "repo to clone")

	user := flag.String("user", "", "auth username")
	pass := flag.String("pass", "", "auth password")

	flag.Usage = func() {
		fmt.Printf("usage: %s [options] [path-to-script]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

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

	if *url == "" {
		out, err := exec.Command("git", "ls-remote", "--get-url", "origin").Output()
		if err != nil {
			panic(err)
		}
		*url = strings.TrimSpace(string(out))
	}

	if *pass == "" {
		cmd := exec.Command("git", "credential", "fill")
		cmd.Stdin = strings.NewReader("protocol=https\nhost=github.com\n")
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

		*pass = info["password"]
		*user = info["username"]
	}

	args := []string{"run", "-e", "USER", "-e", "PASS", "-e", "REF", "-e", "SHA", "-e", "URL"}

	var err error
	script := "bios.sh"
	if len(flag.Args()) == 1 {
		script, err = filepath.Abs(flag.Args()[0])
		if err != nil {
			panic(err)
		}
		args = append(args, "-v", fmt.Sprintf("%s:%s", script, "/tmp/bios.sh"))
	}

	fmt.Printf("Ref: %s\nSHA: %s\nURL: %s\nsh : %s\n", *ref, *sha, *url, script)

	args = append(args, "mixable/bios", "handler.Handle")
	cmd := exec.Command("docker", args...)
	cmd.Env = []string{
		fmt.Sprintf("USER=%s", *user),
		fmt.Sprintf("PASS=%s", *pass),
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
	fmt.Printf("ERR: %+v\n", err)
}
