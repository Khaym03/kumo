package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"syscall"

	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
)

type instanceData struct {
	launcher *launcher.Launcher
	wsURL    string
}

type MultiManager struct {
	*launcher.Manager

	instances map[string]*instanceData
	mu        sync.Mutex
}

func NewMultiManager() *MultiManager {
	return &MultiManager{
		Manager:   launcher.NewManager(),
		instances: make(map[string]*instanceData),
	}
}

func (m *MultiManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/list" {
		m.listInstances(w, r)
		return
	}
	if r.URL.Path == "/close" {
		m.closeInstance(w, r)
		return
	}

	if r.Header.Get("Upgrade") == "websocket" {
		instanceID := r.URL.Query().Get("id")
		if instanceID == "" {
			http.Error(w, "missing instance id", http.StatusBadRequest)
			return
		}
		m.mu.Lock()
		inst, ok := m.instances[instanceID]

		if !ok || !processIsRunning(inst.launcher.PID()) {
			l := launcher.New() // create new launcher
			// apply unique user data directory or other flags based on instanceID
			l.Set(flags.UserDataDir, fmt.Sprintf("/tmp/rod/user-data-%s", instanceID))

			l.Set(flags.NoSandbox)

			host := r.URL.Query().Get("proxyHost")
			port := r.URL.Query().Get("proxyPort")
			if host != "" && port != "" {
				l.Set(flags.ProxyServer, fmt.Sprintf("%s:%s", host, port))
			}

			// call BeforeLaunch hook if set
			if m.BeforeLaunch != nil {
				m.BeforeLaunch(l, w, r)
			}

			// launch browser and get WebSocket debug URL
			u := l.MustLaunch()

			m.instances[instanceID] = &instanceData{
				launcher: l,
				wsURL:    u,
			}
			m.mu.Unlock()

			// proxy websocket between client and launched browser
			m.proxyWebsocket(w, r, u)
			return
		}
		m.mu.Unlock()

		// proxy websocket to existing instance's URL
		m.proxyWebsocket(w, r, inst.wsURL)
		return
	}

	// default route - return general info or 400
	http.Error(w, "unsupported endpoint", http.StatusNotFound)
}

func (m *MultiManager) proxyWebsocket(w http.ResponseWriter, r *http.Request, target string) {
	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputil.NewSingleHostReverseProxy(toHTTP(*targetURL)).ServeHTTP(w, r)
}

func (m *MultiManager) listInstances(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := make([]string, 0, len(m.instances))
	for id := range m.instances {
		ids = append(ids, id)
	}
	b, _ := json.Marshal(ids)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (m *MultiManager) closeInstance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	m.mu.Lock()
	defer m.mu.Unlock()

	inst, ok := m.instances[id]
	if !ok {
		http.Error(w, "instance id not found", http.StatusNotFound)
		return
	}
	inst.launcher.Kill()
	inst.launcher.Cleanup()
	delete(m.instances, id)
	w.Write([]byte("closed"))
}

func toHTTP(u url.URL) *url.URL {
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}

	return &u
}

func processIsRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
