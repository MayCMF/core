package boot

import (
	"github.com/MayCMF/core/src/common/config"
	"github.com/google/gops/agent"
)

// InitMonitor - Initialize service monitoring
func InitMonitor() error {
	if c := config.Global().Monitor; c.Enable {
		err := agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir, ShutdownCleanup: true})
		if err != nil {
			return err
		}
	}
	return nil
}
