package site

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/acme/autocert"
)

func (b *Blog) RunAutoTLS() error {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(b.Config.Web.DomainNames[0:]...),
	}

	dir := cacheDir()
	fmt.Println("Using cache: ", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Printf("warning: autocert.NewListener not using a cache: %v", err)
	} else {
		m.Cache = autocert.DirCache(dir)
	}

	// Here we run the low priority .well-known and redirect to https handler in a goroutine
	// and "runWithManager" the actual web application.
	go http.ListenAndServe(b.Config.Web.ListenAddress+":80", m.HTTPHandler(nil))
	return runWithManager(b.Gin, *m, b.Config.Web.ListenAddress)
}

func runWithManager(r http.Handler, m autocert.Manager, address string) error {
	s := &http.Server{
		Addr:      address + ":443",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   r,
	}

	return s.ListenAndServeTLS("", "")
}

func cacheDir() string {
	const base = "golang-autocert"
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Caches", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(v, base)
			}
		}
		// Worst case:
		return filepath.Join(homeDir(), base)
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(homeDir(), ".cache", base)
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "/"
}
