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
	// Try to load from embedded data first
	err := json.Unmarshal(config.Data, &routes)
	if err != nil {
		fmt.Printf("INIT: Error unmarshaling embedded data: %v\n", err)
		// Fallback to hardcoded configuration
		loadFallbackConfig()
	}

	// If no routes loaded or missing our specific route, use fallback
	if len(routes) == 0 || !hasAboutMeRoute() {
		fmt.Printf("INIT: Using fallback configuration\n")
		loadFallbackConfig()
	}

	fmt.Printf("INIT: Loaded %d routes\n", len(routes))
	for i, route := range routes {
		fmt.Printf("INIT: Route %d - Scope: %s, Rules: %v\n", i, route.Scope, route.Rules)
	}
}

func hasAboutMeRoute() bool {
	for _, route := range routes {
		if route.Scope == "about\\.me" {
			if _, exists := route.Rules["thomas.perdana"]; exists {
				return true
			}
		}
	}
	return false
}

func loadFallbackConfig() {
	routes = []routeData{
		{
			Scope: "about\\.me",
			Rules: map[string]string{
				"/":              "https://about.cashinblue.com",
				"thomas.perdana": "https://about.cashinblue.com",
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
	for _, rt := range routes {
		scope := regexp.MustCompile(rt.Scope)
		if scope.MatchString(r.Host) {
			// Handle root path
			if r.URL.Path == "" || r.URL.Path == "/" {
				if rt.Rules["/"] != "" {
					http.Redirect(w, r, rt.Rules["/"], http.StatusFound)
					return
				}
			}

			// Handle specific paths
			pathKey := ""
			if len(r.URL.Path) > 1 {
				pathKey = r.URL.Path[1:]
			}

			if url := rt.Rules[pathKey]; url != "" {
				http.Redirect(w, r, url, http.StatusFound)
				return
			} else {
				_, _ = fmt.Fprintf(w, "Invalid short name for path: '%s'", pathKey)
				return
			}
		}
	}
	_, _ = fmt.Fprintf(w, "No routes configured for host: %s", r.Host)
}
