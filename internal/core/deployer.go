package core

type Deployer interface {
	Run(egg Egg) error
	Stop(containerID string) error
}
