# ini-parser-FarahTharwat

## How to use 

### What to install first

1. Install Go. I am using version `1.22.4`. You can use this link as a guide: [install-go-in-vscode](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code#1-install-go)
2. Clone the repo.

### How to run 

- If you want to run the tests of the APIs before using them, navigate to the `pkg` directory in your terminal using `cd pkg` and then run `go test -v`.
- If you want to use the APIs in your code, follow this sequence:
    1. First, create an instance of INI Parser using the API `NewIniParser()`.
    2. Then, choose between loading the INI content from a file or as a string (from a variable) using `LoadFromString(text string)` or `LoadFromFile(path string)`. Ensure you provide a valid file path.
    3. After successful parsing, you can:
        - Add key-pair values to a certain section using `Set(section string, key string, value string)`.
        - Get all section names using `GetSectionNames()`.
        - Get the value of a key in a certain section using `Get(section string, key string)`.
        - Get all the sections with their corresponding key-value pairs using `GetSections()`.
        - Print your implemented INI Parser interface through `String()` by calling `fmt.Println(p)` or `p.String()`, where `p` is your implemented instance of your INI interface.
        - Save your instance to a given file path using `SaveToFile(path string)`.
