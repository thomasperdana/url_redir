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

	// Try to load from embedded data first
	err := json.Unmarshal(config.Data, &routes)
	if err != nil {
		fmt.Printf("INIT: Error unmarshaling embedded data: %v\n", err)
		// Fallback to hardcoded configuration
		loadFallbackConfig()
	} else {
		fmt.Printf("INIT: Successfully loaded embedded data\n")
	}

	// Check if we have the required routes
	hasValidConfig := len(routes) > 0 && hasAboutMeRoute()
	fmt.Printf("INIT: Valid config check: routes=%d, hasAboutMe=%v\n", len(routes), hasAboutMeRoute())

	// If no routes loaded or missing our specific route, use fallback
	if !hasValidConfig {
		fmt.Printf("INIT: Using fallback configuration (missing required routes)\n")
		loadFallbackConfig()
	}

	fmt.Printf("INIT: Final loaded %d routes\n", len(routes))
	for i, route := range routes {
		fmt.Printf("INIT: Route %d - Scope: %s, Rules: %v\n", i, route.Scope, route.Rules)
	}
}

func hasAboutMeRoute() bool {
	for _, route := range routes {
		if route.Scope == "about\\.me" {
			// Check if we have both root path and thomas.perdana rules
			hasRoot := false
			hasThomas := false

			if _, exists := route.Rules["/"]; exists {
				hasRoot = true
			}
			if _, exists := route.Rules[""]; exists {
				hasRoot = true
			}
			if _, exists := route.Rules["thomas.perdana"]; exists {
				hasThomas = true
			}

			return hasRoot && hasThomas
		}
	}
	return false
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
				// Try root rule first
				if url, exists := rt.Rules["/"]; exists && url != "" {
					fmt.Printf("Root redirect to: %s\n", url)
					http.Redirect(w, r, url, http.StatusFound)
					return
				}
				// Try empty path rule
				if url, exists := rt.Rules[""]; exists && url != "" {
					fmt.Printf("Empty path redirect to: %s\n", url)
					http.Redirect(w, r, url, http.StatusFound)
					return
				}
				// If no root rules found, show error for root path
				fmt.Printf("No rule found for root path\n")
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

			// If we reach here, no rule was found for this specific path
			fmt.Printf("No rule found for path: '%s'\n", pathKey)
			_, _ = fmt.Fprintf(w, "Invalid short name for path: '%s'", pathKey)
			return
		}
	}
	fmt.Printf("No matching scope for host: %s\n", r.Host)
	_, _ = fmt.Fprintf(w, "No routes configured for host: %s", r.Host)
}
