package tfparser

import (
	"fmt"
	"strings"
)

type state int

const (
	stateTop             state = iota // top level context
	stateModuleName                   // read token 'module', await for module name
	stateModuleOpenBlock              // read module name (stored it in parser), await for open curly brace
	stateModule                       // inside module block, but not parsing any lexems in it

	stateModuleSource      // read token 'source' for module. await '=' sign
	stateModuleSourceValue // await for value of 'source' parameter of module

	stateModuleProvidersDeclared  // read token 'providers' passed to module. await for '=' sign
	stateModuleProvidersOpenBlock // read token '=' for providers block. await for '{' (not using this one!!!!)
	stateModuleProviders          // inside 'providers' block. await for provider alias
	stateModuleProiderAlias       // read provider alias passed into module. Wait for '=' sign
	stateModuleProiderName        // read provider name corresponded to alias.

	stateModuleParameterName  // parameter passed to module read. Await for '='
	stateModuleParameterValue // read '=' after reading parameter name. Await for parameter value
)

// parser for stateTopf state
// it skips/ignores all tokens/blocks except for 'module'
func (p *parser) parseTopLevel() error {
	switch p.state {
	case stateTop:
		switch strings.ToLower(p.peek()) {
		case "module":
			p.state = stateModuleName
			p.pop()
		case "{":
			p.skipBlock()
		// default is a token we are not parsing right now.
		default:
			p.pop()
		}
	case stateModuleName: // next token must be moudule name
		if p.curModName != "" {
			p.err = fmt.Errorf("FSM error, reading module name, while we have current module %#q", p.curModName)
			return p.err
		}
		if p.config.Modules == nil {
			p.config.Modules = make(map[string]*Module)
		}
		p.curModName = p.pop()
		_, exists := p.config.Modules[p.curModName]
		if exists {
			p.err = fmt.Errorf("Duplicated module name found: %#q", p.curModName)
			return p.err
		}
		p.config.Modules[p.curModName] = &Module{make(map[string]string), make(map[string]string), ""}
		p.state = stateModuleOpenBlock
	case stateModuleOpenBlock:
		p.err = p.popToken("{")
		if p.err != nil {
			return p.err
		}
		p.state = stateModule
	default:
		return p.parseModule()
	}
	return nil
}

// parser for module level context
func (p *parser) parseModule() error {
	switch p.state {

	// This is top level of module context
	case stateModule:
		switch p.peek() {
		case "source":
			p.pop()
			p.state = stateModuleSource
		case "providers":
			p.pop()
			p.state = stateModuleProvidersDeclared
		case "}":
			p.pop()
			p.curModName = ""
			p.state = stateTop
		default: // this is module parameters
			if p.curModParName != "" {
				p.err = fmt.Errorf("FSM error, next token is to be parameter, but there is already parameter %#q", p.curModParName)
				return p.err
			}
			p.curModParName = p.pop()
			p.state = stateModuleParameterName
		}

	// Following 2 items are dealing with source
	case stateModuleSource:
		p.popToken("=")
		p.state = stateModuleSourceValue
	case stateModuleSourceValue:
		p.state = stateModule
		if p.curModName == "" {
			p.err = fmt.Errorf("FSM error, mod name expected to be set")
			return p.err
		}
		_, exists := p.config.Modules[p.curModName]
		if !exists {
			p.err = fmt.Errorf("Module object for %#q was not created", p.curModName)
		}
		p.config.Modules[p.curModName].SourcePath = p.pop()
		//mod.SourcePath = p.pop()

	// Following is to do with providers
	case stateModuleProvidersDeclared:
		p.popToken("=")
		p.popToken("{")
		p.state = stateModuleProviders
	case stateModuleProviders, stateModuleProiderAlias, stateModuleProiderName:
		return p.parseProviders()

	// Read parameter
	case stateModuleParameterName:
		p.popToken("=")
		parValue := p.pop()
		// add parameters
		if p.curModParName == "" {
			p.err = fmt.Errorf("FSM error: got param value %#q, but param name is empty", parValue)
			return p.err
		}
		_, exists := p.config.Modules[p.curModName].Parameters[p.curModParName]
		if exists {
			p.err = fmt.Errorf("Duplicated parameter %#q", p.curModParName)
			return p.err
		}
		p.config.Modules[p.curModName].Parameters[p.curModParName] = parValue
		p.state = stateModule
		p.curModParName = ""
	}
	return nil
}

func (p *parser) parseProviders() error {
	switch p.state {
	case stateModuleProviders:
		tok := p.pop()
		switch tok {
		case "}":
			p.state = stateModule
			return nil
		default:
			p.popToken("=")
			provName := p.pop()
			// we are not changing state, we do not want to parse each simple 'alias = value' strings using FSM
			_, exists := p.config.Modules[p.curModName].Providers[tok]
			if exists {
				p.err = fmt.Errorf("Provier alias %#q has already been used in module %#q", tok, p.curModName)
				return p.err
			}
			p.config.Modules[p.curModName].Providers[tok] = provName
		}
	}
	return nil
}
