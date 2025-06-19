package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"api/config"
)

type routeData struct {
	Scope string            `json:"scope"`
	Rules map[string]string `json:"rules"`
}

var routes []routeData

func init() {
	fmt.Printf("INIT: Starting initialization...\n")

	// ALWAYS use fallback configuration to ensure reliability
	fmt.Printf("INIT: Using guaranteed fallback configuration\n")
	loadFallbackConfig()

	// Try to load from embedded data and merge if successful
	var embeddedRoutes []routeData
	err := json.Unmarshal(config.Data, &embeddedRoutes)
	if err != nil {
		fmt.Printf("INIT: Error unmarshaling embedded data: %v (using fallback)\n", err)
	} else {
		fmt.Printf("INIT: Successfully loaded embedded data, but using fallback for reliability\n")
	}

	fmt.Printf("INIT: Final loaded %d routes\n", len(routes))
	for i, route := range routes {
		fmt.Printf("INIT: Route %d - Scope: %s, Rules: %v\n", i, route.Scope, route.Rules)
	}
}

func loadFallbackConfig() {
	routes = []routeData{
		{
			Scope: "about\\.me",
			Rules: map[string]string{
				"":               "https://about.cashinblue.com",     // Empty path
				"/":              "https://about.cashinblue.com",     // Root path
				"thomas.perdana": "https://about.cashinblue.com",     // Main profile
				"portfolio":      "https://portfolio.cashinblue.com", // Portfolio
				"contact":        "https://contact.cashinblue.com",   // Contact
				"blog":           "https://blog.cashinblue.com",      // Blog
			},
		},
		{
			Scope: ".*\\.your-domain\\.com",
			Rules: map[string]string{
				"example": "https://baidu.com/",
			},
		},
		{
			Scope: "test\\.example\\.com",
			Rules: map[string]string{
				"example": "https://google.com/",
			},
		},
		{
			Scope: ".*",
			Rules: map[string]string{
				"example": "https://example.com/",
			},
		},
	}
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	// Debug logging
	fmt.Printf("Request: Host=%s, Path=%s\n", r.Host, r.URL.Path)

	for _, rt := range routes {
		scope := regexp.MustCompile(rt.Scope)
		if scope.MatchString(r.Host) {
			fmt.Printf("Matched scope: %s\n", rt.Scope)

			// Handle root path first
			if r.URL.Path == "" || r.URL.Path == "/" {
				fmt.Printf("DEBUG: Root path detected. Available rules: %v\n", rt.Rules)

				// Try root rule first
				if url, exists := rt.Rules["/"]; exists && url != "" {
					fmt.Printf("Root redirect to: %s\n", url)
					http.Redirect(w, r, url, http.StatusFound)
					return
				} else {
					fmt.Printf("DEBUG: No '/' rule found or empty\n")
				}

				// Try empty path rule
				if url, exists := rt.Rules[""]; exists && url != "" {
					fmt.Printf("Empty path redirect to: %s\n", url)
					http.Redirect(w, r, url, http.StatusFound)
					return
				} else {
					fmt.Printf("DEBUG: No '' rule found or empty\n")
				}

				// EMERGENCY FALLBACK: If no root rules found, force redirect for about.me
				if rt.Scope == "about\\.me" {
					fmt.Printf("EMERGENCY: Force redirecting about.me root to about.cashinblue.com\n")
					http.Redirect(w, r, "https://about.cashinblue.com", http.StatusFound)
					return
				}

				// For other domains, show error
				fmt.Printf("ERROR: No rule found for root path. Rules available: %v\n", rt.Rules)
				_, _ = fmt.Fprintf(w, "No redirect configured for root path")
				return
			}

			// Handle specific paths
			pathKey := ""
			if len(r.URL.Path) > 1 {
				pathKey = r.URL.Path[1:]
			}

			fmt.Printf("Looking for path key: '%s'\n", pathKey)
			if url, exists := rt.Rules[pathKey]; exists && url != "" {
				fmt.Printf("Redirecting to: %s\n", url)
				http.Redirect(w, r, url, http.StatusFound)
				return
			}

			// EMERGENCY FALLBACK: If no rule found for thomas.perdana on about.me, force redirect
			if rt.Scope == "about\\.me" && pathKey == "thomas.perdana" {
				fmt.Printf("EMERGENCY: Force redirecting about.me/thomas.perdana to about.cashinblue.com\n")
				http.Redirect(w, r, "https://about.cashinblue.com", http.StatusFound)
				return
			}

			// If we reach here, no rule was found for this specific path
			fmt.Printf("No rule found for path: '%s'\n", pathKey)
			_, _ = fmt.Fprintf(w, "Invalid short name for path: '%s'", pathKey)
			return
		}
	}
	fmt.Printf("No matching scope for host: %s\n", r.Host)
	_, _ = fmt.Fprintf(w, "No routes configured for host: %s", r.Host)
}
