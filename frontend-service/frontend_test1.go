package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"okx-bot/frontend-service/controllers"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Initialize logger before all tests
var log = logrus.New()

// Test theCreateAccount controller
func TestCreateAccountController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	body := "{\"email\":\"testuser@mail.ru\",\"password\":\"password\"}"

	req, _ := http.NewRequest("POST", "/api/user/new", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Contains(t, r.Body.String(), "\"username\":\"testuser\"")
}

// test theAuthenticate controller
func TestAuthenticateController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	body := "{\"username\":\"testuser\",\"password\":\"password\"}"

	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	// NOTE: exact response code may vary depending on the implementation details of AuthenticateController.
	// Change this assertion based on the expected status of a successful login.
	expectedStatusCode := http.StatusOK
	assert.Equal(t, expectedStatusCode, r.Code)
	assert.Contains(t, r.Body.String(), "\"username\":\"testuser\"")
}

// Test theCreateContact controller
func TestCreateContactController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")

	body := "{\"contact_id\":\"123\",\"user_id\":\"1\",\"details\":\"test contact\"}"

	req, _ := http.NewRequest("POST", "/api/contacts/new", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Contains(t, r.Body.String(), "\"contact_id\":\"123\"")
}

// Test GetContactsFor controller
func TestGetContactsForController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/me/contacts", nil)
	router.ServeHTTP(r, req)

	// NOTE: exact response code may vary depending on the implementation details of GetContactsForController.
	// Change this assertion based on the expected status of a successful request.
	expectedStatusCode := http.StatusOK
	assert.Equal(t, expectedStatusCode, r.Code)
	assert.Contains(t, r.Body.String(), "\"contact_id\"")
}

// Test the ReceiveSignal controller
func TestReceiveSignalController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/signal/receive", controllers.ReceiveSignal).Methods("POST")

	body := "{\"signal\":\"test signal\"}"

	req, _ := http.NewRequest("POST", "/api/signal/receive", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Contains(t, r.Body.String(), "\"signal\":\"test signal\"")
}

// Test the CreateSignal controller
func TestCreateSignalController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/signal/create", controllers.CreateSignal).Methods("POST")

	body := "{\"signal\":\"test signal\"}"

	req, _ := http.NewRequest("POST", "/api/signal/create", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Contains(t, r.Body.String(), "\"signal\":\"test signal\"")
}

// Test the GetAllSignals controller
func TestGetAllSignalsController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/signal", controllers.GetAllSignals).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/signal", nil)
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"signal\"")
}

// Test the GetBots controller
func TestGetBotsController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/signal/bots", controllers.GetBots).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/signal/bots", nil)
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"bot_id\"")
}

// Test the GetAllOkxBots controller
func TestGetAllOkxBotsController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/signal/okx/bots", controllers.GetAllOkxBots).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/signal/okx/bots", nil)
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"bot_id\"")
}

// Test the GetAllOkxSignals controller
func TestGetAllOkxSignalsController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/signal/okx/all", controllers.GetAllOkxSignals).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/signal/okx/all", nil)
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"signal\"")
}

// Test the CreateBot controller
func TestCreateBotController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/bot/create", controllers.CreateBot).Methods("POST")

	body := "{\"bot_id\":\"bot123\",\"name\":\"test bot\",\"details\":\"details\"}"

	req, _ := http.NewRequest("POST", "/api/bot/create", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Contains(t, r.Body.String(), "\"bot_id\":\"bot123\"")
}

// Test the DeleteBot controller
func TestDeleteBotController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/bot/delete", controllers.DeleteBot).Methods("DELETE")
	req, _ := http.NewRequest("DELETE", "/api/bot/delete", nil)

	// Assuming that the DeleteBot controller expects a URL parameter for the bot ID
	req.URL.RawQuery = "bot_id=bot123"

	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"status\":\"deleted\"")
}

// Test the CreateOkxApi controller
func TestCreateOkxApiController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/okx/create", controllers.CreateOkxApi).Methods("POST")

	body := "{\"api_key\":\"testkey\",\"api_secret\":\"testsecret\",\"passphrase\":\"testpass\",\"subaccount\":\"testsub\"}"
	req, _ := http.NewRequest("POST", "/api/okx/create", bytes.NewBufferString(body))
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Contains(t, r.Body.String(), "\"api_key\":\"testkey\"")
}

// Test the GetOkxApiFor controller
func TestGetOkxApiForController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/okx/keys", controllers.GetOkxApiFor).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/okx/keys", nil)
	req.URL.RawQuery = "user_id=1"

	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"api_key\":\"testkey\"")
}

// Test the CheckOkx controller
func TestCheckOkxController(t *testing.T) {
	r := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/check/okx", controllers.CheckOkx).Methods("GET")

	req, _ := http.NewRequest("GET", "/api/check/okx", nil)
	router.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Contains(t, r.Body.String(), "\"status\":\"ok\"")
}

func TestMain(m *testing.M) {
	// Run all the tests
	m.Run()

	// Coverage test
	coverage := m.Run()
	if coverage != 0 {
		log.Errorf("Non-zero code: %d", coverage)
	}
}
