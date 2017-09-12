package services

type Service interface {
	Initialise(configFile string) error
	Finalise() error
	IsRunning() bool
	SetNetwork(string)
	Status() string
	Start() error
	Stop() error
}
