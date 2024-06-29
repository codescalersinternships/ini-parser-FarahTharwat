package pkg

import (
	"testing"
	
)
const (
	stringTestPass1 = `[section1]
					hello=valu(e1()
					key2=value2
					[section2]
					key3=value3
					key4=value
								`
	stringTestFail1 = `[section1]
					hellovalu(e1()
					key2value2
								`
	stringTestFail2 = " "
)

func TestLoadFromString(t *testing.T) {
	p := NewIniParser()
	t.Run( "testing an ini style string with a correct format" , func(t *testing.T){
		err:= p.LoadFromString(stringTestPass1)
		assertNotError(t,err)
	})
	
}

func assertNotError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got error %q want nil", got)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}