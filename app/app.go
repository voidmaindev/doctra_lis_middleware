// Package app provides the interface that defines the methods that an application should implement.
package app

import "github.com/voidmaindev/doctra_lis_middleware/log"

const serverStopTimeout = 10

// App is the interface that defines the methods that an application should implement.
type App interface {
	SetLogger(*log.Logger)
	InitApp() error
	setConfig() error
	Start() error
	Stop() error
}
