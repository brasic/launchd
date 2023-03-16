package launchd

import (
	"fmt"
	"os"
	"path/filepath"
)

// Everything we deal with is limited to the current user's gui domain.
var domain = fmt.Sprintf("gui/%d", os.Getuid())

// Service is a LaunchAgent service.
type Service struct {
	// Name is the fully qualified name of the service e.g. com.id.doom
	Name string
	// ExecutablePath is the absolute path to the executable that will run the service
	ExecutablePath string
	// Argv is the list of arguments to pass to ExecutablePath
	Argv []string
	// RunAtLoad is whether the service should be started at login
	RunAtLoad bool
	// KeepAlive is whether the service should be restarted if it crashes
	KeepAlive bool
}

// ForRunningProgram returns a Service with appropriate daemon defaults for the current running executable.
func ForRunningProgram(name string, argv []string) *Service {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("Could not get executable path: %w", err))
	}
	return &Service{
		Name:           name,
		ExecutablePath: exe,
		Argv:           argv,
		RunAtLoad:      true,
		KeepAlive:      true,
	}
}

// UserSpecifier unambiguously identifies the service in subcommands.
// e.g. gui/501/com.id.doom
// See launchctl(1).
func (s *Service) UserSpecifier() string {
	return fmt.Sprintf("%s/%s", domain, s.Name)
}

// DefinitionPath is the absolute fs path where the service's plist config lives
func (s *Service) DefinitionPath() (string, error) {
	dir, err := launchAgentsDir()
	if err != nil {
		return "", err
	}
	plistFileName := s.Name + ".plist"
	return filepath.Join(dir, plistFileName), nil
}

func launchAgentsDir() (dir string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	dir = filepath.Join(home, "Library", "LaunchAgents")
	mkdirErr := os.MkdirAll(dir, 0700)
	if mkdirErr != nil {
		return "", fmt.Errorf("'%s' couldn't be created: %w", dir, err)
	}
	return
}
