package drugcord

import (
	"fmt"
	"strings"
)

// Types and structs for routing bot commands.
type Target int32

const (
	TargetNone         Target = 0 // Only log the Reply (if there is one).
	TargetRequestor    Target = 1 // Only respond to the user who requested this by PM.
	TargetSameChannel  Target = 2 // Only respond to the channel this message was posted in.
	TargetAdminChannel Target = 3 // Only tell the admins about this.
	TargetOtherChannel Target = 4 // Target a specific channel, specified out of band.
)

// Implement this to receive responses from CommandRouter
type CommandHandler interface {
	Send(response CommandResponse)
	SendAll(response []CommandResponse)
}

type CommandRouter struct {
	globals  []GlobalCommand
	commands map[string]Command
}

type CommandResponse struct {
	Input  *MessageInput
	Reply  []string
	Target Target
}

// Command should be a slice of space separated words
// GlobalCommand will not prevent further processing unless the return boolean is true
type GlobalCommand interface {
	Action(command *MessageInput) ([]CommandResponse, bool)
}

type MessageInput struct {
	Original interface{}
	Content  string
	Split    []string
}

func (m MessageInput) SplitImpl() []string {
	return strings.Split(m.Content, " ")
}

// A non-zero length CommandResponse will prevent further command processing
type Command interface {
	Action(command *MessageInput) []CommandResponse
	Matches(command *MessageInput) bool
}

func (cr *CommandRouter) RegisterCommands(commands map[string]Command) {
	if cr.commands == nil {
		cr.commands = commands
		return
	}
	for k, v := range commands {
		cr.commands[k] = v
	}
}

func (cr *CommandRouter) RegisterGlobals(globals []GlobalCommand) {
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
func (cr *CommandRouter) HandleMessage(handler CommandHandler, message *MessageInput) {
	responses := cr.handleMessageImpl(message)
	if len(responses) > 0 {
		handler.SendAll(responses)
	} else {
		fmt.Printf("No responses for message %s\n", message.Content)
	}
}

func (cr *CommandRouter) handleMessageImpl(message *MessageInput) (response []CommandResponse) {
	response = []CommandResponse{}
	if len(cr.globals) == 0 && len(cr.commands) == 0 {
		fmt.Println("No command handlers registered. Empty response.")
		return response
	}
	// Initial vars
	message.Split = message.SplitImpl()
	cmd := message.Split[0]

	// Parse globals. Maybe exit early.
	for _, glob := range cr.globals {
		r, stop := glob.Action(message)
		response = append(response, r...)
		if stop {
			return response
		}
	}

	// Parse commands.
	// Try the map first.
	command, ok := cr.commands[cmd]
	if ok {
		return command.Action(message)
	}

	// Manually match against every command.
	for _, cmd := range cr.commands {
		if cmd.Matches(message) {
			response = append(response, cmd.Action(message)...)
			if len(response) > 0 {
				break
			}
		}
	}
	return response
}
