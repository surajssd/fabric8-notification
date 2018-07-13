package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"unsafe"
)

// Config contains basic configuration data for Keycloak
type Config struct {
	BaseURL string
	Realm   string
	Broker  string
}

// RealmAuthURL return endpoint for realm auth config "{BaseURL}/auth/realms/{Realm}/broker/{Broker}/token"
func (c Config) RealmAuthURL() string {
	return fmt.Sprintf("%v/auth/realms/%v", c.BaseURL, c.Realm)
}

// BrokerTokenURL return endpoint to fetch Brokern token "{BaseURL}/auth/realms/{Realm}/broker/{Broker}/token"
func (c Config) BrokerTokenURL() string {
	return fmt.Sprintf("%v/auth/realms/%v/broker/%v/token", c.BaseURL, c.Realm, c.Broker)
}

// OpenshiftToken fetches the Openshift token defined for the current user in Keycloak
func OpenshiftToken(config Config, token string) (string, error) {
	ut, err := get(config.BrokerTokenURL(), token)
	if err != nil {
		return "", err
	}

	return ut.AccessToken, nil
}

type usertoken struct {
	AccessToken string `json:"access_token"`
}

func get(url, token string) (*usertoken, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// for debug only
	rb, _ := httputil.DumpRequest(req, true)
	if false {
		fmt.Println(string(rb))
	}

	client := createHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	b := buf.Bytes()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unknown response:\n%v\n%v", *(*string)(unsafe.Pointer(&b)), string(rb))
	}

	var u usertoken
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
