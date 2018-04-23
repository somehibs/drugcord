package drugcord

import (
	"strings"
)

// Types and structs for routing bot commands.
const TargetNone = 0         // Only log the Reply (if there is one).
const TargetRequestor = 1    // Only respond to the user who requested this by PM.
const TargetSameChannel = 2  // Only respond to the channel this message was posted in.
const TargetAdminChannel = 3 // Only tell the admins about this.
const TargetOtherChannel = 4 // Target a specific channel, specified out of band.

type CommandHandler interface {
	Send(response CommandResponse)
	SendAll(response []CommandResponse)
}

type CommandResponse struct {
	Reply  []string
	Target int32
}

type GlobalCommand interface {
	Action(command string) []CommandResponse
}

type Command interface {
	GlobalCommand
	Matches(command string) bool
}

type CommandRouter struct {
	globals  []GlobalCommand
	commands map[string]Command
	handler  *CommandHandler
}

func (cr *CommandRouter) Init(handler *CommandHandler) {
	cr.handler = handler
}

func (cr *CommandRouter) RegisterCommands(commands map[string]Command) {
	for k, v := range commands {
		cr.commands[k] = v
	}
}

func (cr *CommandRouter) RegisterGlobals(globals []Command) {
	for _, x := range globals {
		alreadyExists := false
		for _, g := range cr.globals {
			if x == g {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			cr.globals = append(cr.globals, x)
		}
	}
}

// We expect that you'll have stripped any protocol spaces and jargon so we can parse some plain text.
func (cr *CommandRouter) HandleMessage(message string) []CommandResponse {
	split := strings.Split(message, " ")
	cmd := split[0]
	remainder := message
	if len(split) > 1 {
		remainder = strings.Join(split[1:], " ")
	}
	// Check if this message matches a plugin command. This is a super cheap lookup, so why not try it.
	command, ok := cr.commands[cmd]
	if ok {
		return command.Action(remainder)
	}
	return []CommandResponse{}
}
