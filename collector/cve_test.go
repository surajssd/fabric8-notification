package collector_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fabric8-services/fabric8-notification/auth"
	authApi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/testsupport"
	"github.com/fabric8-services/fabric8-notification/wit"
	witApi "github.com/fabric8-services/fabric8-notification/wit/api"

	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
)

func createLocalClient(t *testing.T, witURL, authURL string) (*witApi.Client, *authApi.Client) {
	witClient, err := wit.NewCachedClient(witURL)
	if err != nil {
		t.Fatal(err)
	}

	authClient, err := auth.NewCachedClient(authURL)
	if err != nil {
		t.Fatal(err)
	}
	return witClient, authClient
}

func TestCVEResolver(t *testing.T) {
	witServer := createServer(serveWITRequest)
	authServer := createServer(serveAuthRequest)

	witURL := "http://" + witServer.Listener.Addr().String() + "/"
	authURL := "http://" + authServer.Listener.Addr().String() + "/"

	witClient, authClient := createLocalClient(t, witURL, authURL)

	cveResolver := collector.NewCVEResolver(authClient, witClient)
	codebaseURL := "git@github.com:testrepo/testproject1.git"
	recvs, _, err := cveResolver(context.Background(), codebaseURL, uuid.NewV4())

	assert.Nil(t, err)
	assert.NotNil(t, recvs)
	assert.Equal(t, 2, len(recvs))
	checkEmails(t, recvs, "testuser1@redhat.com", "testuser2@redhat.com")
}

func createServer(handle func(http.ResponseWriter, *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	return httptest.NewServer(mux)
}

func serveWITRequest(rw http.ResponseWriter, req *http.Request) {
	var err error
	var res string

	reqPath := req.URL.Path
	if reqPath == "/api/search/codebases" {
		codebaseURL := req.URL.Query().Get("url")
		if strings.Contains(codebaseURL, "testproject1") {
			res, err = testsupport.GetFileContent("test-files/cve/codebase.search.testproject1.json")
		}
	} else if strings.HasPrefix(reqPath, "/api/spaces/") {
		if strings.HasSuffix(reqPath, "012be0a3-8ae9-4e09-9e54-35ea34957731") {
			res, err = testsupport.GetFileContent("test-files/cve/space.012be0a3-8ae9-4e09-9e54-35ea34957731.json")
		} else if strings.HasSuffix(reqPath, "f438b002-958b-4bfa-af1c-1580b03c1138") {
			res, err = testsupport.GetFileContent("test-files/cve/space.f438b002-958b-4bfa-af1c-1580b03c1138.json")
		}
	}

	if err != nil {
		return
	}
	rw.Write([]byte(res))
}

func serveAuthRequest(rw http.ResponseWriter, req *http.Request) {
	var err error
	var res string

	reqPath := req.URL.Path
	if strings.HasPrefix(reqPath, "/api/users") {
		if strings.HasSuffix(reqPath, "a06fa7e7-46f3-4947-988b-d708d990f491") {
			res, err = testsupport.GetFileContent("test-files/cve/testuser1.a06fa7e7-46f3-4947-988b-d708d990f491.json")
		} else if strings.HasSuffix(reqPath, "8d72071c-6b89-4f14-90f3-869faa950018") {
			res, err = testsupport.GetFileContent("test-files/cve/testuser2.8d72071c-6b89-4f14-90f3-869faa950018.json")
		}
	}

	if err != nil {
		return
	}
	rw.Write([]byte(res))
}

func checkEmails(t *testing.T, recvs []collector.Receiver, emails ...string) {
	t.Helper()
	emailsNotFound := make([]string, 0, 0)
	for _, email := range emails {
		found := false
		for _, recv := range recvs {
			if recv.EMail == email {
				found = true
				break
			}
		}
		if !found {
			emailsNotFound = append(emailsNotFound, email)
		}
	}
	if len(emailsNotFound) > 0 {
		t.Errorf("%d emails %s not found in receivers", len(emailsNotFound), emailsNotFound)
	}
}
