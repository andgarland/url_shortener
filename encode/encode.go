package encode

import (
    "strings"
)

const (
    chars   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    size    = int64(len(chars))
)

//Encode a database id # to get its respective key (the shortURL). Key generated by translating the id to it's alphanumeric equivalent.
func Encode(id int64) string {
    if id == 0 {
        return string(chars[0])
    }

    //slight adjustment to break up the purely incremental sequence
    id = id * 11

    key := ""

    for id > 0 {
        key = string(chars[id % size]) + key
        id = id / size
    }

    return key
}

//Reverse the process to decode the shortURL back to the id #
func Decode(key string) int64 {
    var id int64

    for i := 0; i < len(key); i++ {
        id = id * size + int64(strings.Index(chars, string(key[i])))
    }

    id = id / 11

    return id
}