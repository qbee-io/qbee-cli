package broker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	connections    ConnectionsCache
}

const (
	defaultConnectionTTL  = 5 * time.Minute
	DefaultListenPort     = "8081"
	DefaultRemotePort     = "80"
	DefaultRemoteProtocol = "http"
	DefaultRemoteHost     = "localhost"
)

// NewService creates a new broker service
func NewService() *Service {

	token := os.Getenv("QBEE_TOKEN")

	remoteHost := DefaultRemoteHost
	if os.Getenv("QBEE_REMOTE_HOST") != "" {
		remoteHost = os.Getenv("QBEE_REMOTE_HOST")
	}

	remotePort := DefaultRemotePort
	if os.Getenv("QBEE_REMOTE_PORT") != "" {
		remotePort = os.Getenv("QBEE_REMOTE_PORT")
	}

	listenPort := DefaultListenPort
	if os.Getenv("QBEE_LISTEN_PORT") != "" {
		listenPort = os.Getenv("QBEE_LISTEN_PORT")
	}

	remoteProtocol := DefaultRemoteProtocol
	if os.Getenv("QBEE_REMOTE_PROTOCOL") != "" {
		remoteProtocol = os.Getenv("QBEE_REMOTE_PROTOCOL")
	}

	return &Service{
		authToken:      token,
		remoteHost:     remoteHost,
		remotePort:     remotePort,
		remoteProtocol: remoteProtocol,
		listenPort:     listenPort,
		connections:    NewConnectionsCache(),
	}
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {

	// Sart a goroutine that re-authenticates the client every minute
	go s.reAuthenticateClient(ctx)
	// Start garbage collector
	go s.garbageCollector()

	router := chi.NewRouter()
	router.Use(s.AuthMiddleware)
	router.HandleFunc("/*", s.Proxy())

	if s.authToken == "" {
		log.Println("Warning: No authentication token provided. Device access will be open")
	}

	log.Printf("Starting server on :%s. Press CTRL+C to stop it.", s.listenPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.listenPort), router)
}

// WithListenPort sets the port to listen on (default: 8081)
func (s *Service) WithListenPort(listenPort string) *Service {
	s.listenPort = listenPort
	return s
}

// WithClient sets the client to use when connecting to the device
func (s *Service) WithClient(client *client.Client) *Service {
	s.client = client
	return s
}

// WithAuthToken sets the auth token to use when connecting to the device
func (s *Service) WithAuthToken(token string) *Service {
	s.authToken = token
	return s
}

// WithRemoteHost sets the host to use when connecting to the device (default: localhost)
func (s *Service) WithRemotePort(port string) *Service {
	s.remotePort = port
	return s
}

// WithRemoteProtocol sets the protocol to use when connecting to the device (http or https)
func (s *Service) WithRemoteProtocol(protocol string) *Service {
	s.remoteProtocol = protocol
	return s
}

// each session contains the username of the user and the time at which it expires
type session struct {
	expiry time.Time
}

var sessions = map[string]session{}
var sessionTimeout = 4 * time.Hour
var sessionCookieName = "session_token"

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
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
		if token == s.authToken {
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(sessionTimeout)

			// Set the session in the session map
			sessions[sessionToken] = session{
				expiry: expiresAt,
			}
			cookie := &http.Cookie{
				Name:     sessionCookieName,
				Value:    sessionToken,
				SameSite: http.SameSiteStrictMode,
				Expires:  expiresAt,
			}
			// Set the cookie in the response
			http.SetCookie(w, cookie)
			// Set the session in the request context
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie(sessionCookieName)
		if err == nil {
			sessionToken := c.Value
			userSession, exists := sessions[sessionToken]
			if exists {
				if !userSession.isExpired() {
					next.ServeHTTP(w, r)
					return
				}
				delete(sessions, sessionToken)
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	if conn, ok := s.connections.Get(connectionId); ok {
		conn.LastConnected = time.Now()
		s.connections.Add(connectionId, conn)
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

	// Connect to the device
	go func() {
		errChan <- s.client.Connect(cancelCtx, deviceNodeID, remoteTarget)
	}()

	// Clean up connection if device disconnects or context is cancelled
	go func() {
		defer func() {
			cancelFunc()
			s.connections.Remove(connectionId)
		}()

		select {
		case <-cancelCtx.Done():
			return
		case <-errChan:
			return
		}
	}()

	// Wait for the port to be ready
	portReady := make(chan bool)
	go waitForPort(localPort, portReady)

	select {
	case err := <-errChan:
		if err != nil {
			cancelFunc()
			return 0, err
		}
		log.Printf("do portforwarding: connection to device %s on port %s closed\n", deviceId, devicePort)

	case isPortReady := <-portReady:
		if !isPortReady {
			cancelFunc()
			return 0, fmt.Errorf("port %d is not ready", localPort)
		}
	}

	s.connections.Add(connectionId, Connection{
		Cancel:        cancelFunc,
		Port:          localPort,
		LastConnected: time.Now(),
	})

	// return the local port if the connection was successful
	return localPort, nil
}

// clientReAuthInterval is the interval at which the client will be re-authenticated
const clientReAuthInterval = 4 * time.Hour

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
func (s *Service) garbageCollector() {
	for {
		time.Sleep(garbageCollectorInterval)
		s.connections.CleanUp()
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
