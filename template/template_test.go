package template_test

import (
	"context"
	"fmt"
	"testing"

	"strings"

	"github.com/fabric8-services/fabric8-notification/auth"
	authApi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/testsupport"
	"github.com/fabric8-services/fabric8-notification/types"
	"github.com/fabric8-services/fabric8-notification/wit"
	witApi "github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	OpenshiftIOAPI     = "http://api.openshift.io"
	OpenShiftIOAuthAPI = "https://auth.openshift.io"
)

func addGlobalVars(vars map[string]interface{}) map[string]interface{} {
	vars["webURL"] = "http://localhost"
	return vars
}

func createClient(t *testing.T) (*witApi.Client, *authApi.Client) {
	c, err := wit.NewCachedClient(OpenshiftIOAPI)
	if err != nil {
		t.Fatal(err)
	}
	authClient, err := auth.NewCachedClient(OpenShiftIOAuthAPI)
	if err != nil {
		t.Fatal(err)
	}
	return c, authClient
}

func TestTrueOnFoundName(t *testing.T) {
	reg := template.AssetRegistry{}

	_, exist := reg.Get(string(types.WorkitemUpdate))
	assert.True(t, exist)
}

func TestFalseOnMissingName(t *testing.T) {
	reg := template.AssetRegistry{}

	_, exist := reg.Get("MISSING")
	assert.False(t, exist)
}

func TestRenderWorkitemCreate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(string(types.WorkitemCreate))
	assert.True(t, exist)

	witClient, authClient := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	_, vars, err := collector.WorkItem(context.Background(), authClient, witClient, &testsupport.DummyCollaboratorCollector{}, wiID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)
	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "Message-ID")
	assert.Contains(t, headers, "X-OSIO-Space")
	assert.Contains(t, headers, "X-OSIO-Area")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "Ruchir Garg"))
	assert.True(t, strings.Contains(body, "1343"))
	assert.True(t, strings.Contains(body, "mode under Backlog")) // part of the Description. Might change since we're on live data
	assert.True(t, strings.Contains(body, "/openshiftio"))       // Space/Area tag. Might change since we're on live data

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderWorkitemCreateMissingDescription(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(string(types.WorkitemCreate))
	assert.True(t, exist)

	witClient, authClient := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	_, vars, err := collector.WorkItem(context.Background(), authClient, witClient, &testsupport.DummyCollaboratorCollector{}, wiID)
	if err != nil {
		t.Fatal(err)
	}

	wi := vars["workitem"].(*witApi.WorkItemSingle)
	delete(wi.Data.Attributes, "system.description")
	delete(wi.Data.Attributes, "system.description.rendered")
	delete(wi.Data.Attributes, "system.description.markup")

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)
	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "Message-ID")
	assert.Contains(t, headers, "X-OSIO-Space")
	assert.Contains(t, headers, "X-OSIO-Area")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "Ruchir Garg"))
	assert.True(t, strings.Contains(body, "1343"))

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderWorkitemUpdate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(string(types.WorkitemUpdate))
	assert.True(t, exist)

	witClient, authClient := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	_, vars, err := collector.WorkItem(context.Background(), authClient, witClient, &testsupport.DummyCollaboratorCollector{}, wiID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "In-Reply-To")
	assert.Contains(t, headers, "References")
	assert.Contains(t, headers, "X-OSIO-Space")
	assert.Contains(t, headers, "X-OSIO-Area")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "1343"))

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderCommentCreate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(string(types.CommentCreate))
	assert.True(t, exist)

	witClient, authClient := createClient(t)
	ciID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	_, vars, err := collector.Comment(context.Background(), authClient, witClient, &testsupport.DummyCollaboratorCollector{}, ciID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "In-Reply-To")
	assert.Contains(t, headers, "References")
	assert.Contains(t, headers, "X-OSIO-Space")
	assert.Contains(t, headers, "X-OSIO-Area")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "1343"))
	assert.True(t, strings.Contains(body, "just shows loading animation,")) // part of the msg

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderCommentUpdate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(string(types.CommentUpdate))
	assert.True(t, exist)

	witClient, authClient := createClient(t)
	ciID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	_, vars, err := collector.Comment(context.Background(), authClient, witClient, &testsupport.DummyCollaboratorCollector{}, ciID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "In-Reply-To")
	assert.Contains(t, headers, "References")
	assert.Contains(t, headers, "X-OSIO-Space")
	assert.Contains(t, headers, "X-OSIO-Area")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "1343"))
	assert.True(t, strings.Contains(body, "just shows loading animation,")) // part of the ms

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderEmailUpdate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(string(types.UserEmailUpdate))
	assert.True(t, exist)

	_, authClient := createClient(t)
	ciID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	_, vars, err := collector.User(context.Background(), authClient, ciID)
	if err != nil {
		t.Fatal(err)
	}

	if vars == nil {
		vars = map[string]interface{}{}
	}
	vars["custom"] = map[string]interface{}{
		"verifyURL": "https://verift.url.openshift.io",
	}
	_, body, _, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)
	assert.True(t, strings.Contains(body, "https://verift.url.openshift.io"))
}

func TestRenderInvitationSpaceNoorg(t *testing.T) {
	reg := template.AssetRegistry{}

	template, exist := reg.Get(string(types.InvitationSpaceNoorg))
	assert.True(t, exist)

	vars := map[string]interface{}{}
	vars["custom"] = map[string]interface{}{
		"inviter":   "John Smith",
		"spaceName": "Customer Orders",
		"roleNames": "Contributor, Project Lead",
		"acceptURL": "http://openshift.io/invitations/accept/12345-ABCDE-FFFFF-99999-88888",
	}

	_, body, _, err := template.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(body, "http://openshift.io/invitations/accept/12345-ABCDE-FFFFF-99999-88888"))

	//ioutil.WriteFile("../invitation-space.html", []byte(body), os.FileMode(0777))
}

func TestRenderInvitationTeamNoorg(t *testing.T) {
	reg := template.AssetRegistry{}

	template, exist := reg.Get(string(types.InvitationTeamNoorg))
	assert.True(t, exist)

	vars := map[string]interface{}{}
	vars["custom"] = map[string]interface{}{
		"teamName":  "Developers",
		"inviter":   "John Smith",
		"spaceName": "Financial Backend",
		"acceptURL": "http://openshift.io/invitations/accept/12345-ABCDE-FFFFF-99999-77777",
	}

	_, body, _, err := template.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(body, "http://openshift.io/invitations/accept/12345-ABCDE-FFFFF-99999-77777"))

	//ioutil.WriteFile("../invitation-team.html", []byte(body), os.FileMode(0777))
}

func TestRenderCVE(t *testing.T) {
	files := []string{"cve.basic", "cve.many"}

	reg := template.AssetRegistry{}
	template, exist := reg.Get(string(types.AnalyticsNotifyCVE))
	assert.True(t, exist)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			vars := make(map[string]interface{})
			payload, err := testsupport.GetFileContent(fmt.Sprintf("test-files/%s.json", file))
			require.NoError(t, err)
			vars["custom"] = testsupport.GetCustomElement(payload)

			sub, body, _, err := template.Render(addGlobalVars(vars))
			require.NoError(t, err)

			custom := toMap(vars["custom"])
			assert.True(t, strings.Contains(sub, toString(custom["repo_url"])))
			checkCVEBody(t, body, custom)
		})
	}
}

func checkCVEBody(t *testing.T, body string, custom map[string]interface{}) {
	t.Helper()
	assert.True(t, strings.Contains(body, toString(custom["repo_url"])))
	assert.True(t, strings.Contains(body, toString(custom["scanned_at"])))

	depArr := toArrMap(custom["vulnerable_deps"])
	assert.NotNil(t, depArr)

	for _, dep := range depArr {
		assert.True(t, strings.Contains(body, toString(dep["name"])))
		cveArr := toArrMap(dep["cves"])
		assert.NotNil(t, depArr)
		for _, cve := range cveArr {
			assert.True(t, strings.Contains(body, toString(cve["CVE"])))
			cvss := fmt.Sprintf("[%s/10]", toString(cve["CVSS"]))
			assert.True(t, strings.Contains(body, cvss))
		}
	}
}

func toString(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func toMap(val interface{}) map[string]interface{} {
	if m, ok := val.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func toArrMap(val interface{}) []map[string]interface{} {
	if arr, ok := val.([]interface{}); ok {
		res := make([]map[string]interface{}, 0, 0)
		for _, e := range arr {
			val := toMap(e)
			if val != nil {
				res = append(res, val)
			} else {
				return nil
			}
		}
		return res
	}
	return nil
}
