package launchd

import (
	"fmt"
	"time"

	"github.com/brasic/launchd/state"
)

// Install sets up a new service by writing a plist file and telling launchd about it.
// It also starts the service and waits for it to come up if RunAtLoad is true.
func (s *Service) Install() (err error) {
	if s.InstallState().Is(state.Installed) {
		// Nothing to do
		return nil
	}
	content, err := s.RenderPlist()
	if err != nil {
		return err
	}

	if err = s.WritePlist(content); err != nil {
		return err
	}

	if _, err = s.Bootstrap(); err != nil {
		return err
	}

	if s.RunAtLoad && s.waitUntilRunning(5*time.Second) {
		return nil
	}
	return fmt.Errorf("Timed out waiting for service to boot")
}

func (s *Service) waitUntilRunning(timeout time.Duration) bool {
	_, timedOut := s.PollUntil(state.Running, timeout)
	fmt.Printf(finalStatus(timedOut), s.UserSpecifier())
	return !timedOut
}

func finalStatus(timedOut bool) string {
	if timedOut {
		return "timed out waiting for service to come up. Something is probably wrong.\nRun launchctl print %s` for more detail.\n"
	}
	return "done!\nRun launchctl print %s` for more detail.\n"
}
