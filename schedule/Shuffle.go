package schedule

import (
    "time"
    "math/rand"
)

func Shuffle(slice interface{}) interface{}{
    switch slice.(type) {
    case []*Group:
        result := slice.([]*Group)
        rand.Seed(int64(time.Now().Nanosecond()))
        for i := len(result)-1; i > 0; i-- {
            j := rand.Intn(i+1)
            result[i], result[j] = result[j], result[i]
        }
        return result
    case []*Subject:
        result := slice.([]*Subject)
        rand.Seed(int64(time.Now().Nanosecond()))
        for i := len(result)-1; i > 0; i-- {
            j := rand.Intn(i+1)
            result[i], result[j] = result[j], result[i]
        }
        return result
    }
    return nil
}
