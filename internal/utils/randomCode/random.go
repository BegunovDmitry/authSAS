package utils_random

import "math/rand/v2"

func RandRange(min, max int) int {
    return rand.IntN(max-min) + min
}