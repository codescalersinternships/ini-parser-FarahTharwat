package pkg

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const (
	CorrectStringTestCase1 = `;
					[section1]
					hello=valu(e1()
					;
					key2=value2
					[section2]
					#key5=123
					key3=value3
					key4=value
								`
	CorrectStringTestCase2 = " "

	WrongStringTestCase1 = `[section1]
					hellovalu(e1()
					key2value2
								`
	CorrectStringTestCase15 = `[section1]
hello=valu(e1()
key2=value2
[section2]
key3=value3
key4=value
`
)

type TestCase struct {
	src  string
	want map[string]map[string]string
}

var OutputTestCase1 = map[string]map[string]string{
	"section1": {
		"hello": "valu(e1()",
		"key2":  "value2",
	},
	"section2": {
		"key3": "value3",
		"key4": "value",
	},
}
var EmptyMapOutput = map[string]map[string]string{}

var OutputTestCase5 = map[string]map[string]string{
	"owner": {
		"name":         "John Doe",
		"organization": "Acme Widgets Inc.",
	},
	"database": {
		"server": "192.0.2.62",
		"port":   "143",
		"file":   "payroll.dat",
	},
}

var OutputTestCase10 = []string{"owner", "database"}

func TestLoadFromString(t *testing.T) {
	t.Run("testing ini style strings with acceptable formats", func(t *testing.T) {
		tests := []TestCase{
			{src: CorrectStringTestCase1, want: OutputTestCase1},
			{src: CorrectStringTestCase2, want: EmptyMapOutput},
		}
		for index, test := range tests {
			t.Run(fmt.Sprintf("test #%v", index), func(t *testing.T) {
				p := NewIniParser()
				_ = p.LoadFromString(test.src)
				assertCorrectResult(t, p.sections, test.want)
			})
		}
	})
	t.Run("testing an ini style string with wrong format : test#3", func(t *testing.T) {
		p := NewIniParser()
		err := p.LoadFromString(WrongStringTestCase1)
		assertError(t, fmt.Sprint(err), ErrorMatchingPairs)
	})

}
func TestLoadFromFile(t *testing.T) {
	t.Run("testing not existing file testcase#4", func(t *testing.T) {
		test := TestCase{
			src: "/mnt/h/recovery-keys/tests.txt",
		}
		p := NewIniParser()
		err := p.LoadFromFile(test.src)
		assertError(t, fmt.Sprint(err), ErrorOpeningFile)

	})
	t.Run("testing existing file testcase#5", func(t *testing.T) {
		test := TestCase{
			src:  "../example.ini",
			want: OutputTestCase5,
		}
		p := NewIniParser()
		_ = p.LoadFromFile(test.src)
		assertCorrectResult(t, p.sections, test.want)

	})
}
func TestGet(t *testing.T) {
	t.Run("testing get function with existing key testcase#6", func(t *testing.T) {
		p := NewIniParser()
		_, err := p.Get("owner", "name")
		assertError(t, fmt.Sprint(err), ErrorEmptyMap)

	})
	t.Run("testing get function with existing key testcase#7", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		got, _ := p.Get("owner", "name")
		assertCorrectResult(t, got, "John Doe")

	})
	t.Run("testing get function with non-existing section testcase#8", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		_, err := p.Get("house", "gender")
		assertError(t, fmt.Sprint(err), "does not exist")
	})
	t.Run("testing get function with non-existing key testcase#9", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		_, err := p.Get("owner", "gender")
		assertError(t, fmt.Sprint(err), "was not found")
	})
}

func TestGetSectionNames(t *testing.T) {
	t.Run("testing section names testcase#10", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		sectionNames, _ := p.GetSectionNames()
		assertCorrectResult(t,sectionNames, OutputTestCase10)
	})
	t.Run("testing section names from empty map testcase#11", func(t *testing.T) {
		p := NewIniParser()
		_, err := p.GetSectionNames()
		assertError(t, fmt.Sprint(err), ErrorEmptyMap)
	})
}

func TestGetSections(t *testing.T) {
	t.Run("testing sections map testcase#12", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		sections, _ := p.GetSections()
		assertCorrectResult(t, sections, OutputTestCase5)
	})
	t.Run("testing sections map from empty map testcase#13", func(t *testing.T) {
		p := NewIniParser()
		_, err := p.GetSectionNames()
		assertError(t, fmt.Sprint(err), ErrorEmptyMap)
	})
}

func TestSet(t *testing.T) {
	t.Run("testing set value to key in map testcase#14", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		p.Set("owner", "gender", "female")
		assertCorrectResult(t, p.sections["owner"]["gender"], "female")
	})
	t.Run("testing set value to key for a section that did not exist before testcase#15", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase5
		p.Set("employee", "gender", "female")
		assertCorrectResult(t, p.sections["employee"]["gender"], "female")
	})
}

func TestString(t *testing.T) {
	t.Run("testing converting ini object to string testcase#16", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase1
		content := p.String()
		p.LoadFromString(content)
		assertCorrectResult(t,p.sections,OutputTestCase1)
	})

	t.Run("testing converting empty ini object to string testcase#17", func(t *testing.T) {
		p := NewIniParser()
		p.sections = EmptyMapOutput
		got := p.String()
		assertCorrectResult(t, got, "")
	})

}

func TestSaveToFile(t *testing.T) {
	t.Run("testing not existing file or directory to save into the ini object testcase#18", func(t *testing.T) {
		p := NewIniParser()
		p.sections = EmptyMapOutput
		err := p.SaveToFile("/mnt/h/recovery-keys/tests.txt")
		assertError(t, fmt.Sprint(err), ErrorOpeningFile)
	})
	t.Run("testing existing file to save into the ini object testcase#19", func(t *testing.T) {
		p := NewIniParser()
		p.sections = OutputTestCase1
		p.SaveToFile("/mnt/h/recovery-keys/test.txt")
		p2:= NewIniParser()
		p2.LoadFromFile("/mnt/h/recovery-keys/test.txt")
		sections,_:=p2.GetSections()
		assertCorrectResult(t,sections,p.sections)
	})
}

func assertCorrectResult(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q \n want %q", got, want)
	}
}

func assertError(t testing.TB, errGotten, want string) {
	t.Helper()
	if !strings.Contains(errGotten, want) {
		t.Errorf("got error %q ", errGotten)
	}
}
