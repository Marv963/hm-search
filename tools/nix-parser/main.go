package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Declaration repräsentiert eine einzelne Deklaration
type Declaration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Declarations ist ein Slice mit Custom Unmarshaler für flexible Formate
type Declarations []Declaration

func (d *Declarations) UnmarshalJSON(data []byte) error {
	// Fall 1: null
	if string(data) == "null" {
		*d = nil
		return nil
	}

	// Fall 2: Einzelner String (z.B. "<nixpkgs/...>")
	var singleStr string
	if err := json.Unmarshal(data, &singleStr); err == nil {
		*d = []Declaration{{Name: singleStr, URL: ""}}
		return nil
	}

	// Fall 3: Array von Strings
	var strSlice []string
	if err := json.Unmarshal(data, &strSlice); err == nil {
		decls := make([]Declaration, len(strSlice))
		for i, s := range strSlice {
			decls[i] = Declaration{Name: s, URL: ""}
		}
		*d = decls
		return nil
	}

	// Fall 4: Array von Objekten (Standardfall)
	var objSlice []Declaration
	if err := json.Unmarshal(data, &objSlice); err == nil {
		*d = objSlice
		return nil
	}

	// Fallback: leeres Array
	*d = []Declaration{}
	return nil
}

// RawOption aus Nix
type RawOption struct {
	Declarations Declarations    `json:"declarations"`
	Description  string          `json:"description"`
	Default      json.RawMessage `json:"default"`
	Example      json.RawMessage `json:"example"`
	Loc          []string        `json:"loc"`
	ReadOnly     bool            `json:"readOnly"`
	Type         string          `json:"type"`
}

// Option for Astro.js
type Option struct {
	Name         string       `json:"name"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Type         string       `json:"type"`
	Default      string       `json:"default"`
	Example      string       `json:"example"`
	ReadOnly     bool         `json:"readOnly"`
	Loc          []string     `json:"loc"`
	Declarations Declarations `json:"declarations"`
}

type Output struct {
	LastUpdate string   `json:"last_update"`
	Count      int      `json:"count"`
	Options    []Option `json:"options"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input.json> [output.json]")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := "options-parsed.json"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	// 1. Einlesen
	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var rawOptions map[string]RawOption
	if err := json.Unmarshal(data, &rawOptions); err != nil {
		panic(err)
	}

	// 2. Transformieren
	var options []Option

	for name, raw := range rawOptions {
		// Filter: _module.args und andere interne _module Einträge überspringen
		if name == "_module.args" || (len(name) > 8 && name[:8] == "_module.") {
			continue
		}

		opt := Option{
			Name:         name,
			Title:        name,
			Description:  raw.Description,
			Type:         raw.Type,
			Default:      resolveValue(raw.Default),
			Example:      resolveValue(raw.Example),
			ReadOnly:     raw.ReadOnly,
			Loc:          raw.Loc,
			Declarations: raw.Declarations,
		}
		options = append(options, opt)
	}

	// Sortieren nach Name (für bessere UX in Astro)
	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})

	// 3. Output bauen
	output := Output{
		LastUpdate: time.Now().UTC().Format("January 2, 2006 at 15:04 UTC"),
		Count:      len(options),
		Options:    options,
	}

	// 4. Schreiben
	out, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(outputFile, out, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Parsed %d options to %s\n", len(options), outputFile)
}

// resolveValue wandelt Nix literalExpression in plain text um
func resolveValue(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	// Versuch als literalExpression zu parsen
	var lit struct {
		Type string `json:"_type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &lit); err == nil && lit.Type == "literalExpression" {
		return lit.Text
	}

	// Versuch als plain string
	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		return str
	}

	// Versuch als andere primitive Typen (number, bool)
	var iface any
	if err := json.Unmarshal(raw, &iface); err == nil {
		return fmt.Sprintf("%v", iface)
	}

	// Fallback: raw JSON als String
	return string(raw)
}
