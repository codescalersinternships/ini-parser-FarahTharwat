package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	ErrOpeningFile   = "no such file or directory"
	ErrMatchingPairs = "key-value pairs must be in the format key=value"
	ErrSectionSyntax = "incorrect syntax for sections"
	ErrEmptyMap      = "map is empty !"
)

type IniParser struct {
	sections map[string]map[string]string
}

func NewIniParser() *IniParser {
	return &IniParser{
		sections: make(map[string]map[string]string),
	}
}
func (p *IniParser) parse(scanner *bufio.Scanner) error {
	section := ""
	bracketPattern := `^\[.+\]$`
	pairPattern := `^\s*[a-zA-Z][a-zA-Z0-9_]*[0-9]?\s*=\s*[^"]*$`
	bracketRegex := regexp.MustCompile(bracketPattern)
	pairRegex := regexp.MustCompile(pairPattern)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "[") && !bracketRegex.MatchString(line) {
			return fmt.Errorf("invalid format at line: '%s', %v", line, ErrSectionSyntax)
		}
		if strings.HasPrefix(line, "[") {
			section = strings.Trim(line, "[]")
			p.sections[section] = make(map[string]string)
			continue
		}
		if !pairRegex.MatchString(line) {
			return fmt.Errorf("invalid format at line: '%s', %v ", line, ErrMatchingPairs)
		}
		pair := strings.SplitN(line, "=", 2)
		if len(pair) == 2 {
			key := strings.TrimSpace(pair[0])
			value := strings.TrimSpace(pair[1])
			p.sections[section][key] = value
		}
	}
	return nil
}
func (p *IniParser) LoadFromString(text string) error {
	scanner := bufio.NewScanner(strings.NewReader(text))
	err := p.parse(scanner)
	if err != nil {
		return err
	}
	return nil
}
func (p *IniParser) LoadFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return p.LoadFromString(string(content))
}

func (p *IniParser) Get(section string, key string) (string, error) {
	if len(p.sections) == 0 {
		return " ", fmt.Errorf(ErrEmptyMap)
	}
	values, isMapContainsKey := p.sections[section]
	if !isMapContainsKey {
		return " ", fmt.Errorf("section %s does not exist", key)
	}
	if val, ok := values[key]; ok {
		return val, nil
	}
	return " ", fmt.Errorf("key %s was not found", key)
}
func (p *IniParser) GetSectionNames() (sectionNames []string) {
	var sections []string
	for section := range p.sections {
		sections = append(sections, section)
	}
	return sections
}
func (p *IniParser) GetSections() (map[string]map[string]string) {
	return p.sections
}

func (p *IniParser) Set(section string, key string, value string) {
	if _, ok := p.sections[section]; !ok {
		log.Printf("warning : section %s did not exit however it has been created", section)
		p.sections[section] = make(map[string]string)
	}
	p.sections[section][key] = value
}

func (p *IniParser) String() string {
	text := ""
	if p.sections == nil {
		return text
	}
	for section, pairs := range p.sections {
		text += "[" + section + "]\n"
		for key, value := range pairs {
			text += key + "=" + value + "\n"
		}
	}
	return text
}
func (p *IniParser) SaveToFile(path string) error {
	file, err := os.OpenFile(path,os.O_APPEND|os.O_WRONLY,0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(p.String())
	return err
}

