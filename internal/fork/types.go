package fork

type Configuration struct {
	DefaultOrganization string
	IsInteractive bool
}

type ForkSetupResult struct {
	WasAlreadyConfigured bool
	ForkCreated bool
	ForkName string
	UpstreamName string
}

type ForkStatus struct {
	IsInFork bool
	HasCorrectRemotes bool
	HasCorrectDefault bool
	ForkName string
	UpstreamName string
}
