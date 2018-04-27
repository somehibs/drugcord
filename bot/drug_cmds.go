package drugcord

import (
	"fmt"
	"strings"

	"github.com/somehibs/tripapi/api"
)

var DrugCommands = map[string]Command{"drug": DrugCmd{DiscordFormatter{}}}

type DrugCmd struct {
	Formatter Formatter
}

func (d DrugCmd) Matches(command *MessageInput) bool {
	return strings.HasPrefix(command.Content, "drug")
}

func (d DrugCmd) Action(command *MessageInput) (response []CommandResponse) {
	if command == nil {
		panic("Cannot see command?")
	}
	drugName := command.Split[1]
	drug := tripapi.GetDrug(drugName)

	reply := []string{}
	if drug != nil && d.Formatter != nil {
		// Format the drug with a nonexistent formatter.
		reply = append(reply, "`drug` "+drug.PrettyName)
		reply = append(reply, d.Formatter.FormatAll(drug))
	} else if d.Formatter == nil {
		reply = []string{"Formatter doesn't exist (error)."}
	} else {
		reply = []string{fmt.Sprintf("Could not find drug %s", drugName)}
	}
	if len(reply) > 0 {
		response = append(response, CommandResponse{command, reply, TargetSameChannel})
	}

	return response
}
