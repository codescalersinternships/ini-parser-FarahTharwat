package main
import (
	"INI_PARSER/pkg"
	"fmt"
)
func main(){
	p := pkg.NewIniParser()
	err := p.LoadFromString("[section1]\nhello=valu(e1()\nkey2=value2\n[section2]\nkey3=value3\nkey4=value")
	if err != nil {
		panic(err)
	}
	sections := p.GetSectionsNames()
	fmt.Printf("Sections: %v\n", sections)
	section1, err := p.Get("section1", "hello")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Section1: %v\n", section1)
	section2, err := p.Get("section2", "key3")
	if err != nil {
		panic(err)
	}		
	fmt.Printf("Section2: %v\n", section2)
	err = p.Set("section1", "hello", "new value")
	if err != nil {
		panic(err)
	}
	section1, err = p.Get("section1", "hello")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Section1: %v\n", section1)
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
	//Change The Path to the location of the file
	err = p.LoadFromFile("H:/Me/CodeScalers/Golang/INI_Parser/example.INI")
	if err != nil {
		panic(err)
	}
	sections = p.GetSectionsNames()
	fmt.Printf("Sections: %v\n", sections)
	fmt.Printf("Section1: %v\n", section1)
    // getsections = p.GetSections()
    // if getsections == nil {
    //     fmt.Println("No sections found.")
    // } else {
    //     for section, kv := range getsections {
    //         fmt.Println("Section:", section)
    //         for key, value := range kv {
    //             fmt.Printf("%s = %s\n", key, value)
    //         }
    //     }
    // }
	text := p.ToString()
	fmt.Printf("Text: %s\n", text)
	}