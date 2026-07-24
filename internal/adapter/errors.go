package adapter

import "errors"

// ErrToolNotFound is returned when the OSS tool binary is not present in PATH.
// The scanner uses this sentinel to distinguish "tool absent" from "tool ran with 0 findings".
var ErrToolNotFound = errors.New("tool not in PATH")
