package beat

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"github.com/quantfall/holmes"
)

type Beat interface {
	Install(config BeatConfig, h *holmes.Logger) error
	Run(h *holmes.Logger) error
	Uninstall(h *holmes.Logger) error
}

// NewBeatsInstance creates a new Beats instance.
func NewBeatsInstance() (*[]Beat, error) {
	var beats []Beat
	switch runtime.GOOS {
	case "windows":
		beats = append(beats, &Winlogbeat{})
	case "linux":
		beats = append(beats, &Filebeat{})
	default:
		return nil, fmt.Errorf("operative system no supported")
	}
	return &beats, nil
}

type BeatConfig struct {
	IP         string
	CA         string
	ClientCert string
	ClientKey  string
	Path       string
}

// InstallBeats installs, configures adn runs Beats on the specified IP address.
// Returns an error in case of failure to install or configure Beats.
func InstallBeats(ip string, cons configuration.ConstConfig, h *holmes.Logger) {
	path, err := utils.GetMyPath()
	if err != nil {
		h.FatalError("error getting current path: %v", err)
	}

	configBeat := BeatConfig{
		IP:         ip,
		CA:         filepath.Join(path, "certs", cons.TLSCA),
		ClientCert: filepath.Join(path, "certs", cons.TLSCRT),
		ClientKey:  filepath.Join(path, "certs", cons.TLSKEY),
	}

	beats, err := NewBeatsInstance()
	if err != nil {
		h.FatalError("error getting beats instance: %v", err)
	}
	var wg sync.WaitGroup

	for _, b := range *beats {
		beat := b
		wg.Add(1)
		go func() {
			err = beat.Install(configBeat, h)
			if err != nil {
				fmt.Printf("error configuring beats: %v", err)
				h.FatalError("error configuring beats: %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// RunBeats runs previously installed beats.
func RunBeats(h *holmes.Logger) {
	beats, err := NewBeatsInstance()
	if err != nil {
		h.FatalError("error getting beats instance: %v", err)
	}
	var wg sync.WaitGroup

	for _, b := range *beats {
		beat := b
		wg.Add(1)
		go func() {
			err := beat.Run(h)
			if err != nil {
				fmt.Printf("error running beats: %v", err)
				h.FatalError("error running beats: %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// DisableBeats disables previously installed beats.
func DisableBeats(h *holmes.Logger) {
	beats, err := NewBeatsInstance()
	if err != nil {
		h.FatalError("error getting beats instance: %v", err)
	}
	var wg sync.WaitGroup

	for _, b := range *beats {
		beat := b
		wg.Add(1)
		go func() {
			err := beat.Uninstall(h)
			if err != nil {
				fmt.Printf("error disabling beats: %v", err)
				h.FatalError("error disabling beats: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
