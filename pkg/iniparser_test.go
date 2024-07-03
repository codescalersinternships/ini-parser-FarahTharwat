package pkg

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const (
	validIniStringStartingWithComment = `;
	[section1]
	hello=valu(e1()
	;
	key2=value2
	[section2]
	#key5=123
	key3=value3
	key4=value
	`
	validIniStringEmpty     = " " // also test against ""
	invalidIniStringKeyPair = `[section1]
					hellovalu(e1()
					key2value2
					`
	validIniStringExampleWithIntegerValuePair = `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = payroll.dat`
)

type TestCase struct {
	src  string
	want map[string]map[string]string
}

var mapOutputValidIniStringStartingWithComment = map[string]map[string]string{
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

var mapOutputValidIniStringExampleWithIntegerValuePair = map[string]map[string]string{
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

var emptySlice = []string{}

func TestLoadFromString(t *testing.T) {
	t.Run("testing loading ini style strings with acceptable formats wrt their maps outputs tescases#1-2", func(t *testing.T) {
		tests := []TestCase{
			{src: validIniStringStartingWithComment, want: mapOutputValidIniStringStartingWithComment},
			{src: validIniStringEmpty, want: EmptyMapOutput},
		}
		for index, test := range tests {
			t.Run(fmt.Sprintf("test #%v", index), func(t *testing.T) {
				p := NewIniParser()
				_ = p.LoadFromString(test.src)
				assertEqual(t, p.sections, test.want)
			})
		}
	})
	t.Run("testing an ini style string with wrong format : test#3", func(t *testing.T) {
		p := NewIniParser()
		err := p.LoadFromString(invalidIniStringKeyPair)
		assertError(t, fmt.Sprint(err), ErrMatchingPairs)
	})

}
func TestLoadFromFile(t *testing.T) {
	t.Run("testing not existing file testcase#4", func(t *testing.T) {
		test := TestCase{
			src: "/mnt/h/recovery-keys/tests.txt",
		}
		p := NewIniParser()
		err := p.LoadFromFile(test.src)
		assertError(t, fmt.Sprint(err), ErrOpeningFile)

	})
	t.Run("testing existing file testcase#5", func(t *testing.T) {
		test := TestCase{
			src:  "../example.ini",
			want: mapOutputValidIniStringExampleWithIntegerValuePair,
		}
		p := NewIniParser()
		_ = p.LoadFromFile(test.src)
		assertEqual(t, p.sections, test.want)
	})
}
func TestGet(t *testing.T) {
	t.Run("testing get function with existing key testcase#6", func(t *testing.T) {
		p := NewIniParser()
		_, err := p.Get("owner", "name")
		assertError(t, fmt.Sprint(err), ErrEmptyMap)

	})
	t.Run("testing get function with existing key testcase#7", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		got, _ := p.Get("owner", "name")
		assertEqual(t, got, "John Doe")

	})
	t.Run("testing get function with non-existing section testcase#8", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		_, err := p.Get("house", "gender")
		assertError(t, fmt.Sprint(err), "does not exist")
	})
	t.Run("testing get function with non-existing key testcase#9", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		_, err := p.Get("owner", "gender")
		assertError(t, fmt.Sprint(err), "was not found")
	})
}

func TestGetSectionNames(t *testing.T) {
	t.Run("testing section names testcase#10", func(t *testing.T) {
		var outputMap = []string{"owner", "database"}
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		sectionNames := p.GetSectionNames()
		assertEqual(t, sectionNames, outputMap)
	})
	t.Run("testing section names from empty map testcase#11", func(t *testing.T) {
		p := NewIniParser()
		sectionNames := p.GetSectionNames()
		assertEqual(t, len(sectionNames), len(emptySlice))
	})
}

func TestGetSections(t *testing.T) {
	t.Run("testing sections map testcase#12", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		sections := p.GetSections()
		assertEqual(t, sections, mapOutputValidIniStringExampleWithIntegerValuePair)
	})
	t.Run("testing sections map from empty map testcase#13", func(t *testing.T) {
		p := NewIniParser()
		emptySectionNames := p.GetSections()
		assertEqual(t, emptySectionNames, EmptyMapOutput)
	})
}

func TestSet(t *testing.T) {
	t.Run("testing set value to key in map testcase#14", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		p.Set("owner", "gender", "female")
		assertEqual(t, p.sections["owner"]["gender"], "female")
	})
	t.Run("testing set value to key for a section that did not exist before testcase#15", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		p.Set("employee", "gender", "female")
		assertEqual(t, p.sections["employee"]["gender"], "female")
	})
}

func TestString(t *testing.T) {
	t.Run(fmt.Sprint("testing converting correct ini object into string testcases#16-17"), func(t *testing.T) {
		tests := []TestCase{
			{src: validIniStringStartingWithComment, want: mapOutputValidIniStringStartingWithComment},
			{src: "", want: EmptyMapOutput},
		}
		for index, test := range tests {
			t.Run(fmt.Sprintf("test #%v", index), func(t *testing.T) {
				p := NewIniParser()
				_ = p.LoadFromString(test.src)
				newout := p.String()
				p2 := NewIniParser()
				_ = p2.LoadFromString(newout)
				//asserting that the state of p1 and p2 are the same
				assertEqual(t, test.want, p2.GetSections())
			})
		}
	})
}
func TestSaveToFile(t *testing.T) {
	t.Run("testing not existing file or directory to save into the ini object testcase#18", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringEmpty)
		err := p.SaveToFile("/mnt/h/recovery-keys/tests.txt")
		assertError(t, fmt.Sprint(err), ErrOpeningFile)
	})
	t.Run("testing existing file to save into the ini object testcase#19", func(t *testing.T) {
		
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringStartingWithComment)
		_ = p.SaveToFile("../test.txt")
		p2 := NewIniParser()
		_ = p2.LoadFromFile("../test.txt")
		sections := p2.GetSections()
		assertEqual(t, sections, p.sections)
	})
}

func TestStringer(t *testing.T) {
	t.Run("testing stringer interface testcase#20", func(t *testing.T) {
		p := NewIniParser()
		_ = p.LoadFromString(validIniStringExampleWithIntegerValuePair)
		pContent := fmt.Sprint(p)
		p2 := NewIniParser()
		p2.LoadFromString(pContent)
		//asserting that the state of p1 and p2 are the same
		assertEqual(t, mapOutputValidIniStringExampleWithIntegerValuePair, p2.GetSections())
	})
}

func assertEqual(t *testing.T, got, want any) {
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
