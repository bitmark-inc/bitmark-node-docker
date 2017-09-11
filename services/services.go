package services

type Service interface {
	Initialise(configFile string) error
	Finalise() error
	IsRunning() bool
	Status() string
	Start() error
	Stop() error
}
