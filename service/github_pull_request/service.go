package githubpullrequest

type Action string

const (
	ActionOpened   Action = "opened"
	ActionSync     Action = "synchronize"
	ActionAssigned Action = "assigned"
)

func (a Action) ShouldValidate() bool {
	return a == ActionOpened || a == ActionSync
}

type Status string

const (
	StatusError   Status = "error"
	StatusFailure Status = "failure"
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
)
