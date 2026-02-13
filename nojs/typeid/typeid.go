//go:build js || wasm
// +build js wasm

package typeid

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

// GenerateTypeID creates a deterministic uint32 TypeID for a component type.
// It uses MD5 hash of the package path and component name to ensure
// consistency across builds and collision resistance.
func GenerateTypeID(packagePath, componentName string) uint32 {
	// Hash the fully-qualified component name
	h := md5.Sum([]byte(fmt.Sprintf("%s.%s", packagePath, componentName)))

	// Convert first 4 bytes to uint32
	return binary.BigEndian.Uint32(h[:4])
}
