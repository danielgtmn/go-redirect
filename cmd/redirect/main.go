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

func main() {
	currentLogLevel = parseLogLevel(os.Getenv("LOG_LEVEL"))

	targetURL := os.Getenv("REDIRECT_TARGET")
	if targetURL == "" {
		log.Fatal("REDIRECT_TARGET environment variable is required")
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

		target := targetURL
		if preservePath {
			target = targetURL + r.URL.Path
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
		}

		logInfo("%s %s %s -> %s (%d)", anonIP, r.Method, r.URL.Path, target, redirectCode)
		http.Redirect(w, r, target, redirectCode)
	})

	logInfo("Starting redirect server on :%s", port)
	logInfo("Target: %s (Code: %d, Preserve Path: %v, Block Scanners: %v)", targetURL, redirectCode, preservePath, blockScanners)
	logDebug("Log level: %s", os.Getenv("LOG_LEVEL"))
	if blockScanners && len(blockedPaths) > 0 {
		logInfo("Blocking %d path patterns", len(blockedPaths))
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
