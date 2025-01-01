package extensions

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/gouef/diago"
	"github.com/gouef/router"
	"log"
)

type DiagoRouteExtension struct {
	router                *router.Router
	currentRoute          string
	data                  DiagoRouteData
	PanelGenerator        diago.PanelGenerator
	TemplateProvider      diago.TemplateProvider
	JSTemplateProvider    diago.TemplateProvider
	PopupTemplateProvider diago.TemplateProvider
}

type DiagoRoute struct {
	Actual  bool
	Name    string
	Pattern string
	Method  string
}

type DiagoRouteData struct {
	CurrentRoute string
	Routes       []DiagoRoute
}

type DefaultRouteTemplateProvider struct{}

func (p *DefaultRouteTemplateProvider) GetTemplate() string {
	return GetDiagoRoutePanelPopupTemplate()
}

func NewDefaultTemplateProvider() *DefaultRouteTemplateProvider {
	return &DefaultRouteTemplateProvider{}
}

type DefaultRouteJSTemplateProvider struct{}

func (p *DefaultRouteJSTemplateProvider) GetTemplate() string {
	return GetDiagoRoutePanelJSTemplate()
}

func NewDefaultJSTemplateProvider() *DefaultRouteJSTemplateProvider {
	return &DefaultRouteJSTemplateProvider{}
}

type DefaultRoutePopupTemplateProvider struct{}

func (p *DefaultRoutePopupTemplateProvider) GetTemplate() string {
	return GetDiagoRoutePanelPopupTemplate()
}

func NewDefaultPopupTemplateProvider() *DefaultRoutePopupTemplateProvider {
	return &DefaultRoutePopupTemplateProvider{}
}

func NewDiagoRouteExtension(r *router.Router) *DiagoRouteExtension {
	generator := diago.NewDefaultPanelGenerator()
	tmpProvider := NewDefaultTemplateProvider()
	popupProvider := NewDefaultPopupTemplateProvider()
	jsProvider := NewDefaultJSTemplateProvider()
	return &DiagoRouteExtension{
		router:                r,
		TemplateProvider:      tmpProvider,
		PopupTemplateProvider: popupProvider,
		JSTemplateProvider:    jsProvider,
		PanelGenerator:        generator,
		data: DiagoRouteData{
			Routes:       []DiagoRoute{},
			CurrentRoute: "",
		},
	}
}

func (e *DiagoRouteExtension) SetTemplateProvider(provider diago.TemplateProvider) {
	e.TemplateProvider = provider
}

func (e *DiagoRouteExtension) GetTemplateProvider() diago.TemplateProvider {
	return e.TemplateProvider
}

func (e *DiagoRouteExtension) SetPanelGenerator(generator diago.PanelGenerator) {
	e.PanelGenerator = generator
}

func (e *DiagoRouteExtension) GetPanelGenerator() diago.PanelGenerator {
	return e.PanelGenerator
}

func (e *DiagoRouteExtension) GetHtml(c *gin.Context) string {
	if router.IsRelease() {
		return ""
	}
	result, err := e.PanelGenerator.GenerateHTML("diagoRoutePanelPopup", e.PopupTemplateProvider, e.data)

	if err != nil {
		log.Printf("Diago Route Extension: %s", err.Error())
	}
	return result
}
func (e *DiagoRouteExtension) GetJSHtml(c *gin.Context) string {
	if router.IsRelease() {
		return ""
	}

	//result, err := e.generateDiagoPanelJSHTML()
	result, err := e.PanelGenerator.GenerateHTML("diagoRoutePanelJS", e.JSTemplateProvider, e.data)

	if err != nil {
		log.Printf("Diago Route Extension: %s", err.Error())
	}
	return result
}
func (e *DiagoRouteExtension) GetPanelHtml(c *gin.Context) string {
	if router.IsRelease() {
		return ""
	}
	result, err := e.PanelGenerator.GenerateHTML("diagoRoutePanel", e.TemplateProvider, e.data)

	if err != nil {
		log.Printf("Diago Route Extension: %s", err.Error())
		return ""
	}
	return result
}

func (e *DiagoRouteExtension) BeforeNext(c *gin.Context) {
	if router.IsRelease() {
		return
	}
	e.currentRoute = c.FullPath()
	var routes []DiagoRoute
	for _, route := range e.router.GetRoutes() {
		routes = append(routes, DiagoRoute{
			Actual:  e.currentRoute == route.GetPattern(),
			Name:    route.GetName(),
			Pattern: route.GetPattern(),
			Method:  route.GetMethod().String(),
		})
	}

	e.data = DiagoRouteData{
		CurrentRoute: e.currentRoute,
		Routes:       routes,
	}
}
func (e *DiagoRouteExtension) AfterNext(c *gin.Context) {
	if router.IsRelease() {
		return
	}
	e.currentRoute = c.FullPath()
	var routes []DiagoRoute
	for _, route := range e.router.GetRoutes() {
		routes = append(routes, DiagoRoute{
			Actual:  e.currentRoute == route.GetPattern(),
			Name:    route.GetName(),
			Pattern: route.GetPattern(),
			Method:  route.GetMethod().String(),
		})
	}

	e.data = DiagoRouteData{
		CurrentRoute: e.currentRoute,
		Routes:       routes,
	}
}
