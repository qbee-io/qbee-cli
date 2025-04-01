package broker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.qbee.io/client"
)

type Service struct {
	client         *client.Client
	authToken      string
	remoteHost     string
	remotePort     string
	remoteProtocol string
	listenPort     string
}

const (
	defaultConnectionTTL  = 5 * time.Minute
	defaultDevicePort     = "80"
	defaultListenPort     = "8081"
	defaultRemoteProtocol = "http"
	defaultRemoteHost     = "localhost"
)

type Connection struct {
	Cancel        context.CancelFunc
	Port          int
	LastConnected time.Time
}

var connectionMutex sync.Mutex

var connections = make(map[string]Connection)

func NewService() *Service {

	token := os.Getenv("QBEE_TOKEN")

	remoteHost := defaultRemoteHost
	if os.Getenv("QBEE_REMOTE_HOST") != "" {
		remoteHost = os.Getenv("QBEE_REMOTE_HOST")
	}

	remotePort := defaultDevicePort
	if os.Getenv("QBEE_REMOTE_PORT") != "" {
		remotePort = os.Getenv("QBEE_REMOTE_PORT")
	}

	listenPort := defaultListenPort
	if os.Getenv("QBEE_LISTEN_PORT") != "" {
		listenPort = os.Getenv("QBEE_LISTEN_PORT")
	}

	remoteProtocol := defaultRemoteProtocol
	if os.Getenv("QBEE_REMOTE_PROTOCOL") != "" {
		remoteProtocol = os.Getenv("QBEE_REMOTE_PROTOCOL")
	}

	return &Service{
		authToken:      token,
		remoteHost:     remoteHost,
		remotePort:     remotePort,
		remoteProtocol: remoteProtocol,
		listenPort:     listenPort,
	}
}

func (s *Service) Start(ctx context.Context) error {

	// Sart a goroutine that re-authenticates the client every minute
	go s.reAuthenticateClient(ctx)
	// Start garbage collector
	go garbageCollector()

	router := chi.NewRouter()
	router.Use(s.AuthMiddleware)
	router.HandleFunc("/*", s.Proxy())

	if s.authToken == "" {
		log.Println("Warning: No authentication token provided. Device access will be open")
	}

	log.Println("Starting server on :8081. Press CTRL+C to stop it.")
	return http.ListenAndServe(fmt.Sprintf(":%s", s.listenPort), router)
}

func (s *Service) WithPort(port string) *Service {
	s.listenPort = port
	return s
}

func (s *Service) WithClient(client *client.Client) *Service {
	s.client = client
	return s
}

func (s *Service) WithAuthToken(token string) *Service {
	s.authToken = token
	return s
}

func (s *Service) WithRemotePort(port string) *Service {
	s.remotePort = port
	return s
}

func (s *Service) WithRemoteProtocol(protocol string) *Service {
	s.remoteProtocol = protocol
	return s
}

// AuthMiddleware is a middleware that checks for the presence of an auth token
func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If no token has been configured, bypass auth
		if s.authToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Look for auth header
		token := r.Header.Get("X-Qbee-Authorization")
		if token != s.authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Proxy forwards the request to the device
func (s *Service) Proxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		deviceId, err := resolveDeviceId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		devicePort := r.Header.Get("X-Qbee-Device-Port")
		if devicePort == "" {
			devicePort = s.remotePort
		}

		localPort, err := s.doPortForwarding(ctx, deviceId, devicePort)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		localUrl := fmt.Sprintf("http://localhost:%d", localPort)
		parsedUrl, err := url.Parse(localUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

		proxy.Director = func(req *http.Request) {
			req.Header = r.Header
			req.Host = parsedUrl.Host
			req.URL.Scheme = parsedUrl.Scheme
			req.URL.Host = parsedUrl.Host
			req.URL.Path = r.URL.Path
		}

		proxy.ServeHTTP(w, r)

	}
}

// doPortForwarding establishes a port forwarding connection using random local port to the device on a specific port
func (s *Service) doPortForwarding(ctx context.Context, deviceId, devicePort string) (int, error) {
	connectionId := fmt.Sprintf("%s:%s", deviceId, devicePort)
	if conn, ok := connections[connectionId]; ok {
		connectionMutex.Lock()
		defer connectionMutex.Unlock()
		conn.LastConnected = time.Now()
		connections[connectionId] = conn
		return conn.Port, nil
	}

	//log.Printf("Connecting to device %s on port %s\n", deviceId, devicePort)
	deviceNodeID, err := s.resolveDeviceID(ctx, deviceId)
	if err != nil {
		return 0, err
	}

	//log.Printf("Resolved device %s to node %s\n", deviceId, deviceNodeID)
	localPort, err := GetFreePort()
	if err != nil {
		return 0, err
	}

	remoteTarget := []client.RemoteAccessTarget{
		{
			RemotePort: devicePort,
			LocalPort:  fmt.Sprintf("%d", localPort),
			Protocol:   "tcp",
			LocalHost:  "localhost",
			RemoteHost: s.remoteHost,
		},
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	errChan := make(chan error)
	go func() {
		errChan <- s.client.Connect(cancelCtx, deviceNodeID, remoteTarget)
	}()

	portReady := make(chan bool)
	go waitForPort(localPort, portReady)

	select {
	case err := <-errChan:
		if err != nil {
			cancelFunc()
			return 0, err
		}
	case isPortReady := <-portReady:
		if !isPortReady {
			cancelFunc()
			return 0, fmt.Errorf("port %d is not ready", localPort)
		}
	}

	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	connections[connectionId] = Connection{
		Cancel:        cancelFunc,
		Port:          localPort,
		LastConnected: time.Now(),
	}
	// return the local port if the connection was successful
	return localPort, nil
}

// Re-authenticate the client every minute
const clientReAuthInterval = 10 * time.Minute

// reAuthenticateClient re-authenticates the client periodically
func (s *Service) reAuthenticateClient(ctx context.Context) {
	for {
		time.Sleep(clientReAuthInterval)
		// Should do
		//log.Printf("Re-authenticating client\n")
		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			log.Fatal(err)
		}
		s.client = cli
	}
}

const garbageCollectorInterval = 1 * time.Minute

// garbageCollector Clean up connections that have not been used for a while
func garbageCollector() {
	for {
		time.Sleep(garbageCollectorInterval)
		for k, v := range connections {
			if time.Since(v.LastConnected) > defaultConnectionTTL {
				connectionMutex.Lock()
				// Should do debug logging here
				//log.Printf("Removing connection %s", k)
				v.Cancel()
				delete(connections, k)
				connectionMutex.Unlock()
			}
		}
	}
}

// resolve t
func (s *Service) resolveDeviceID(ctx context.Context, deviceId string) (string, error) {

	// Assume we have public key digest
	if _, err := uuid.Parse(deviceId); err != nil {
		return deviceId, nil
	}

	// UUID parsed, we need to resolve it to public key digest
	response, err := s.client.ListDeviceInventory(ctx, client.InventoryListQuery{
		Search: client.InventoryListSearch{
			UUID: deviceId,
		},
	})

	if err != nil {
		return "", err
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("device not found")
	}

	if len(response.Items) > 1 {
		return "", fmt.Errorf("multiple devices found")
	}

	return response.Items[0].PubKeyDigest, nil
}
