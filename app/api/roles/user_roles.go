package roles

const (
	ADMIN      = "admin"
	INTERNAL   = "internal"
	WORKER     = "worker"
	FREELANCER = "freelancer"
	CLIENT     = "client"
	PLAYER     = "player"
)

func FetchWorkflowRoles() []string {
	return []string{ADMIN, INTERNAL}
}
