// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package did

import (
	"fmt"
	"strings"
)

func Parse(input string) (*Identifier, error) {
	p := &parser{input: input, out: &Identifier{
		data: &identifierData{},
	}}

	parserState := p.checkLength
	for parserState != nil {
		parserState = parserState()
	}

	if p.err != nil {
		return nil, fmt.Errorf("%v: %v", p.err, "invalid DID")
	}

	p.out.data.Path = strings.Join(p.out.data.PathSegments[:], "/")

	return p.out, nil
}

type parser struct {
	input        string
	currentIndex int
	out          *Identifier
	err          error
}

type parserStep func() parserStep

func (p *parser) checkLength() parserStep {
	if inputLength := len(p.input); inputLength < 7 {
		return p.errorf(inputLength, "input length is less than 7")
	}
	return p.parseScheme
}

func (p *parser) parseScheme() parserStep {
	currentIndex := 3
	if p.input[:currentIndex+1] != prefix {
		return p.errorf(currentIndex, "input does not begin with 'did:' prefix: %s", p.input[:currentIndex+1])
	}
	p.currentIndex = currentIndex
	return p.parseMethod
}

func (p *parser) parseMethod() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	for {
		if currentIndex == inputLength {
			return p.errorf(currentIndex, "input does not have a second `:` marking end of method name")
		}

		char := input[currentIndex]

		if char == ':' {
			if currentIndex == startIndex {
				return p.errorf(currentIndex, "method is empty")
			}
			break
		}

		currentIndex = currentIndex + 1
	}

	p.currentIndex = currentIndex
	p.out.data.Method = input[startIndex:currentIndex]

	return p.parseID
}

func (p *parser) parseID() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var next parserStep
	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]

		if char == ':' {
			next = p.parseID
			break
		}

		if char == ';' {
			next = p.parseParamName
			break
		}

		if char == '/' {
			next = p.parsePath
			break
		}

		if char == '?' {
			next = p.parseQuery
			break
		}

		if char == '#' {
			next = p.parseFragment
			break
		}

		currentIndex = currentIndex + 1
	}

	if currentIndex == startIndex {
		return p.errorf(currentIndex, "idstring must be at least one char long")
	}

	p.currentIndex = currentIndex
	p.out.data.ID = input[startIndex:currentIndex]

	return next
}

func (p *parser) parseParamName() parserStep {
	input := p.input
	startIndex := p.currentIndex + 1
	next := p.paramTransition()
	currentIndex := p.currentIndex

	if currentIndex == startIndex {
		return p.errorf(currentIndex, "Param name must be at least one char long")
	}

	p.out.data.Params = append(p.out.data.Params, Param{Name: input[startIndex:currentIndex], Value: ""})

	return next
}

func (p *parser) parseParamValue() parserStep {
	input := p.input
	startIndex := p.currentIndex + 1
	next := p.paramTransition()
	currentIndex := p.currentIndex

	p.out.data.Params[len(p.out.data.Params)-1].Value = input[startIndex:currentIndex]

	return next
}

func (p *parser) paramTransition() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1

	var indexIncrement int
	var next parserStep

	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]

		if char == ';' {
			next = p.parseParamName
			break
		}

		if char == '=' {
			next = p.parseParamValue
			break
		}

		if char == '/' {
			next = p.parsePath
			break
		}

		if char == '?' {
			next = p.parseQuery
			break
		}

		if char == '#' {
			next = p.parseFragment
			break
		}

		if char == '%' {
		} else {
			indexIncrement = 1
		}

		currentIndex = currentIndex + indexIncrement
	}

	p.currentIndex = currentIndex

	return next
}

func (p *parser) parsePath() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var next parserStep

	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]

		if char == '/' {
			next = p.parsePath
			break
		}

		if char == '?' {
			next = p.parseQuery
			break
		}

		if char == '%' {
		} else {
			indexIncrement = 1
		}

		currentIndex = currentIndex + indexIncrement
	}

	if currentIndex == startIndex && len(p.out.data.PathSegments) == 0 {
		return p.errorf(currentIndex, "first path segment must have at least one character")
	}

	p.currentIndex = currentIndex
	p.out.data.PathSegments = append(p.out.data.PathSegments, input[startIndex:currentIndex])
	return next
}

func (p *parser) parseQuery() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var next parserStep

	for {
		if currentIndex == inputLength {
			break
		}

		char := input[currentIndex]

		if char == '#' {
			next = p.parseFragment
			break
		}

		if char == '%' {
		} else {
			indexIncrement = 1
		}

		currentIndex = currentIndex + indexIncrement
	}

	p.currentIndex = currentIndex
	p.out.data.Query = input[startIndex:currentIndex]
	return next
}

func (p *parser) parseFragment() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int

	for {
		if currentIndex == inputLength {
			break
		}

		char := input[currentIndex]

		if char == '%' {
		} else {
			indexIncrement = 1
		}

		currentIndex = currentIndex + indexIncrement
	}

	p.currentIndex = currentIndex
	p.out.data.Fragment = input[startIndex:currentIndex]

	return nil
}

func (p *parser) errorf(index int, format string, args ...interface{}) parserStep {
	p.currentIndex = index
	p.err = fmt.Errorf(format, args...)
	return nil
}
