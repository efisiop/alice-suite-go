package useragent

import "strings"

// DeviceLabel returns a short human-readable label for the given User-Agent string,
// e.g. "Chrome on macOS", "Safari on iPhone". Returns "Unknown device" if empty or unparseable.
func DeviceLabel(ua string) string {
	ua = strings.TrimSpace(ua)
	if ua == "" {
		return "Unknown device"
	}
	u := strings.ToLower(ua)

	// Detect OS / device (order matters for mobile)
	var os string
	switch {
	case strings.Contains(u, "iphone"):
		os = "iPhone"
	case strings.Contains(u, "ipad"):
		os = "iPad"
	case strings.Contains(u, "android"):
		os = "Android"
	case strings.Contains(u, "mac os x") || strings.Contains(u, "macintosh"):
		os = "macOS"
	case strings.Contains(u, "windows"):
		os = "Windows"
	case strings.Contains(u, "linux"):
		os = "Linux"
	default:
		os = "Unknown"
	}

	// Detect browser (order matters: Edge/Opera first; Safari before Chrome because Safari UA can contain "Chrome")
	var browser string
	switch {
	case strings.Contains(u, "edg/"):
		browser = "Edge"
	case strings.Contains(u, "opr/") || strings.Contains(u, "opera"):
		browser = "Opera"
	case strings.Contains(u, "firefox"):
		browser = "Firefox"
	case (strings.Contains(u, "safari") || strings.Contains(u, "applewebkit")) && !strings.Contains(u, "chrome/") && !strings.Contains(u, "crios"):
		browser = "Safari"
	case strings.Contains(u, "chrome") || strings.Contains(u, "crios"):
		browser = "Chrome"
	default:
		browser = "Browser"
	}

	if os == "Unknown" {
		return browser
	}
	return browser + " on " + os
}
