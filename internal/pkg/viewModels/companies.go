package viewModels

import (
	"html/template"
	"htmx-example/internal/pkg/models"
	"htmx-example/internal/pkg/storage"
	"htmx-example/internal/pkg/web"
	"net/http"
	"time"
)

type CompaniesViewModel struct {
	templates   *template.Template
	jsonStorage *storage.JsonStorage
}

func NewCompaniesViewModel(templates *template.Template, jsonStorage *storage.JsonStorage) *CompaniesViewModel {
	return &CompaniesViewModel{
		templates:   templates,
		jsonStorage: jsonStorage,
	}
}

func (instance CompaniesViewModel) Index(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	return web.RenderResponse(http.StatusOK, instance.templates, "index.html", companies.All(), nil)
}

func (instance CompaniesViewModel) AddCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"enterEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row-add.html", nil, headers)
}

func (instance CompaniesViewModel) SaveNewCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	row := models.Company{}

	err := request.ParseForm()
	if err != nil {
		return web.GetEmptyResponse(http.StatusBadRequest, nil)
	}

	row.Company = request.Form.Get("company")
	row.Contact = request.Form.Get("contact")
	row.Country = request.Form.Get("country")
	row.Employees = request.Form.Get("employees")
	row.Employees = request.Form.Get("employees")

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	companies.Add(&row)

	err = instance.jsonStorage.Write(companies)
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", row, headers)
}

func (instance CompaniesViewModel) CancelSaveNewCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.GetEmptyResponse(http.StatusOK, headers)
}

func (instance CompaniesViewModel) EditCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	row := companies.GetByID(id)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"enterEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row-edit.html", row, headers)
}

func (instance CompaniesViewModel) SaveExistingCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	row := companies.GetByID(id)

	err = request.ParseForm()
	if err != nil {
		return web.GetEmptyResponse(http.StatusBadRequest, nil)
	}

	row.Company = request.Form.Get("company")
	row.Contact = request.Form.Get("contact")
	row.Country = request.Form.Get("country")
	row.Employees = request.Form.Get("employees")
	companies.Update(row)

	err = instance.jsonStorage.Write(companies)
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", row, headers)
}

func (instance CompaniesViewModel) CancelSaveExistingCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	row := companies.GetByID(id)
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", row, headers)
}

func (instance CompaniesViewModel) DeleteCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	id := request.PathValue("id")

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	companies.Delete(id)

	err = instance.jsonStorage.Write(companies)
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	return web.GetEmptyResponse(http.StatusOK, nil)
}
