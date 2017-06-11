package cmder

import (
	"os"

	"github.com/go-mango/mango/logger"
)

//systemd script template.
const systemd = `[Unit]
Description=%s
After=network.target

[Service]
PrivateTmp=true
PIDFile=/var/run/%s.pid
ExecStart=/usr/local/bin/%s
ExecStop=kill -9 $MAINPID

[Install]
WantedBy=multi-user.target`

//Systemd systemd
type Systemd struct {
	Desc string
}

//SetDesc sets description of program.
func (sys *Systemd) SetDesc(s string) {
	sys.Desc = s
}

//Install installs program.
func (sys *Systemd) Install() {
	if os.Getuid() != 0 {
		logger.NewLogger().Fatal("installer must be run as root role.")
		os.Exit(1)
	}
}
