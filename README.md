# DuckDuckGo images api
This package allows you to get pictures from [duckduckgo](https://duckduckgo.com).

# Usage
## Import
```go
import ddg "github.com/beshenkaD/go-ddg-images-api"
```

## Use
```go
func main() {
        r, err := ddg.Do(ddg.Query{
                Keywords: "duck",
                Moderate: true,
        })

        if err != nil {
                panic(err)
        }

        for _, res := range r.Results {
                println(res.Image)
        }
}
```