package config

import _ "embed"

// Embedded redirect configuration
//
//go:embed redirect.json
var Data []byte
