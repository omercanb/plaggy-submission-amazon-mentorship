package main

import (
	"aiplag-agent/common/config"
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
)

func main() {
	svcConfig := &service.Config{
		Name:        "com.plaggy.daemon",
		DisplayName: "Plaggy Daemon",
		Description: "Daemon for monitoring filesystem changes and assignments.",
		Executable:  config.DaemonExecutablePath(),
		Option: map[string]interface{}{
			"StartType": "automatic",
			"RunAtLoad": true,
			// Windows is the only OS where StartType: "automatic" actually auto-starts without extra commands, after the service is installed.
			// linux and macos require some shell commands of their own for launch on boot, could work on this later
			// but it seems like that will be done when we actually have installer stuff instead of dev commands
		},
	}

	d := &Daemon{} // <- type from daemon.go
	s, err := service.New(d, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err := service.Control(s, os.Args[1])
		if err != nil {
			logger.Errorf("Valid actions: install, uninstall, start, stop")
			log.Fatal(err)
		}
		fmt.Println("Service action executed:", os.Args[1])
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
