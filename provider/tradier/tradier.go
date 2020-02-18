package tradier

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	tradierHost               = "api.tradier.com"
	tradierURL                = "https://" + tradierHost + "/v1/"
	tradierFormat             = "application/json"
	tradierAuthFile           = "../../resources/provider-auth/tradier"
	tradierTimeoutSec         = 30
	tradierRequestDelayMillis = 500 // Avoid rate limiting by keeping requests below 2 per second.

	nanoToMilli = 1e6
)

var tradierProvider *Tradier

type Tradier struct {
	sync.Mutex
	client     *http.Client
	lastMillis int64  // Last request timestamp
	auth       string // Authorization token
}

func init() {
	// Setup HTTPS client.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: tradierHost, // Required for the SNI client requirement.
				// VerifyPeerCertificate: showCertDNSNames, // To display DNS names on the server certificate.
			},
		},
		Timeout: tradierTimeoutSec * time.Second,
	}

	// Get authorization token.
	authToken, err := ioutil.ReadFile(tradierAuthFile)
	if err != nil {
		log.Fatalln("Could not retrieve Tradier authorization token:", err)
	}

	// Initialize Tradier provider.
	tradierProvider = &Tradier{
		client:     client,
		lastMillis: time.Now().UnixNano() / nanoToMilli,
		auth:       string(authToken),
	}
}

func New() *Tradier {
	return tradierProvider
}
