package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4"
	git "gopkg.in/src-d/go-git.v4"

	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

// Runtime regroups runtime parameters
type Runtime struct {
	bomPath         string
	basePreviousTag string
	baseCurrentTag  string
	githubToken     string
}

// Bom partial definition
type Bom struct {
	Version  string `yaml:"version"`
	Features []struct {
		Name    string `yaml:"name"`
		Image   string `yaml:"image"`
		Version string `yaml:"version"`
	} `yaml:"features"`
}

var (
	runtime = &Runtime{}
	repo    = "scalair/robokops"
	// used in commit signature:
	authorName  = "scalaircloudops"
	authorEmail = "cloudops@scalair.fr"
)

func main() {
	log.SetLevel(log.TraceLevel)
	log.Info("=== Robokops Features Packager Script ===")

	initEnv()

	bom, err := parseBom(runtime.bomPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, feature := range bom.Features {
		featureName := feature.Name
		log.Infof("--- Packaging %s ---", featureName)

		// in-memory filesystem abstraction used for cloning
		// repo without relying on actual disk filesystem
		fs := memfs.New()

		// differentiate terraform feature from others
		featurePath := featureName
		if featureName != "terraform" {
			featurePath = "k8s/" + featureName
		}

		// clone repo in memory
		r, err := clone(fs)
		if err != nil {
			log.Error(err)
			continue
		}

		// get latest version from changelog
		version, err := findInFile(
			&fs,
			featurePath+"/CHANGELOG.md",
			`v?[0-9]+\.[0-9]+\.[0-9]+`,
		)
		if err != nil {
			log.Error(err)
			continue
		}

		// increment version by using robokops-base version
		version, err = increment(version)
		if err != nil {
			log.Error(err)
			continue
		}

		// create a branch named after the feature and its new version
		// and checkout to it
		w, err := checkout(r, featureName+"/"+version)
		if err != nil {
			log.Error(err)
			continue
		}

		// update robokops-base version in dockerfile
		err = updateFile(
			&fs,
			featurePath+"/Dockerfile",
			"scalair/robokops-base:.*",
			"scalair/robokops-base:"+runtime.baseCurrentTag,
		)
		if err != nil {
			log.Error(err)
			continue
		}

		// add a new entry to the changelog
		err = updateFile(
			&fs,
			featurePath+"/CHANGELOG.md",
			"^# Changelog",
			fmt.Sprintf("# Changelog\n\n## %s - %s\n### Changed\n- Release %s", version, time.Now().Format("2006-01-02"), version),
		)
		if err != nil {
			log.Error(err)
			continue
		}
		// update feature version in bom
		for i, f := range bom.Features {
			if f.Name == featureName {
				log.Infof("Update bom file for %s to version %s", featureName, version)
				bom.Features[i].Version = version
				data, err := yaml.Marshal(&bom)
				if err != nil {
					log.Error(err)
					break
				}
				bomFile, err := fs.Create("bom.yaml")
				if err != nil {
					log.Error(err)
					break
				}
				defer bomFile.Close()
				_, err = bomFile.Write(data)
				if err != nil {
					log.Error(err)
					break
				}
				break
			}
		}
		if err != nil {
			log.Error(err)
			continue
		}

		// commit and push all changes
		err = add(w, ".")
		if err != nil {
			log.Error(err)
			continue
		}
		err = commit(w, fmt.Sprintf("Release %s %s", featureName, version))
		if err != nil {
			log.Error(err)
			continue
		}
		err = push(r)
		if err != nil {
			log.Error(err)
			continue
		}

		err = merge(featureName, version, featureName+"/"+version)
		if err != nil {
			log.Error(err)
			continue
		}

		rem, err := r.Remote("origin")
		if err != nil {
			log.Error(err)
			continue
		}
		err = deleteBranch(rem, featureName+"/"+version)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

// parse command-line and environment to retrieve needed parameters
func initEnv() error {
	flag.StringVar(&runtime.bomPath, "b", runtime.bomPath, "Path of bom.yaml file")
	flag.StringVar(&runtime.basePreviousTag, "p", runtime.basePreviousTag, "Previous tag of robokops-base image")
	flag.StringVar(&runtime.baseCurrentTag, "c", runtime.baseCurrentTag, "Current tag of robokops-base image")
	flag.StringVar(&runtime.githubToken, "g", runtime.githubToken, "Githunilb token to branch, push and merge to Robokops repo")

	if val, ok := os.LookupEnv("BOM_PATH"); ok {
		runtime.bomPath = val
	}
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

	if runtime.bomPath == "" {
		runtime.bomPath = "bom.yaml"
	}
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

// parse the bom.yaml file and set features list
func parseBom(bomPath string) (Bom, error) {
	// check if we need to use relative or absolute path
	if _, err := os.Stat(bomPath); os.IsNotExist(err) {
		absPath, err := os.Executable()
		if err != nil {
			return Bom{}, err
		}
		bomPath = filepath.Dir(absPath) + string(os.PathSeparator) + bomPath
		if _, err := os.Stat(bomPath); os.IsNotExist(err) {
			return Bom{}, err
		}
	}

	source, err := ioutil.ReadFile(bomPath)
	if err != nil {
		return Bom{}, err
	}

	var bom Bom
	err = yaml.Unmarshal(source, &bom)
	if err != nil {
		return Bom{}, err
	}

	return bom, nil
}

// find and extract a specific string in a file
func findInFile(fs *billy.Filesystem, filename, expression string) (string, error) {
	log.Infof("Find a string in the file %s", filename)
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

// increment version based on some rules
func increment(version string) (string, error) {
	log.Infof("Increment version")

	featureVersion, err := toVersion(version)
	if err != nil {
		return "", fmt.Errorf("Feature version must be to the form [v]x.y.z: %s", err.Error())
	}

	basePreviousVersion, err := toVersion(runtime.basePreviousTag)
	if err != nil {
		return "", errors.New("Previous base version must be to the form '[v]x.y.z', found: " + runtime.basePreviousTag)
	}

	baseNextVersion, err := toVersion(runtime.baseCurrentTag)
	if err != nil {
		return "", errors.New("Current base version must be to the form '[v]x.y.z', found: " + runtime.basePreviousTag)
	}

	newFeatureVersion := featureVersion.Increment(basePreviousVersion, baseNextVersion)

	log.Infof("Increment version %s to version %s", version, newFeatureVersion)

	return newFeatureVersion.String(), nil
}

// clone git repo with Github token
func clone(fs billy.Filesystem) (*git.Repository, error) {
	log.Infof("Clone the repository %s", repo)
	r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		Auth:     &http.BasicAuth{Username: "empty", Password: runtime.githubToken},
		URL:      fmt.Sprintf("https://github.com/%s.git", repo),
		Progress: nil,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// create a new branch and checkout to it
func checkout(r *git.Repository, branch string) (*git.Worktree, error) {
	log.Infof("Create the branch %s ", branch)
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
	log.Infof("Checkout the branch %s", branch)
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

	log.Debugf("File %s content:\n%s", filename, cnt)

	return nil
}

// add changes of filename in staging area
func add(w *git.Worktree, filename string) error {
	log.Infof("Add %s", filename)
	_, err := w.Add(filename)
	if err != nil {
		return err
	}
	return nil
}

// commit staged changes
func commit(w *git.Worktree, message string) error {
	log.Infof("Commit with message: %s", message)
	_, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// push branch to repo
func push(r *git.Repository) error {
	log.Infof("Push")
	err := r.Push(&git.PushOptions{
		Auth:     &http.BasicAuth{Username: "empty", Password: runtime.githubToken},
		Progress: nil,
	})
	if err != nil {
		return err
	}
	return nil
}

// create a pull request for the feature in Github and merge it to master
func merge(feature, version, branch string) error {
	splitRepo := strings.Split(repo, "/")
	if len(splitRepo) != 2 {
		return errors.New("Repository name must be to the form 'owner/repository'")
	}
	repoOwner := splitRepo[0]
	repoName := splitRepo[1]

	ctx := context.Background()

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: runtime.githubToken})))

	log.Infof("Pull request: Release %s %s", feature, version)
	pr, gr, err := client.PullRequests.Create(ctx, repoOwner, repoName, &github.NewPullRequest{
		Title:               github.String("Release " + feature + " " + version),
		Head:                github.String(branch),
		Base:                github.String("master"),
		MaintainerCanModify: github.Bool(true),
	},
	)
	if err != nil {
		return err
	}

	log.Trace(gr)

	prNum := pr.GetNumber()

	// wait until the branch is merged or timeout is reached
	log.Infof("Merge branch %s to master", branch)
	timeout := time.Now().Add(30 * time.Second)
	for {
		mr, gr, err := client.PullRequests.Merge(ctx, repoOwner, repoName, prNum, *pr.Title, &github.PullRequestOptions{MergeMethod: "squash"})
		if err == nil || time.Now().After(timeout) {
			log.Trace(mr)
			log.Trace(gr)
			break
		} else {
			log.Error(err)
			log.Info("Retrying...")
		}
	}
	log.Info("Merge done")

	return nil
}

// delete branch in the remote
func deleteBranch(remote *git.Remote, branch string) error {
	log.Infof("Delete branch %s", branch)
	return remote.Push(&git.PushOptions{
		Auth:     &http.BasicAuth{Username: "empty", Password: runtime.githubToken},
		RefSpecs: []config.RefSpec{config.RefSpec(":refs/heads/" + branch)},
		Progress: nil,
	})
}

//
// Version type
//

// Version to the form [v]x.y.z
type Version struct {
	Prefix bool // v
	Major  int
	Minor  int
	Patch  int
}

// Version to string
func (v Version) String() string {
	var prefix string
	if v.Prefix {
		prefix = "v"
	}
	return fmt.Sprintf("%s%d.%d.%d", prefix, v.Major, v.Minor, v.Patch)
}

// Increment version based on the update from base to next versions
func (v Version) Increment(base, next Version) Version {
	version := v

	if next.Major > base.Major {
		version.Major++
		version.Minor = 0
		version.Patch = 0
	} else if next.Minor > base.Minor {
		version.Minor++
		version.Patch = 0
	} else if next.Patch > base.Patch {
		version.Patch++
	}

	return version
}

func toVersion(strVersion string) (Version, error) {
	var version Version

	hasPrefix := strings.HasPrefix(strVersion, "v")
	strVersion = strings.Trim(strVersion, "v")

	tokens := strings.Split(strVersion, ".")
	if len(tokens) != 3 {
		return Version{}, errors.New("Incorrect version format")
	}

	major, err := strconv.Atoi(tokens[0])
	if err != nil {
		return Version{}, errors.New("Incorrect major version format")
	}
	minor, err := strconv.Atoi(tokens[1])
	if err != nil {
		return Version{}, errors.New("Incorrect minor version format")
	}
	patch, err := strconv.Atoi(tokens[2])
	if err != nil {
		return Version{}, errors.New("Incorrect patch version format")
	}

	version.Prefix = hasPrefix
	version.Major = major
	version.Minor = minor
	version.Patch = patch

	return version, nil
}
