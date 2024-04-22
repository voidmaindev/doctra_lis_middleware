package app

// App is the interface that defines the methods that an application should implement.
type App interface {
	InitApp() error
}
