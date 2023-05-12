package beat

import (
	"fmt"
	"path/filepath"

	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"github.com/quantfall/holmes"
)

type Winlogbeat struct {
}

// Generates the Winlogbeat configuration file from a template using the provided values
// Starts the Winlogbeat service on Windows with the settings specified in the winlogbeat.yml file.
func (w *Winlogbeat) Install(config BeatConfig, h *holmes.Logger) error {
	path, err := utils.GetMyPath()
	if err != nil {
		return err
	}
	// Configure Winlogbeat
	configFile := filepath.Join(path, "beats", "winlogbeat", "winlogbeat.yml")
	templateFile := filepath.Join(path, "templates", "winlogbeat.yml")
	err = utils.GenerateFromTemplate(config, templateFile, configFile)
	if err != nil {
		return err
	}
	fmt.Println("Winlogbeat was configured correctly.")
	h.Info("winlogbeat was configured correctly.")
	return nil
}

// Start Winlogbeat
func (w *Winlogbeat) Run(h *holmes.Logger) error {
	path, err := utils.GetMyPath()
	if err != nil {
		return err
	}

	running, err := utils.IsProcessRunning("winlogbeat.exe")
	if err != nil {
		return fmt.Errorf("error checking if winlogbeat is running: %v", err)
	}

	if running {
		err = utils.StopProcess("winlogbeat.exe")
		if err != nil {
			return fmt.Errorf("error stoping winlogbeat: %v", err)
		}
	}

	result, errB := utils.Execute(
		filepath.Join(path, "beats", "winlogbeat", "winlogbeat.exe"),
		filepath.Join(path, "beats", "winlogbeat"),
		"--strict.perms=false",
		"-c",
		"winlogbeat.yml",
	)
	if errB {
		return fmt.Errorf("error running winlogbeat: %s", result)
	}
	fmt.Println("Winlogbeat ran successfully.")
	h.Info("winlogbeat ran successfully.")

	return nil
}

// Stops the Winlogbeat.
func (w *Winlogbeat) Uninstall(h *holmes.Logger) error {
	return utils.StopProcess("winlogbeat.exe")
}
