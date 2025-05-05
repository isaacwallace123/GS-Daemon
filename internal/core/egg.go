package core

type Egg interface {
	GetName() string
	GetStartupCommand() string
	GetEnvironment() map[string]string
	GetImage() string
	GetPorts() []int
	GetVolumes() []string
}
