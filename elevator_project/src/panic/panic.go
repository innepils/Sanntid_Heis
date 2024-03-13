package panic

import (
	"fmt"
)

func Recover() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
}
