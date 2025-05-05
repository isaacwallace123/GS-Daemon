package core

type Installer interface {
	ShouldInstall(egg Egg) bool
	Install(egg Egg) error
}
