package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4"
	git "gopkg.in/src-d/go-git.v4"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// Runtime regroups runtime parameters
type Runtime struct {
	basePreviousTag string
	baseCurrentTag  string
	githubToken     string
}

var (
	runtime  = &Runtime{}
	repo     = "scalair/robokops"
	features = []string{
		"terraform",
    "cluster-init",
		"cluster-autoscaler",
		"monitoring",
		"elastic-stack",
		"external-dns",
		"ingress-nginx",
		"aws-alb-ingress-controller",
		"aws-efs-csi-driver",
		"gitlabci",
		"dashboard",
		"velero",
		"kubewatch",
		"jenkins",
	}
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("=== Robokops Packaging Script ===")

	initEnv()

	for _, feature := range features {
		log.Infof("--- Packaging %s ---", feature)

		// in-memory filesystem abstraction used for cloning
		// repo without relying on actual disk filesystem
		fs := memfs.New()

		// differentiate terraform feature from others
		featurePath := feature
		if feature != "terraform" {
			featurePath = "k8s/" + feature
		}

		// clone repo in memory
		r, err := clone(fs)
		if err != nil {
			log.Fatal(err)
		}

		// get latest version from changelog
		version, err := findInFile(
			&fs,
			featurePath+"/CHANGELOG.md",
			`v?[0-9]+\.[0-9]+\.[0-9]+`,
		)
		if err != nil {
			log.Fatal(err)
		}

		// increment version by using robokops-base version
		version, err = newVersion(version)
		if err != nil {
			log.Fatal(err)
		}

		// create a branch named after the feature and its new version
		// and checkout to it
		w, err := checkout(r, feature+"/"+version)
		if err != nil {
			log.Fatal(err)
		}

		// update robokops-base version in dockerfile
		err = updateFile(
			&fs,
			featurePath+"/Dockerfile",
			"scalair/robokops-base:.*",
			"scalair/robokops-base:"+runtime.baseCurrentTag,
		)
		if err != nil {
			log.Fatal(err)
		}

		// add a new entry to the changelog
		err = updateFile(
			&fs,
			featurePath+"/CHANGELOG.md",
			"^# Changelog",
			fmt.Sprintf("# Changelog\n\n## %s - %s\n### Changed\n- Release %s", version, time.Now().Format("2006-01-02"), version),
		)
		if err != nil {
			log.Fatal(err)
		}

		// commit and push all changes
		err = add(w, ".")
		if err != nil {
			log.Fatal(err)
		}
		err = commit(w, fmt.Sprintf("Release %s %s", feature, version))
		if err != nil {
			log.Fatal(err)
		}
		err = push(r)
		if err != nil {
			log.Fatal(err)
		}

		err = pullRequest(feature, version, feature+"/"+version)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// parse command-line and environment to retrieve needed parameters
func initEnv() error {
	flag.StringVar(&runtime.basePreviousTag, "p", runtime.basePreviousTag, "Previous tag of robokops-base image")
	flag.StringVar(&runtime.baseCurrentTag, "c", runtime.baseCurrentTag, "Current tag of robokops-base image")
	flag.StringVar(&runtime.githubToken, "g", runtime.githubToken, "Github token to branch, push and merge to Robokops repo")

	if val, ok := os.LookupEnv("BASE_PREVIOUS_TAG"); ok {
		runtime.basePreviousTag = val
	}
	if val, ok := os.LookupEnv("BASE_CURRENT_TAG"); ok {
		runtime.baseCurrentTag = val
	}
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		runtime.githubToken = val
	}

	flag.Parse()

	if runtime.basePreviousTag == "" {
		return errors.New("Previous tag command-line parameter is not set")
	}
	if runtime.baseCurrentTag == "" {
		return errors.New("Current tag command-line parameter is not set")
	}
	if runtime.githubToken == "" {
		return errors.New("Github token is not set")
	}

	log.Debugf("Base previous tag %s", runtime.basePreviousTag)
	log.Debugf("Base current tag %s", runtime.baseCurrentTag)

	return nil
}

// find and extract a specific string in a file
func findInFile(fs *billy.Filesystem, filename, expression string) (string, error) {
	log.Info("Get current version")
	file, err := (*fs).Open(filename)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	re := regexp.MustCompile(expression)
	version := re.FindString(buf.String())
	file.Close()
	return version, nil
}

// which version digit to increment
func incrementChoice() (string, error) {
	previousTagTokens := strings.Split(runtime.basePreviousTag, ".")
	previousTagTokens[0] = strings.Trim(previousTagTokens[0], "v")

	currentTagTokens := strings.Split(runtime.baseCurrentTag, ".")
	currentTagTokens[0] = strings.Trim(currentTagTokens[0], "v")

	prevMaj, err := strconv.Atoi(previousTagTokens[0])
	if err != nil {
		return "", err
	}
	currMaj, err := strconv.Atoi(currentTagTokens[0])
	if err != nil {
		return "", err
	}
	if prevMaj < currMaj {
		return "major", nil
	}

	prevMin, err := strconv.Atoi(previousTagTokens[1])
	if err != nil {
		return "", err
	}
	currMin, err := strconv.Atoi(currentTagTokens[1])
	if err != nil {
		return "", err
	}
	if prevMin < currMin {
		return "minor", nil
	}

	prevPatch, err := strconv.Atoi(previousTagTokens[2])
	if err != nil {
		return "", err
	}
	currPatch, err := strconv.Atoi(currentTagTokens[2])
	if err != nil {
		return "", err
	}
	if prevPatch < currPatch {
		return "patch", nil
	}

	return "", errors.New("Could not increment version")
}

// increment version based on some rules
func newVersion(version string) (string, error) {
	log.Info("Increment version")
	hasV := strings.HasPrefix(version, "v")
	tokens := strings.Split(version, ".")
	tokens[0] = strings.Trim(tokens[0], "v")

	log.Debugf("Previous feature version: %s.%s.%s", tokens[0], tokens[1], tokens[2])

	digitMaj, _ := strconv.Atoi(tokens[0])
	digitMin, _ := strconv.Atoi(tokens[1])
	digitPatch, _ := strconv.Atoi(tokens[2])

	increment, err := incrementChoice()
	if err != nil {
		return "", err
	}

	if increment == "major" {
		digitMaj++
		digitMin = 0
		digitPatch = 0
	} else if increment == "minor" {
		digitMin++
		digitPatch = 0
	} else if increment == "patch" {
		digitPatch++
	} else {
		return "", errors.New("Previous tag and current tag are identical")
	}

	tokens[0] = strconv.Itoa(digitMaj)
	tokens[1] = strconv.Itoa(digitMin)
	tokens[2] = strconv.Itoa(digitPatch)

	if hasV {
		tokens[0] = "v" + tokens[0]
	}

	version = strings.Join(tokens, ".")

	log.Debugf("New feature version: %s", version)

	return version, nil
}

// clone git repo with Github token
func clone(fs billy.Filesystem) (*git.Repository, error) {
	log.Infof("Git clone %s", repo)
	r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "empty",
			Password: runtime.githubToken,
		},
		URL: fmt.Sprintf("https://github.com/%s.git", repo),
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// create a new branch and checkout to it
func checkout(r *git.Repository, branch string) (*git.Worktree, error) {
	log.Infof("Git branch %s ", branch)
	headRef, err := r.Head()
	if err != nil {
		return nil, err
	}
	ref := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+branch), headRef.Hash())
	err = r.Storer.SetReference(ref)
	if err != nil {
		return nil, err
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	err = w.Checkout(&git.CheckoutOptions{Branch: ref.Name()})
	if err != nil {
		return nil, err
	}
	return w, nil
}

// replace one line (src) of a file by another (dst)
func updateFile(fs *billy.Filesystem, filename, src, dst string) error {

	log.Infof("Update %s", filename)
	file, err := (*fs).Open(filename)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	re := regexp.MustCompile(src)
	cnt := re.ReplaceAllString(string(buf.Bytes()), dst)
	file.Close()

	newFile, err := (*fs).Create(filename)
	if err != nil {
		return err
	}
	_, err = newFile.Write([]byte(cnt))
	if err != nil {
		return err
	}
	newFile.Close()

	return nil
}

// add changes of filename in staging area
func add(w *git.Worktree, filename string) error {
	log.Infof("Git add %s", filename)
	_, err := w.Add(filename)
	if err != nil {
		return err
	}
	status, _ := w.Status()
	log.Debugf("Add status:\n%s", status)
	return nil
}

// commit staged changes
func commit(w *git.Worktree, message string) error {
	log.Infof("Git commit -m '%s'", message)
	_, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "auto",
			Email: "cloudops@scalair.fr",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	status, _ := w.Status()
	log.Debugf("Commit status:\n%s", status)

	return nil
}

// push branch to repo
func push(r *git.Repository) error {
	log.Infof("Git push")
	err := r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "empty",
			Password: runtime.githubToken,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// create a pull request for the feature in Github
func pullRequest(feature, version, branch string) error {
	log.Info("Pull request")
	authToken := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: runtime.githubToken})
	authClient := oauth2.NewClient(context.Background(), authToken)

	client := github.NewClient(authClient)

	pr := &github.NewPullRequest{
		Title:               github.String("Release " + feature + " " + version),
		Head:                github.String(branch),
		Base:                github.String("master"),
		Body:                github.String("Pull request created automatically by the Robokops features packager script"),
		MaintainerCanModify: github.Bool(true),
	}

	orgRepo := strings.Split(repo, "/")
	if len(orgRepo) != 2 {
		return errors.New("Repository name must be to the form 'owner/repository'")
	}

	_, _, err := client.PullRequests.Create(context.Background(), orgRepo[0], orgRepo[1], pr)
	if err != nil {
		return err
	}

	return nil
}
