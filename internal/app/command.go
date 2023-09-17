package app

import (
	"regexp"
)

type DataType int

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	Undefined DataType = iota
	Number
	Title
	Select
)

type ReturnCommand struct {
	Fields map[string]Field
}

type Field struct {
	DataType DataType
	Value    interface{}
}

type Command struct {
	command string
	regex   string
}

type CommandHandler struct {
	commands []*Command
}

func NewCommand(c string, r string) *Command {
	return &Command{
		command: c,
		regex:   r,
	}
}

func NewCommandHandler(commands []*Command) CommandHandler {
	return CommandHandler{
		commands: commands,
	}
}
func (ch *CommandHandler) WhichCommand(s string) *Command {
	for _, e := range ch.commands {
		if e.checkIsCommand(s) {
			return e
		}
	}
	return NewCommand("unknown", "")
}

func (c *Command) checkIsCommand(s string) bool {
	reg := regexp.MustCompile(c.regex)
	subMatchs := reg.FindStringSubmatch(s)
	return len(subMatchs) > 2
}

// interface หรือป่าว
func (c *Command) extract(s string, srv *Service, f func(s []string, srv *Service) map[string]Field) ReturnCommand {
	reg := regexp.MustCompile(c.regex)
	subMatchs := reg.FindStringSubmatch(s)
	resp := f(subMatchs, srv)
	return ReturnCommand{resp}
}
