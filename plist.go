package launchd

import (
	"bytes"
	"fmt"
	"os"

	// Imported for embedded launchagent template below.
	_ "embed"
	"text/template"
)

//go:embed launchagent.tmpl
var plistTemplateSrc string
var plistTemplate = template.Must(
	template.New("launchctl").Parse(plistTemplateSrc),
)

// RenderPlist returns the contents of a plist file content for a launchd service.
func (s *Service) RenderPlist() ([]byte, error) {
	var buf bytes.Buffer
	err := plistTemplate.Execute(&buf, *s)
	return buf.Bytes(), err
}

// WritePlist writes the contents of the service's plist file content to the appropriate directory.
func (s *Service) WritePlist(content []byte) error {
	content, err := s.RenderPlist()
	if err != nil {
		return fmt.Errorf("could not render Plist: %w", err)
	}
	path, err := s.DefinitionPath()
	if err != nil {
		return fmt.Errorf("could not find definition path: %w", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", path, err)
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}
