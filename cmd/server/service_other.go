// +build !windows

package main

import (
	"fmt"
)

// isWindowsService sempre retorna false em plataformas não-Windows
func isWindowsService() (bool, error) {
	return false, nil
}

// runAsWindowsService não é suportado em plataformas não-Windows
func runAsWindowsService() error {
	return fmt.Errorf("Windows Service não é suportado nesta plataforma")
}
