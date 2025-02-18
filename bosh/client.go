package bosh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/EngineerBetter/control-tower/resource"

	"github.com/EngineerBetter/control-tower/iaas"

	"github.com/EngineerBetter/control-tower/terraform"

	"github.com/EngineerBetter/control-tower/bosh/internal/boshcli"
	"github.com/EngineerBetter/control-tower/bosh/internal/workingdir"
	"github.com/EngineerBetter/control-tower/config"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// StateFilename is default name for bosh-init state file
const StateFilename = "director-state.json"

// CredsFilename is default name for bosh-init creds file
const CredsFilename = "director-creds.yml"

//counterfeiter:generate . IClient
// IClient is a client for performing bosh-init commands
type IClient interface {
	Deploy([]byte, []byte, bool) ([]byte, []byte, error)
	Cleanup() error
	Instances() ([]Instance, error)
	CreateEnv([]byte, []byte, string) ([]byte, []byte, error)
	Recreate() error
	Locks() ([]byte, error)
}

// Instance represents a vm deployed by BOSH
type Instance struct {
	Name  string
	IP    string
	State string
}

// ClientFactory creates a new IClient
type ClientFactory func(config config.ConfigView, outputs terraform.Outputs, stdout, stderr io.Writer, provider iaas.Provider, versionFile []byte) (IClient, error)

//New returns an IAAS specific implementation of BOSH client
func New(config config.ConfigView, outputs terraform.Outputs, stdout, stderr io.Writer, provider iaas.Provider, versionFile []byte) (IClient, error) {
	workingdir, err := workingdir.New()
	if err != nil {
		return nil, err
	}

	boshCLIPath, err := resource.DownloadBOSHCLI()
	if err != nil {
		return nil, fmt.Errorf("failed to determine BOSH CLI path: [%v]", err)
	}

	boshCLI := boshcli.New(boshCLIPath, exec.Command)

	switch provider.IAAS() {
	case iaas.AWS:
		return NewAWSClient(config, outputs, workingdir, stdout, stderr, provider, boshCLI)
	case iaas.GCP:
		return NewGCPClient(config, outputs, workingdir, stdout, stderr, provider, boshCLI)
	}
	return nil, fmt.Errorf("IAAS not supported: %s", provider.IAAS())
}

func instances(boshCLI boshcli.ICLI, ip, password, ca string) ([]Instance, error) {
	output := new(bytes.Buffer)

	if err := boshCLI.RunAuthenticatedCommand(
		"instances",
		ip,
		password,
		ca,
		false,
		output,
		"--json",
	); err != nil {
		return nil, fmt.Errorf("Error [%s] running `bosh instances`. stdout: [%s]", err, output.String())
	}

	jsonOutput := struct {
		Tables []struct {
			Rows []struct {
				Instance     string `json:"instance"`
				IPs          string `json:"ips"`
				ProcessState string `json:"process_state"`
			} `json:"Rows"`
		} `json:"Tables"`
	}{}

	if err := json.NewDecoder(output).Decode(&jsonOutput); err != nil {
		return nil, err
	}

	instances := []Instance{}

	for _, table := range jsonOutput.Tables {
		for _, row := range table.Rows {
			instances = append(instances, Instance{
				Name:  row.Instance,
				IP:    row.IPs,
				State: row.ProcessState,
			})
		}
	}

	return instances, nil
}

func saveFilesToWorkingDir(workingdir workingdir.IClient, provider iaas.Provider, creds []byte) error {
	concourseVersionsContents, _ := provider.Choose(iaas.Choice{
		AWS: awsConcourseVersions,
		GCP: gcpConcourseVersions,
	}).([]byte)
	concourseSHAsContents, _ := provider.Choose(iaas.Choice{
		AWS: awsConcourseSHAs,
		GCP: gcpConcourseSHAs,
	}).([]byte)

	filesToSave := map[string][]byte{
		concourseVersionsFilename:      concourseVersionsContents,
		concourseSHAsFilename:          concourseSHAsContents,
		concourseManifestFilename:      concourseManifestContents,
		concourseCompatibilityFilename: concourseCompatibility,
		concourseGrafanaFilename:       concourseGrafana,
		concourseGitHubAuthFilename:    concourseGitHubAuth,
		credsFilename:                  creds,
		extraTagsFilename:              extraTags,
	}

	for filename, contents := range filesToSave {
		_, err := workingdir.SaveFileToWorkingDir(filename, contents)
		if err != nil {
			return fmt.Errorf("failed to save %s to working directory: [%v]", filename, err)
		}
	}
	return nil
}
