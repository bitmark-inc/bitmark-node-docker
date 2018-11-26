package services

type Service interface {
	Initialise(configFile string) error
	Finalise() error
	IsRunning() bool
	SetNetwork(string)
	Status() map[string]interface{}
	Start() error
	Stop() error
	GetPath() string
	GetNetwork() string
}
