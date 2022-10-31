package launchd

import (
	"fmt"
	"time"

	"github.com/brasic/launchd/state"
)

// IsHealthy returns true if the service is running and healthy.
func (s *Service) IsHealthy() bool {
	return s.InstallState().Is(state.Installed) && s.RunState().Is(state.Running)
}

// InstallState returns information about whether or not a service is installed.
func (s *Service) InstallState() *state.Install {
	return state.NewInstall(s)
}

// RunState returns information about the runtime status of a service.
func (s *Service) RunState() *state.Run {
	return state.NewRun(s)
}

// PollUntil polls the service until it reaches the desired state or the timeout is reached.
func (s *Service) PollUntil(desired state.RunState, timeout time.Duration) (current *state.Run, timedOut bool) {
	current = s.RunState()
	fmt.Print(current.Pretty())
	deadline := time.Now().Add(timeout) // Don't wait forever
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for now := range ticker.C {
		current = s.RunState()
		fmt.Print("\r" + current.Pretty() + " ")
		if current.Is(desired) {
			break
		} else if now.After(deadline) {
			timedOut = true
			break
		}
	}
	return
}
