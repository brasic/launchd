package launchd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderPlistWithBlankService(t *testing.T) {
	svc := &Service{}
	assertRenders(t, svc, `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string></string>
    <key>ProgramArguments</key>
    <array>
      <string></string>
    </array>
    <key>RunAtLoad</key>
    <false/>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>
`)
}

func TestRenderPlistWithRealService(t *testing.T) {
	svc := mkService()
	assertRenders(t, svc, `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.example.testservice</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/testservice</string>
      <string>-config</string>
      <string>/etc/testservice.conf</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    
    <key>KeepAlive</key>
    <dict>
      <key>SuccessfulExit</key>
      <false/>
    </dict>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>
`)
}

func TestRenderPlistNoKeepalive(t *testing.T) {
	svc := mkService()
	svc.KeepAlive = false
	assertRenders(t, svc, `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.example.testservice</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/testservice</string>
      <string>-config</string>
      <string>/etc/testservice.conf</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>
`)
}

func TestRenderPlistNoRunAtLoad(t *testing.T) {
	svc := mkService()
	svc.RunAtLoad = false
	assertRenders(t, svc, `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.example.testservice</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/testservice</string>
      <string>-config</string>
      <string>/etc/testservice.conf</string>
    </array>
    <key>RunAtLoad</key>
    <false/>
    
    <key>KeepAlive</key>
    <dict>
      <key>SuccessfulExit</key>
      <false/>
    </dict>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>
`)
}

func mkService() *Service {
	return &Service{
		Name:           "com.example.testservice",
		ExecutablePath: "/usr/local/bin/testservice",
		Argv:           []string{"-config", "/etc/testservice.conf"},
		KeepAlive:      true,
		RunAtLoad:      true,
	}
}

func assertRenders(t *testing.T, svc *Service, expected string) {
	result, err := svc.RenderPlist()
	require.NoError(t, err)
	require.Equal(t, expected, string(result))
}
