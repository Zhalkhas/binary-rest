package repos

import "fmt"

var ErrIndexNotFound = fmt.Errorf("index not found")
var ErrInvalidIndicesFile = fmt.Errorf("invalid indices file")
var ErrParsingIndicesFile = fmt.Errorf("error parsing indices file")
