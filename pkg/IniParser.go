package pkg

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type stringer interface {
	String() string
}
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
		} else if strings.HasPrefix(line, "[") {
			if !bracketRegex.MatchString(line) {
				return fmt.Errorf("invalid format at line: '%s', incorrect syntax for sections", line)
			}
			section = strings.Trim(line, "[]")
			p.sections[section] = make(map[string]string)
		} else if !pairRegex.MatchString(line) {
			return fmt.Errorf("invalid format at line: '%s', key-value pairs must be in the format key=value ", line)
		} else {
			pair := strings.SplitN(line, "=", 2)
			if len(pair) == 2 {
				key := strings.TrimSpace(pair[0])
				value := strings.TrimSpace(pair[1])
				p.sections[section][key] = value
			}
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
	if val, ok := p.sections[section][key]; ok {
		return val, nil
	}
	fmt.Printf("section %s does not exist \n", section)
	return " ", fmt.Errorf("key %s was not found", key)
}
func (p *IniParser) GetSectionNames() []string {
	var sections []string
	for section := range p.sections {
		sections = append(sections, section)
	}
	return sections
}
func (p *IniParser) GetSections() map[string]map[string]string {
	if p.sections == nil {
		return nil
	}
	return p.sections
}
func (p *IniParser) Set(section string, key string, value string) error {
	if _, ok := p.sections[section]; !ok {
		fmt.Printf("warning : section %s did not exit however it has been created", section)
	}
	p.sections[section][key] = value
	return nil
}
func (p *IniParser) String() string {
	text := " "
	if p.sections == nil {
		return " "
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
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(p.String())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (p *IniParser) print() {
	getsections := p.GetSections()
	if getsections == nil {
		fmt.Println("No sections found.")
	} else {
		for section, kv := range getsections {
			fmt.Println("Section:", section)
			for key, value := range kv {
				fmt.Printf("%s = %s\n", key, value)
			}
		}
	}
}
