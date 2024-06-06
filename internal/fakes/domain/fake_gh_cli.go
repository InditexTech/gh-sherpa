package domain

type FakeGhCli struct{}

func NewFakeGhCli() FakeGhCli {
	return FakeGhCli{}
}

func (f FakeGhCli) Execute(result any, args []string) (err error) {
	return nil
}
