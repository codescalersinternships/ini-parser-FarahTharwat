package pkg

import (
	"INI_PARSER/pkg/utils"
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type IniParser struct {
	Sections map[string]map[string]string
}

func NewIniParser() *IniParser {
	return &IniParser{
		Sections: make(map[string]map[string]string),
	}
}
func (p * IniParser)Parse(scanner *bufio.Scanner , src string) (error){
	section := ""
	bracket_pattern := `^\[.+\]$`
	pair_pattern:= `^\s*[a-zA-Z][a-zA-Z0-9_]*[0-9]?\s*=\s*[^"]*$`
	bracket_regex := regexp.MustCompile(bracket_pattern)
	pair_regex := regexp.MustCompile(pair_pattern)
	for scanner.Scan() {
		line:=scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || len(line)==0 {
			continue
		}else if strings.HasPrefix(line, "[") {
			if !bracket_regex.MatchString(line) {
				return fmt.Errorf("invalid format at line: '%s', incorrect syntax for sections", line)
			}
			section=strings.Trim(line, "[]")
			p.Sections[section]=make(map[string]string)
		} else if !pair_regex.MatchString(line) {
			return fmt.Errorf("invalid format at line: '%s', key-value pairs must be in the format key=value ", line)
		}else {
			pair:=strings.Split(line, "=")
			if len(pair)==2 {
				key:=strings.TrimSpace(pair[0])
				value:=strings.TrimSpace(pair[1])
				p.Sections[section][key]=value
			}
		}
	}
	return nil
}
func (p *IniParser)LoadFromString(text string) error {
	scanner:=bufio.NewScanner(strings.NewReader(text))
	err:= p.Parse(scanner, text)
	if err != nil {
		return err
	}
	return nil
}

func (p *IniParser)LoadFromFile(path string) error {
	file,err := utils.OpenFile(path)
	if err != nil {
		return err}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	err = p.Parse(scanner, path)
	if err != nil {
			return err}
	return nil
}

func (p *IniParser)Get(section string, key string) (string, error) {
	if val, ok := p.Sections[section][key]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}
func (p *IniParser)GetSectionsNames()([]string) {
		var sections []string
		for section := range p.Sections {
			sections = append(sections, section)
		}
		return sections
}
func (p *IniParser)GetSections()(map[string]map[string]string) {
	if p.Sections == nil {
		return nil
	}
	return p.Sections
}
func (p *IniParser)Set(section string,key string,value string) (error) {
	if _, ok := p.Sections[section]; ok {
		p.Sections[section][key] = value
		return nil
	}
	return errors.New("section not found")
}
func (p *IniParser)ToString()(string) {
	var text string
	if p.Sections == nil {
		return " "
	}
	for section, pairs := range p.Sections {
		text += "[" + section + "]\n"
		for key, value := range pairs {
			text += key + "=" + value + "\n"
		}
	}
	return text
}
func (p *IniParser)SaveToFile(filename string,path string)(error) {
	file, err := utils.CreateFile(path,filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(p.ToString())
	if err != nil {
		return err
	}
	return nil
}

