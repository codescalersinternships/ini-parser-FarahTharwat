package pkg
import (
	"bufio"
	"errors"
	"regexp"
	"strings"
	"INI_PARSER/pkg/utils"
)

type IniParser struct {
	Sections map[string]map[string]string
}

func NewIniParser() *IniParser {
	return &IniParser{
		Sections: make(map[string]map[string]string),
	}
}
func (p *IniParser)LoadFromString(text string) error {
	scanner:=bufio.NewScanner(strings.NewReader(text))
	section := ""
	bracket_pattern := `^\[.+\]$`
	pair_pattern := `^\s*[a-zA-Z][a-zA-Z0-9_]*[0-9]?\s*=\s*[^"]*$`
	//`^\s*[a-zA-Z\s]+\s*=\s*[^"]*$`
	// file : `^[a-zA-Z]+\s*=\s*(?:"[^"]*"|\d+)$`
	//`^[a-zA-Z]+="?[a-zA-Z][^=]*"?$`
	//^[a-zA-Z]+\s*=\s*"?[^"]*"?
	bracket_regex := regexp.MustCompile(bracket_pattern)
	pair_regex := regexp.MustCompile(pair_pattern)
	if text[0] != '[' {
		return errors.New("invalid format , sections must start with [")
	}
	for scanner.Scan() {
		line:=scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || len(line)==0 {
			continue
		}else if strings.HasPrefix(line, "[") {
			if !bracket_regex.MatchString(line) {
				return errors.New("invalid format , sections must end with ]")
			}
			section=strings.Trim(line, "[]")
			p.Sections[section]=make(map[string]string)
		} else if !pair_regex.MatchString(line) {
			return errors.New("invalid format , key value pairs must be in the format key=value")
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

func (p *IniParser)LoadFromFile(filename string) error {
	file,err := utils.OpenFile(filename)
	if err != nil {
		return err}
	defer file.Close()
	return err
	
}


func main (){
	p := NewIniParser()
	err := p.LoadFromString("[section1]\nhello=valu(e1()\nkey2=value2\n[section2]\nkey3=value3\nkey4=value")
	if err != nil {
		panic(err)
	}

}

