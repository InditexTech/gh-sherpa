package domain

type ForkConfiguration struct {
	DefaultOrganization string
	IsInteractive       bool
}

type ForkSetupResult struct {
	WasAlreadyConfigured bool
	ForkCreated          bool
	ForkName             string
	UpstreamName         string
}

type ForkStatus struct {
	IsInFork          bool
	HasCorrectRemotes bool
	HasCorrectDefault bool
	ForkName          string
	UpstreamName      string
}

type ForkProvider interface {
	IsRepositoryFork() (bool, error)
	CreateFork(forkName string) error
	ForkExists(forkName string) (bool, error)
	SetDefaultRepository(repo string) error
	GetRemoteConfiguration() (map[string]string, error)
	ConfigureRemotesForExistingFork(forkName string) error
}
