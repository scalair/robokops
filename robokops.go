package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// Global var
var runningOrder []string
var containerID string
var ctx context.Context
var dockerCli *client.Client
var featuresMap map[string]map[string]string
var robokopsVersion string

// Colors
var boldGreen *color.Color
var boldRed *color.Color
var boldWhite *color.Color

// Arguments
var config *string
var terraformAction *string
var kubeAction *string
var ssh *string
var dev *bool
var kubeTargets *[]string
var env *[]string

// Bom struct
type Bom struct {
	Version  string
	Features []Feature
}
type Feature struct {
	Name    string
	Image   string
	Version string
}

func main() {
	// Parse the command line arguments..
	argParse()
	// Parse the bom.yaml file
	parseBom()
	// Initialise some stuff
	initialise()

	// Go Robokops !
	figure.NewFigure("Robokops", "", true).Print()

	// Run the features
	if *terraformAction == "destroy" {
		// When destroying the cluster, start by deleting the features
		runKubernetes("delete")
		runTerraform("destroy")
	} else {
		if *terraformAction != "" {
			runTerraform(*terraformAction)
		}
		if *kubeAction != "" {
			runKubernetes(*kubeAction)
		}
	}

	// Done
	breakLine()
	boldGreen.Println("Robokops finished successfully")
}

// Arguments parser
func argParse() {
	parser := argparse.NewParser("Robokops", "Manage Kubernetes clusters and deploy common features")

	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Return Robokops current version"})
	config = parser.String("c", "config", &argparse.Options{Required: false, Help: "Customer configuration folder. Must contains a \"conf\" and a \"terraform\" subfolders"})
	terraformAction = parser.String("T", "terraform", &argparse.Options{Required: false, Help: "Plan, apply or destroy infrastructure. Choose between: plan|apply|destroy"})
	kubeAction = parser.String("a", "action", &argparse.Options{Required: false, Help: "Action to execute. Choose between: deploy|delete|dry-run"})
	ssh = parser.String("s", "ssh", &argparse.Options{Required: false, Help: "Path of the .ssh directory (use only by Terraform to clone private modules)"})
	dev = parser.Flag("d", "dev", &argparse.Options{Required: false, Help: "Add this flag to use local docker image instead of the remote registry"})
	kubeTargets = parser.List("t", "target", &argparse.Options{Required: false, Help: "Target of the action. If not provided will execute against all matching configuration"})
	env = parser.List("e", "env", &argparse.Options{Required: false, Help: "Define environment variables to pass to containers. You can use \"--env all\" to map all env vars available in your OS context"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// Return the current version
	if *version {
		parseBom()
		fmt.Println("Version: " + robokopsVersion)
		os.Exit(0)
	}

	if !*version && *config == "" {
		fmt.Print(parser.Usage("Not enough arguments"))
		os.Exit(1)
	}

	// Special behavior, if env is passed with the value "all"
	// then env will be set with all environment variables
	// available in the OS
	if len(*env) == 1 && (*env)[0] == "all" {
		*env = []string{}
		for _, envVar := range os.Environ() {
			*env = append(*env, envVar)
		}
	}
}

// Parse the bom.yaml file and set featuresMap
func parseBom() {
	var bom Bom
	var err error
	var source []byte
	if _, err := os.Stat("bom.yaml"); err == nil {
		source, err = ioutil.ReadFile("bom.yaml")
	} else if _, err := os.Stat(os.Getenv("GOPATH") + "/src/github.com/scalair/robokops/bom.yaml"); err == nil {
		source, err = ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/scalair/robokops/bom.yaml")
	} else {
		ppError("Cannot find bom.yaml file", err)
	}

	if err != nil {
		ppError("Fail to read the bom.yaml file", err)
	}
	if err = yaml.Unmarshal(source, &bom); err != nil {
		ppError("Fail to parse the bom.yaml file", err)
	}

	featuresMap = map[string]map[string]string{}
	for i := 0; i < len(bom.Features); i++ {
		featureMap := map[string]string{}
		featureMap["image"] = bom.Features[i].Image
		featureMap["version"] = bom.Features[i].Version
		featuresMap[bom.Features[i].Name] = featureMap
	}

	robokopsVersion = bom.Version
}

// Init func
func initialise() {
	runningOrder = []string{
		"cluster-init",
		"cluster-autoscaler",
		"monitoring",
		"elastic-stack",
		"external-dns",
		"ingress-nginx",
		"aws-alb-ingress-controller",
		"gitlabci",
		"dashboard",
	}

	// Catch sigterm signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	// Docker client
	ctx = context.Background()
	var err error
	dockerCli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		ppError("Failed to use docker client. Make sure docker is installed and running.", err)
	}

	// Define awesome colors
	boldGreen = color.New(color.FgGreen, color.Bold)
	boldWhite = color.New(color.FgWhite, color.Bold)
	boldRed = color.New(color.FgRed, color.Bold)

	setTerminalWidth()
}

// terraform feature
func runTerraform(action string) {
	if *ssh == "" {
		*ssh = os.Getenv("HOME") + "/.ssh"
	}

	configAsbPath, _ := filepath.Abs(*config)
	binds := []string{
		configAsbPath + "/terraform:/conf/terraform",
		*ssh + ":/ssh",
	}

	runContainer("terraform", action, binds)
}

// k8s features
func runKubernetes(action string) {
	// When action is delete, reverse the order of execution
	if action == "delete" {
		runningOrder = reverse(runningOrder)
	}

	configAsbPath, _ := filepath.Abs(*config)
	binds := []string{
		configAsbPath + "/conf:/conf",
		"/tmp:/local",
	}

	// If kubeTargets is empty, set it with the list of directory in *config/conf
	if len(*kubeTargets) == 0 {
		dirs, err := ioutil.ReadDir(*config + "/conf")
		if err != nil {
			ppError("Fail to list directory in :"+*config+"/conf", err)
		}
		for _, d := range dirs {
			if d.IsDir() {
				*kubeTargets = append(*kubeTargets, d.Name())
			}
		}
	}

	// Create the features list by ordering kubeTargets based on runningOrder
	var features []string
	for _, feature := range runningOrder {
		if contains(*kubeTargets, feature) {
			features = append(features, feature)
		}
	}

	// Run the features, one at a time
	for _, feature := range features {
		runContainer(feature, action, binds)
	}
}

// Run the feature
func runContainer(feature string, action string, binds []string) {
	// In dev mode, use the local latest image
	var image string
	if *dev {
		image = strings.Split(featuresMap[feature]["image"], "/")[1] + ":latest"
	} else {
		image = featuresMap[feature]["image"] + ":" + featuresMap[feature]["version"]
	}

	envList := ""
	for _, e := range *env {
		envList += e + ";"
	}

	// Print what is gonna run
	breakLine()
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Feature", "Action", "Image"})
	t.AppendRow([]interface{}{feature, action, image})
	boldGreen.Println(t.Render() + "\n")

	// Pull image
	if !*dev {
		color.Green("Pulling image " + image)
		out, err := dockerCli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			ppError("Failed to pull image: "+image, err)
		}
		defer out.Close()
		// The io.ReadCloser object return must be read, otherwise
		// the docker image won't be pulled
		rd := bufio.NewReader(out)
		for {
			_, err := rd.ReadString('\n')
			if err == io.EOF {
				break
			}
		}
		fmt.Print("\n")
	}

	// Create container
	color.Green("Creating container\n\n")
	resp, err := dockerCli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Cmd:   strings.Fields(action),
			Env:   *env,
		},
		&container.HostConfig{
			Binds: binds,
		},
		nil,
		"",
	)
	if err != nil {
		ppError("Failed to create container for image: "+image, err)
	}
	containerID = resp.ID

	// Start container
	color.Green("Starting container with id [" + containerID + "]\n\n")
	if err := dockerCli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		ppError("Failed to start container: "+containerID, err)
	}

	// Output container logs
	out, err := dockerCli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		ppError("Failed to show logs of container: "+containerID, err)
	}
	defer out.Close()
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

// ----------------------
// --- Utils function ---
// ----------------------

// Pretty print error
func ppError(msg string, err error) {
	boldRed.Printf("\nError: ")
	boldWhite.Println(msg)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// Calculate the width of the terminal screen
var terminalWidth int

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Set the terminal with which can be use to display break line
func setTerminalWidth() {
	ws := &winsize{}
	retCode, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		fmt.Println("Failed to calculate terminal width (will use default value 100): ", err)
		terminalWidth = 100
	} else {
		terminalWidth = int(uint(ws.Col))
	}
}

// Print a break line
func breakLine() {
	fmt.Println(strings.Repeat("┅", terminalWidth))
}

// Return true if element is in strSlice
func contains(strSlice []string, element string) bool {
	for _, str := range strSlice {
		if str == element {
			return true
		}
	}
	return false
}

// Reverse the order of strSlice
func reverse(strSlice []string) []string {
	reversedSlice := []string{}
	for i := range strSlice {
		n := strSlice[len(strSlice)-1-i]
		reversedSlice = append(reversedSlice, n)
	}
	return reversedSlice
}

// Remove the running container after receiving sigterm signal
func cleanup() {
	fmt.Println("\nKill signal received. Remove container:", containerID)
	if err := dockerCli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		ppError("Fail to terminate the container: "+containerID, err)
	}
}
