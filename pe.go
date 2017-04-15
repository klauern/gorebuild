package main

import (
	"debug/pe"
	"debug/gosym"
	"errors"
	"os"
)



func getPeTable(file string) (*gosym.Table, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var textStart uint64
	var symtab, pclntab []byte

	obj, err := pe.NewFile(f)
	if err == nil {
		if sect := obj.Section(".text"); sect == nil {
			return nil, errors.New("empty .text")
		}
		if sect := obj.Section(".gosymtab"); sect != nil {
			if symtab, err = sect.Data(); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("empty .gosymtab")
		}
		if sect := obj.Section(".gopclntab"); sect != nil {
			if pclntab, err = sect.Data(); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("empty .gopclntab")
		}

	}
		pcln := gosym.NewLineTable(pclntab, textStart)
	return gosym.NewTable(symtab, pcln)
}

func getPeMainPath(file string) (string, error) {
	table, err := getPeTable(file)
	if err != nil {
		return "", err
	}
	path, _, _ := table.PCToLine(table.LookupFunc("main.main").Entry)
	return path, nil
}
