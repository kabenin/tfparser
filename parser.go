/*
Package tfparser provides a parser that allows getting a very specific piece of information from terraform
configuration (as text, file or directory):
for all modules used in the configuration it reads all parameters and providers passed into module. It also
reads source path for the module.

Data is returned as type TFConfig, which consists of map of types 'Module'
*/
package tfparser

import (
	"fmt"
	"log"
	"strings"
)

// Module represents a call to a module
type Module struct {
	Providers  map[string]string
	Parameters map[string]string
	SourcePath string
}

// TFconfig represents a tf configiration
type TFconfig struct {
	Modules map[string]*Module
}

type parser struct {
	data          string    // tf file(s) as string
	i             int       // index in data
	config        *TFconfig // TFconfig struct we are building
	state         state     // FSM state
	err           error
	curModName    string // name of the module we are parsing
	curModParName string // If we are parsing module parametes, what it name is
}

// ParseString parses a string with tf configurarion
func ParseString(s string) (*TFconfig, error) {
	return (&parser{strings.TrimSpace(s), 0, &TFconfig{}, stateTop, nil, "", ""}).parse()
}

func (p *parser) parse() (*TFconfig, error) {
	for {
		if p.i >= len(p.data) {
			err := p.validate()
			return p.config, err
		}
		err := p.parseTopLevel()
		if err != nil {
			p.logError()
			return nil, err
		}
	}
}

func (p *parser) validate() error {
	if p.curModName != "" {
		return fmt.Errorf("Did not find the closing curly brace when parsing module %v", p.curModName)
	}
	if p.curModParName != "" {
		return fmt.Errorf("Did not find the value for %v param of module %v", p.curModParName, p.curModName)
	}
	return nil
}

func (p *parser) logError() {
	if p.err == nil {
		return
	}
	log.Fatal(fmt.Sprintln(p.err))
}
