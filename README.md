# ini-parser-FarahTharwat
## How to use 
 > What to install first <br> <br>
[go-install] : https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code#1-install-go 
  1. install golang . I am using version ==1.22.4==. You can use this link as a guide [install go in vscode][go-install] <br> 
  2. clone the repo into your vs code <br> 

> How to run <br> <br>
  - If you want to run the testing of the apis before using them just write in your terminal ``` cd pkg ``` then hit ``` go go test -v ```
  - If you want to use the apis in your code then you should use the api in this sequence :
      > 1- First create a instance of INI Parser using the api `NewIniParser()` <br>
        2- Then you can choose between loading the Ini Content from a file or as string (from a variable..) using `LoadFromString(text string)` or `LoadFromFile(path string)` but take care you should give a valid  file path <br>
        3- After successful parsing , you can then choose between :
          >  - adding key-pair values to a certain section using `Set(section string, key string, value string)`
      >      - getting all sections names using `GetSectionNames()`
      >      - getting a certain value of a key in a certain section through `Get(section string, key string)`
      >      - getting all the sections with its corresponding key-value pairs using `GetSections()`
      >      - printing your implemented INI Parser interface through `String()` by calling ``` go fmt.Println(p)``` or ``` go p.String() ``` where p is your implemented instance of your INI interface
      >      - saving your instance into a given file path using `SaveToFile(path string)`

