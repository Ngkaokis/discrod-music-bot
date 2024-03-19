package util

import "regexp"

var(
	UrlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	SearchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
)

