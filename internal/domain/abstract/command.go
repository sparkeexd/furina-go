package abstract

import "github.com/sparkeexd/furina/internal/domain/entity"

// Service's slash commands to be registered.
type CommandService interface {
	Commands() map[string]entity.Command
}
