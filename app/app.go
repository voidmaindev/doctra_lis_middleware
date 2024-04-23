// Package app provides the interface that defines the methods that an application should implement.
package app

import "github.com/voidmaindev/doctra_lis_middleware/log"

// App is the interface that defines the methods that an application should implement.
type App interface {
	SetLogger(*log.Logger)
	InitApp() error
}
