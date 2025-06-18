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
	_ = json.Unmarshal(config.Data, &routes)
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	// Debug info for troubleshooting
	fmt.Printf("DEBUG: Host=%s, Path=%s, Routes=%d\n", r.Host, r.URL.Path, len(routes))

	for i, rt := range routes {
		fmt.Printf("DEBUG: Route %d - Scope=%s, Rules=%v\n", i, rt.Scope, rt.Rules)
		scope := regexp.MustCompile(rt.Scope)
		if scope.MatchString(r.Host) {
			fmt.Printf("DEBUG: Host '%s' matches scope '%s'\n", r.Host, rt.Scope)
			if rt.Rules["/"] != "" && r.URL.Path == "" {
				fmt.Printf("DEBUG: Root redirect to %s\n", rt.Rules["/"])
				http.Redirect(w, r, rt.Rules["/"], http.StatusFound)
				return
			}
			pathKey := r.URL.Path[1:]
			fmt.Printf("DEBUG: Looking for path key '%s'\n", pathKey)
			if url := rt.Rules[pathKey]; url != "" {
				fmt.Printf("DEBUG: Redirecting '%s' to '%s'\n", pathKey, url)
				http.Redirect(w, r, url, http.StatusFound)
				return
			} else {
				fmt.Printf("DEBUG: No rule found for path '%s'\n", pathKey)
				_, _ = fmt.Fprintf(w, "Invalid short name for path: %s", pathKey)
				return
			}
		} else {
			fmt.Printf("DEBUG: Host '%s' does not match scope '%s'\n", r.Host, rt.Scope)
		}
	}
	fmt.Printf("DEBUG: No routes matched for host '%s'\n", r.Host)
	_, _ = fmt.Fprintf(w, "No routes configured for host: %s", r.Host)
}
