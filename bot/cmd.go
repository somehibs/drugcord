package drugcord

// Types and structs for routing commands.
const TargetNone = 0 // Only log the Reply (if there is one).
const TargetRequestor = 1 // Only respond to the user who requested this by PM.
const TargetSameChannel = 2 // Only respond to the channel this message was posted in.
const TargetAdminChannel = 3 // Only tell the admins about this.
const TargetOtherChannel = 4 // Target a specific channel, specified out of band.

type CommandResponse struct {
	Reply string
	Target int32
}

type Command interface {
	Matches(command string) bool
	Action(command string) []CommandResponse
}

type CommandRouter struct {
	global map[string]Command
	commands map[string]Command
}
