// +build windows

package main

import (
	"wsicrmrest/internal/service"

	"golang.org/x/sys/windows/svc"
)

// isWindowsService verifica se o processo está rodando como Windows Service
func isWindowsService() (bool, error) {
	return svc.IsWindowsService()
}

// runAsWindowsService executa a aplicação como Windows Service
func runAsWindowsService() error {
	return service.RunAsWindowsService()
}
