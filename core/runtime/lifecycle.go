package runtime

type Lifecycle string

const (
	Created Lifecycle = "Created"
	Running Lifecycle = "Running"
	Stopped Lifecycle = "Stopped"
)