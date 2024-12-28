package tests

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/gouef/router"
	"github.com/gouef/router/extensions"
	"github.com/stretchr/testify/assert"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDiagoPanelGenerator struct{}

func (m *mockDiagoPanelGenerator) GenerateHTML(data interface{}) (string, error) {
	return "<div>Mocked HTML</div>", nil
}

func TestDiagoRouteExtension(t *testing.T) {
	gin.SetMode(router.TestMode)
	r := router.NewRouter()
	r.EnableTest()
	d := diago.NewDiago()
	n := r.GetNativeRouter()

	routeExtension := extensions.NewDiagoRouteExtension(r)
	d.AddExtension(routeExtension)
	n.Use(diago.DiagoMiddleware(r, d))

	r.AddRouteGet("test", "/test", func(c *gin.Context) {
		panelHtml := template.HTML("<div>Test</div>")
		c.HTML(http.StatusOK, "", panelHtml)
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	n.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "Routes")
}

func TestDiagoRouteExtension_GetHtml(t *testing.T) {
	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	html := routeExtension.GetHtml(nil)
	assert.Contains(t, html, "Routes")
}

func TestDiagoRouteExtension_GetJSHtml(t *testing.T) {
	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	jsHtml := routeExtension.GetJSHtml(nil)
	assert.Contains(t, jsHtml, "function closeRoutesPopup()")
}

type mockDiagoPanelGeneratorWithError struct{}

func (m *mockDiagoPanelGeneratorWithError) GenerateHTML(name string, templateProvider diago.TemplateProvider, data interface{}) (string, error) {
	return "", errors.New("mock error generating HTML")
}

func TestDiagoRouteExtension_GetPanelHtml_ErrorHandling(t *testing.T) {
	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	gen := &mockDiagoPanelGeneratorWithError{}
	routeExtension.SetPanelGenerator(gen)

	var logOutput string
	log.SetOutput(&logWriter{&logOutput})

	panelHtml := routeExtension.GetPanelHtml(nil)
	assert.Empty(t, panelHtml, "Panel HTML should be empty when there's an error")

	assert.Contains(t, logOutput, "Diago Route Extension: mock error generating HTML", "Error message should be logged")
}

type mockTemplateProviderWithParseError struct{}

func (m *mockTemplateProviderWithParseError) GetTemplate() string {
	return "{{ .Latencys }}"
}

type mockTemplateProviderWithExecuteError struct{}

func (m *mockTemplateProviderWithExecuteError) GetTemplate() string {
	return `{{ .NonExistentField }}`
}

func TestGenerateDiagoPanelHTML_TemplateParseError(t *testing.T) {
	mockProvider := &mockTemplateProviderWithParseError{}

	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	result, err := routeExtension.PanelGenerator.GenerateHTML("error", mockProvider, struct {
	}{})

	assert.Error(t, err, "Expected error while parsing template")
	assert.Empty(t, result, "Expected empty result when parsing fails")
}

type mockInvalidTemplateProvider struct{}

func (m *mockInvalidTemplateProvider) GetDiagoLatencyPanelTemplate() string {
	return "{{ .InvalidField"
}

func TestGenerateDiagoPanelHTML_TemplateExecuteError(t *testing.T) {
	mockProvider := &mockTemplateProviderWithExecuteError{}

	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	result, err := routeExtension.PanelGenerator.GenerateHTML("test", mockProvider, nil)

	assert.Error(t, err, "Expected error while executing template")
	assert.Empty(t, result, "Expected empty result when execution fails")
}

func TestDiagoRouteExtension_TemplateProvider(t *testing.T) {
	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	mockProvider := &mockTemplateProviderWithParseError{}

	routeExtension.SetTemplateProvider(mockProvider)

	assert.Equal(t, mockProvider, routeExtension.GetTemplateProvider(), "TemplateProvider should be set correctly")
}

func TestDiagoRouteExtension_PanelGenerator(t *testing.T) {
	r := router.NewRouter()
	routeExtension := extensions.NewDiagoRouteExtension(r)

	mockProvider := &mockTemplateProviderWithParseError{}

	routeExtension.SetTemplateProvider(mockProvider)

	assert.Equal(t, mockProvider, routeExtension.GetTemplateProvider(), "TemplateProvider should be set correctly")
}

type logWriter struct {
	output *string
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	*lw.output = string(p)
	return len(p), nil
}
