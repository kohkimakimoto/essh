package zssh

import(
	"fmt"
	"errors"
)

type Macro struct {
	Name        string
	Description string
	Parallel    bool
	Command     string                                               `gluamapper:"-"`
	CommandFunc func(payload string, host *Host) (string, error)     `gluamapper:"-"`
	Confirm     bool                                                 `gluamapper:"-"`
	ConfirmText string                                               `gluamapper:"-"`
	OnServers   []string                                             `gluamapper:"-"`
	OnTags      map[string][]string                                  `gluamapper:"-"`
	RunLocally  bool                                                 `gluamapper:"-"`
}

var Macros []*Macro = []*Macro{}

func GetMacro(name string) (*Macro, error) {
	for _, macro := range Macros {
		if macro.Name == name {
			return macro, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("not found '%s' task.", name))
}

func (m *Macro) Run(payload string) error {

	return nil
}