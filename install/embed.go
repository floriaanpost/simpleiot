package install

import "embed"

// Content is a FS that holds the web UI assets
//
//go:embed siot.service*
var Content embed.FS
