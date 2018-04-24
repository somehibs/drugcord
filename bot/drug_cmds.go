package drugcord

import (
	"fmt"
	"github.com/somehibs/tripapi/api"
	"strings"
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

	if drug != nil && d.Formatter != nil {
		// Format the drug with a nonexistent formatter.
		response = append(response, CommandResponse{command, []string{d.Formatter.FormatAll(drug)}, TargetSameChannel})
	} else if d.Formatter == nil {
		response = append(response, CommandResponse{command, []string{"Formatter doesn't exist (error)."}, TargetSameChannel})
	} else {
		response = append(response, CommandResponse{command, []string{fmt.Sprintf("Could not find drug %s", drugName)}, TargetSameChannel})
	}

	return response
}
