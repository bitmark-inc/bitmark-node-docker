package services

type Service interface {
	Initialise(configFile string) error
	Finalise() error
	IsRunning() bool
	Start() error
	Stop() error
}
