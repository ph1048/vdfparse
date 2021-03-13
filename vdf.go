package vdfparse

import (
	"errors"
	"strconv"
	"strings"
	"text/scanner"
)

type VdfNode map[string]interface{}

func parseNode(s *scanner.Scanner, first bool) (*VdfNode, error) {
	result := &VdfNode{}
	expectKey := true
	curKey := ""

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tt := s.TokenText()
		switch tt {
		case "{":
			if expectKey {
				return nil, errors.New("{ unexpected at line " + strconv.Itoa(s.Position.Line))
			}
			nod, err := parseNode(s, false)
			if err != nil {
				return nil, err
			}
			(*result)[curKey] = nod
			expectKey = true
		case "}":
			return result, nil
		default:
			str, err := strconv.Unquote(tt)
			if err == nil {
				tt = str
			}
			if expectKey {
				curKey = tt
				expectKey = false
			} else {
				(*result)[curKey] = tt
				expectKey = true
			}
		}
	}
	if first {
		return result, nil
	}
	return nil, errors.New("Unexpected EOF")
}

func ParseVdf(vdf string) (*VdfNode, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(vdf))
	s.Mode = scanner.ScanComments | scanner.SkipComments | scanner.ScanStrings | scanner.ScanIdents
	return parseNode(&s, true)
}
