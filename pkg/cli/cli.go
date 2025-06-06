package cli

import (
	"github.com/spf13/cobra"
)

type cobraCmd = *cobra.Command
type cobraArgs = []string

// OkFlag

type OkFlag struct {
	Name      string
	ShortName string // Optional
	ValueType OkFlagType
	Usage     string
	Default   []any
}

type OkFlagType string

const (
	OkFlagString OkFlagType = "string"
	OkFlagInt    OkFlagType = "int"
)

func NewOkFlag(flagType OkFlagType, name, shortName, usage string, defaultVal ...any) OkFlag {
	return OkFlag{
		Name:      name,
		ShortName: shortName,
		ValueType: flagType,
		Usage:     usage,
		Default:   defaultVal,
	}
}

// OkCmd

type OkCmd struct {
	Use           string
	Short         string
	Long          string
	Run           func(cmd cobraCmd, args cobraArgs)
	Flags         []OkFlag
	Aliases       []string
	RequiredFlags []string
	Args          cobra.PositionalArgs
}

type OkParentCmd = OkCmd

// CLI

func Init() cobraCmd {
	return RegisterCommands(RootCommandList{})
}

type CommandList interface {
	GetParentCommand() *OkParentCmd
	GetCommands() []*OkCmd
	GetSubcommands() []CommandList
}

// Go through commands recursively and build tree of commands
func RegisterCommands(commandList CommandList) cobraCmd {
	if commandList == nil {
		return nil
	}

	var parentCommand = cobraCmdFromOkCmd(commandList.GetParentCommand())

	for _, command := range commandList.GetCommands() {
		parentCommand.AddCommand(cobraCmdFromOkCmd(command))
	}

	for _, subcommand := range commandList.GetSubcommands() {
		parentCommand.AddCommand(RegisterCommands(subcommand))
	}

	return parentCommand
}

func cobraCmdFromOkCmd(command *OkCmd) cobraCmd {
	cmd := &cobra.Command{
		Use:     command.Use,
		Short:   command.Short,
		Long:    command.Long,
		Run:     command.Run,
		Aliases: command.Aliases,
		Args:    command.Args,
	}

	if command.Flags != nil {
		for _, flag := range command.Flags {
			switch flag.ValueType {
			case "string":
				defaultVal := ""
				if len(flag.Default) > 0 {
					defaultVal = flag.Default[0].(string)
				}
				cmd.Flags().StringP(flag.Name, flag.ShortName, defaultVal, flag.Usage)
			case "int":
				defaultVal := 0
				if len(flag.Default) > 0 {
					defaultVal = flag.Default[0].(int)
				}
				cmd.Flags().IntP(flag.Name, flag.ShortName, defaultVal, flag.Usage)
			}
		}
	}

	if len(command.RequiredFlags) > 0 {
		cmd.MarkFlagsRequiredTogether(command.RequiredFlags...)
		for _, f := range command.RequiredFlags {
			err := cmd.MarkFlagRequired(f)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	return cmd
}
