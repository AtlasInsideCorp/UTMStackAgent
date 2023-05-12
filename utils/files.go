package utils

import (
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// writeToFile writes the given `body` string to a file with the specified `fileName`.
// If the file does not exist, it will be created. If the file already exists, its contents
// will be overwritten with the new `body` string.
func WriteToFile(fileName string, body string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(body)
	return err
}

// WriteYAML writes the provided data to the specified file URL in YAML format.
// Returns an error if any error occurs during the process.
func WriteYAML(url string, data interface{}) error {
	config, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	err = WriteToFile(url, string(config[:]))
	if err != nil {
		return err
	}
	return nil
}

// ReadYAML reads the YAML data from the specified file URL and deserializes it into the provided result interface{}.
// Returns an error if any error occurs during the process.
func ReadYAML(url string, result interface{}) error {
	f, err := os.Open(url)
	if err != nil {
		return err
	}
	defer f.Close()
	d := yaml.NewDecoder(f)
	if err := d.Decode(result); err != nil {
		return err
	}
	return nil
}

// GetMyPath returns the directory path where the currently running executable is located.
// Returns a string representing the directory path, and an error if any error occurs during the process.
func GetMyPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

// GenerateFromTemplate generates a file from a template file using the provided data.
// Returns an error if any error occurs during the process.
func GenerateFromTemplate(data interface{}, tfile string, cfile string) error {
	_, fileName := filepath.Split(tfile)
	ut, err := template.New(fileName).ParseFiles(tfile)
	if err != nil {
		return err
	}
	writer, err := os.OpenFile(cfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	err = ut.Execute(writer, data)
	if err != nil {
		return err
	}
	return nil
}
