package viewModels

import (
	"html/template"
	"htmx-example/internal/pkg/models"
	"htmx-example/internal/pkg/web"
	"net/http"
	"time"
)

type CompaniesViewModel struct {
	templates *template.Template
	companies *models.Companies
}

func NewCompaniesViewModel(templates *template.Template, companies *models.Companies) *CompaniesViewModel {
	return &CompaniesViewModel{
		templates: templates,
		companies: companies,
	}
}

func (instance CompaniesViewModel) Index(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "index.html", instance.companies.Companies(), nil)
}

func (instance CompaniesViewModel) AddCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "row-add.html", nil, nil)
}

func (instance CompaniesViewModel) SaveNewCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	row := models.Company{}

	err := request.ParseForm()
	if err != nil {
		return web.GetEmptyResponse(http.StatusBadRequest)
	}

	row.Company = request.Form.Get("company")
	row.Contact = request.Form.Get("contact")
	row.Country = request.Form.Get("country")
	instance.companies.Add(&row)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", row, nil)
}

func (instance CompaniesViewModel) CancelSaveNewCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "companies.html", instance.companies.Companies(), nil)
}

func (instance CompaniesViewModel) EditCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")
	row := instance.companies.GetByID(id)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "row-edit.html", row, nil)
}

func (instance CompaniesViewModel) SaveExistingCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")
	row := instance.companies.GetByID(id)

	err := request.ParseForm()
	if err != nil {
		return web.GetEmptyResponse(http.StatusBadRequest)
	}

	row.Company = request.Form.Get("company")
	row.Contact = request.Form.Get("contact")
	row.Country = request.Form.Get("country")
	instance.companies.Update(row)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", row, nil)
}

func (instance CompaniesViewModel) CancelSaveExistingCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")
	row := instance.companies.GetByID(id)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", row, nil)
}

func (instance CompaniesViewModel) DeleteCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")
	instance.companies.Delete(id)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)
	return web.GetEmptyResponse(http.StatusOK)
}
