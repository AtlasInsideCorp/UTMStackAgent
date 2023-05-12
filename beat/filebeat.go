package beat

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"github.com/quantfall/holmes"
)

type Filebeat struct {
}

// Configures Filebeat
// Generating configuration files from templates, and enabling the necessary modules.
func (f *Filebeat) Install(config BeatConfig, h *holmes.Logger) error {
	path, err := utils.GetMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "linux":
		configFile := filepath.Join("/", "etc", "filebeat", "filebeat.yml")
		templateFile := filepath.Join(path, "templates", "filebeat-linux.yml")
		family, err := utils.DetectLinuxFamily()
		if err != nil {
			return err
		}
		switch family {
		case "debian":
			result, err := utils.Execute("dpkg", filepath.Join(path, "beats", "filebeat"), "-i", filepath.Join(path, "beats", "filebeat", "filebeat-oss-8.5.3-amd64.deb"))
			if err {
				return fmt.Errorf("%s", result)
			}
		case "rhel":
			result, err := utils.Execute("yum", filepath.Join(path, "beats", "filebeat"), "localinstall", "-y", filepath.Join(path, "beats", "filebeat", "filebeat-oss-8.5.3-x86_64.rpm"))
			if err {
				return fmt.Errorf("%s", result)
			}
		}

		if family == "debian" || family == "rhel" {
			err = utils.GenerateFromTemplate(config, templateFile, configFile)
			if err != nil {
				return err
			}
			result, err := utils.Execute("filebeat", filepath.Join(path, "beats", "filebeat"), "modules", "enable", "system")
			if err {
				return fmt.Errorf("%s", result)
			}
			result, err = utils.Execute("systemctl", filepath.Join(path, "beats", "filebeat"), "enable", "filebeat")
			if err {
				return fmt.Errorf("%s", result)
			}
			result, err = utils.Execute("systemctl", filepath.Join(path, "beats", "filebeat"), "restart", "filebeat")
			if err {
				return fmt.Errorf("%s", result)
			}
		}
	}
	fmt.Println("Filebeat was configured correctly.")
	h.Info("filebeat was configured correctly.")
	return nil
}

// Runs the Filebeat service. This function doesn't do anything for Filebeat,
// as the service is started in the Install() function.
func (f *Filebeat) Run(h *holmes.Logger) error {
	return nil
}

// Stops and uninstalls the Filebeat service
func (f *Filebeat) Uninstall(h *holmes.Logger) error {
	path, err := utils.GetMyPath()
	if err != nil {
		return err
	}
	result, errB := utils.Execute("apt-get", filepath.Join(path, "beats", "filebeat"), "remove", "--purge", "-y", "filebeat")
	if errB {
		return fmt.Errorf("%s", result)
	}

	result, errB = utils.Execute("rm", filepath.Join(path, "beats", "filebeat"), "-rf", "/etc/filebeat")
	if errB {
		return fmt.Errorf("%s", result)
	}

	fmt.Println("Filebeat uninstalled successfully.")
	h.Info("filebeat uninstalled successfully.")
	return nil
}
