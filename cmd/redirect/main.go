package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelNone
)

var currentLogLevel LogLevel

func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "none", "off":
		return LogLevelNone
	default:
		return LogLevelInfo
	}
}

func logDebug(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func logInfo(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelInfo {
		log.Printf("[INFO] "+format, v...)
	}
}

func logWarn(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelWarn {
		log.Printf("[WARN] "+format, v...)
	}
}

func anonymizeIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "-"
	}

	if ip.To4() != nil {
		return "x.x.x.x"
	}

	return "x:x:x:x:x:x:x:x"
}

func loadBlockedPaths() []string {
	path := os.Getenv("BLOCKED_PATHS_FILE")
	if path == "" {
		path = "/scanner_paths.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logWarn("Could not load %s: %v", path, err)
		return nil
	}

	var paths []string
	if err := json.Unmarshal(data, &paths); err != nil {
		logWarn("Invalid JSON in %s: %v", path, err)
		return nil
	}

	return paths
}

func isBlockedPath(path string, blockedPaths []string) bool {
	lowerPath := strings.ToLower(path)
	for _, blocked := range blockedPaths {
		if strings.Contains(lowerPath, strings.ToLower(blocked)) {
			return true
		}
	}
	return false
}

func mapDomain(host, sourceDomain, targetDomain string) (string, bool) {
	host = strings.ToLower(host)
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	if host == sourceDomain {
		return targetDomain, true
	}

	suffix := "." + sourceDomain
	if strings.HasSuffix(host, suffix) {
		prefix := host[:len(host)-len(suffix)]
		return prefix + "." + targetDomain, true
	}

	return "", false
}

func main() {
	currentLogLevel = parseLogLevel(os.Getenv("LOG_LEVEL"))

	targetURL := os.Getenv("REDIRECT_TARGET")
	domainMap := os.Getenv("REDIRECT_DOMAIN_MAP")

	var sourceDomain, targetDomain string
	if domainMap != "" {
		parts := strings.SplitN(domainMap, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			log.Fatal("REDIRECT_DOMAIN_MAP must be in format 'source:target' (e.g. 'domain.com:domain.de')")
		}
		sourceDomain = strings.ToLower(parts[0])
		targetDomain = strings.ToLower(parts[1])
	}

	if targetURL == "" && domainMap == "" {
		log.Fatal("REDIRECT_TARGET or REDIRECT_DOMAIN_MAP environment variable is required")
	}

	scheme := os.Getenv("REDIRECT_SCHEME")
	if scheme == "" {
		scheme = "https"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redirectCode := 301
	if code := os.Getenv("REDIRECT_CODE"); code != "" {
		if c, err := strconv.Atoi(code); err == nil && (c == 301 || c == 302) {
			redirectCode = c
		}
	}

	preservePath := os.Getenv("PRESERVE_PATH") == "true"
	blockScanners := os.Getenv("BLOCK_SCANNERS") == "true"

	var blockedPaths []string
	if blockScanners {
		blockedPaths = loadBlockedPaths()
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		anonIP := anonymizeIP(r.RemoteAddr)

		if blockScanners && len(blockedPaths) > 0 && isBlockedPath(r.URL.Path, blockedPaths) {
			logInfo("%s %s %s -> BLOCKED", anonIP, r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}

		var target string

		if domainMap != "" {
			if newHost, ok := mapDomain(r.Host, sourceDomain, targetDomain); ok {
				target = scheme + "://" + newHost + r.URL.Path
				if r.URL.RawQuery != "" {
					target += "?" + r.URL.RawQuery
				}
			} else if targetURL != "" {
				target = targetURL
				if preservePath {
					target = targetURL + r.URL.Path
					if r.URL.RawQuery != "" {
						target += "?" + r.URL.RawQuery
					}
				}
			} else {
				logInfo("%s %s %s -> NOT FOUND (no mapping for host %s)", anonIP, r.Method, r.URL.Path, r.Host)
				http.NotFound(w, r)
				return
			}
		} else {
			target = targetURL
			if preservePath {
				target = targetURL + r.URL.Path
				if r.URL.RawQuery != "" {
					target += "?" + r.URL.RawQuery
				}
			}
		}

		logInfo("%s %s %s -> %s (%d)", anonIP, r.Method, r.URL.Path, target, redirectCode)
		http.Redirect(w, r, target, redirectCode)
	})

	logInfo("Starting redirect server on :%s", port)
	if domainMap != "" {
		logInfo("Domain mapping: *.%s -> *.%s (Scheme: %s, Code: %d, Block Scanners: %v)", sourceDomain, targetDomain, scheme, redirectCode, blockScanners)
	}
	if targetURL != "" {
		logInfo("Target: %s (Code: %d, Preserve Path: %v, Block Scanners: %v)", targetURL, redirectCode, preservePath, blockScanners)
	}
	logDebug("Log level: %s", os.Getenv("LOG_LEVEL"))
	if blockScanners && len(blockedPaths) > 0 {
		logInfo("Blocking %d path patterns", len(blockedPaths))
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
