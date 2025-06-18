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
	fmt.Printf("INIT: Raw embedded data: %s\n", string(config.Data))
	err := json.Unmarshal(config.Data, &routes)
	if err != nil {
		fmt.Printf("INIT: Error unmarshaling: %v\n", err)
	}
	fmt.Printf("INIT: Loaded %d routes\n", len(routes))
	for i, route := range routes {
		fmt.Printf("INIT: Route %d - Scope: %s, Rules: %v\n", i, route.Scope, route.Rules)
	}
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	// Debug info for troubleshooting
	fmt.Printf("DEBUG: Host=%s, Path=%s, Routes=%d\n", r.Host, r.URL.Path, len(routes))

	for i, rt := range routes {
		fmt.Printf("DEBUG: Route %d - Scope=%s, Rules=%v\n", i, rt.Scope, rt.Rules)
		scope := regexp.MustCompile(rt.Scope)
		if scope.MatchString(r.Host) {
			fmt.Printf("DEBUG: Host '%s' matches scope '%s'\n", r.Host, rt.Scope)

			// Handle root path
			if r.URL.Path == "" || r.URL.Path == "/" {
				if rt.Rules["/"] != "" {
					fmt.Printf("DEBUG: Root redirect to %s\n", rt.Rules["/"])
					http.Redirect(w, r, rt.Rules["/"], http.StatusFound)
					return
				}
			}

			// Handle specific paths
			pathKey := ""
			if len(r.URL.Path) > 1 {
				pathKey = r.URL.Path[1:]
			}
			fmt.Printf("DEBUG: Looking for path key '%s'\n", pathKey)

			if url := rt.Rules[pathKey]; url != "" {
				fmt.Printf("DEBUG: Redirecting '%s' to '%s'\n", pathKey, url)
				http.Redirect(w, r, url, http.StatusFound)
				return
			} else {
				fmt.Printf("DEBUG: No rule found for path '%s'\n", pathKey)
				_, _ = fmt.Fprintf(w, "Invalid short name for path: '%s'", pathKey)
				return
			}
		} else {
			fmt.Printf("DEBUG: Host '%s' does not match scope '%s'\n", r.Host, rt.Scope)
		}
	}
	fmt.Printf("DEBUG: No routes matched for host '%s'\n", r.Host)
	_, _ = fmt.Fprintf(w, "No routes configured for host: %s", r.Host)
}
