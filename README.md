# dotenv
Loads environment variables from .env for go projects

## Usage

```go
package main

import "github.com/Kaiser925/dotenv"

func main() {
	err := dotenv.Load()
	if err != nil {
		return 
	}
}
```