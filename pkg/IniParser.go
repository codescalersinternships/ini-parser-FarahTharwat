package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

/*    Error constants used all over the codebase    */
const (
	ErrOpeningFile   = "no such file or directory"
	ErrMatchingPairs = "key-value pairs must be in the format key=value"
	ErrSectionSyntax = "incorrect syntax for sections"
	ErrEmptyMap      = "map is empty !"
)

/*
		A struct that defines/represents an INI Parser . It has sections which stores each section
		with its corresponding key-pair values
*/
type IniParser struct {
	sections map[string]map[string]string
}

/* 		A function to initialise the Parser's sections map    */
func NewIniParser() *IniParser {
	return &IniParser{
		sections: make(map[string]map[string]string),
	}
}

/*
		This function parses each line in the text and checks the syntax of key-value pairs,sections...
		This function returns if an error occured but also transforms the INI string into an INI Parser instance
*/
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

/*
		This function takes an INI string format as an input and returns if an error occurs. It also uses the parse function
	    to fill at the end the Ini sections map.
*/
func (p *IniParser) LoadFromString(text string) error {
	scanner := bufio.NewScanner(strings.NewReader(text))
	err := p.parse(scanner)
	if err != nil {
		return err
	}
	return nil
}

/*
		This function takes a file path as an input and returns if an error occurs (ex: file does not exist, wrong Ini file format..).
	     It also uses the LoadFromString after converting the content of the file into string.
*/
func (p *IniParser) LoadFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return p.LoadFromString(string(content))
}

/*
		This function takes the section name and the key as an input and returns if an error occurs (key does not exist..).
	    If no erros , it retruns the value corresponding to the given key and section
*/
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

/* 		This function returns a string slice containg the names of every existing section */
func (p *IniParser) GetSectionNames() (sectionNames []string) {
	var sections []string
	for section := range p.sections {
		sections = append(sections, section)
	}
	return sections
}

/* 		This function returns a string slice containg the names of every existing section */
func (p *IniParser) GetSections() map[string]map[string]string {
	return p.sections
}

/*
		This function sets a key-value pair to the given section.
	    If section did not exist, a new section will be created
*/
func (p *IniParser) Set(section string, key string, value string) {
	if _, ok := p.sections[section]; !ok {
		log.Printf("warning : section %s did not exit however it has been created", section)
		p.sections[section] = make(map[string]string)
	}
	p.sections[section][key] = value
}

/*
		This function implements the Stringer interface as to be able to use fmt.Println(p) normally .
	    It prints the Ini Parser struct into string
*/
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

/* 		This function takes a file path of a file and saves the Ini Parser interface into string in this file */
func (p *IniParser) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(p.String())
	return err
}
