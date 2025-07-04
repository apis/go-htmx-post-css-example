package viewModels

import (
	"html/template"
	"htmx-example/internal/pkg/models"
	"htmx-example/internal/pkg/storage"
	"htmx-example/internal/pkg/web"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// CompanyViewModel wraps a Company and adds an OrdinalID field and sort parameters
type CompanyViewModel struct {
	Company    models.Company
	OrdinalID  int
	SortColumn string
	SortDir    string
}

// CompaniesTableViewModel represents the data needed to render the companies table
type CompaniesTableViewModel struct {
	Companies  []CompanyViewModel
	SortColumn string
	SortDir    string
}

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

// getCompaniesViewModel is a helper method that contains the common logic for retrieving
// and preparing the companies view model
func (instance CompaniesViewModel) getCompaniesViewModel(request *http.Request, simulatedDelay int, templateName string) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	// Get sort parameters from query string
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	companies, err := instance.jsonStorage.Read()
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	// Convert Company to CompanyViewModel
	companyViewModels := make([]CompanyViewModel, 0, len(companies.All()))
	for _, company := range companies.All() {
		companyViewModels = append(companyViewModels, CompanyViewModel{
			Company:    company,
			OrdinalID:  0, // Will be set after sorting
			SortColumn: sortColumn,
			SortDir:    sortDir,
		})
	}

	// Sort the companies based on the sort column and direction
	sortCompanies(companyViewModels, sortColumn, sortDir)

	// Assign ordinal IDs based on the current sort order
	for i := range companyViewModels {
		companyViewModels[i].OrdinalID = i + 1
	}

	// Create the view model for the template
	tableViewModel := CompaniesTableViewModel{
		Companies:  companyViewModels,
		SortColumn: sortColumn,
		SortDir:    sortDir,
	}

	return web.RenderResponse(http.StatusOK, instance.templates, templateName, tableViewModel, nil)
}

func (instance CompaniesViewModel) Index(request *http.Request, simulatedDelay int) *web.Response {
	return instance.getCompaniesViewModel(request, simulatedDelay, "index.html")
}

func (instance CompaniesViewModel) Companies(request *http.Request, simulatedDelay int) *web.Response {
	return instance.getCompaniesViewModel(request, simulatedDelay, "companies.html")
}

// sortCompanies sorts the companies based on the specified column and direction
func sortCompanies(companies []CompanyViewModel, column, direction string) {
	sort.Slice(companies, func(i, j int) bool {
		// Determine sort order (ascending or descending)
		ascending := direction == "asc"

		// Compare based on the specified column
		switch column {
		case "ID":
			if ascending {
				return companies[i].Company.ID < companies[j].Company.ID
			}
			return companies[i].Company.ID > companies[j].Company.ID
		case "Company":
			if ascending {
				return companies[i].Company.Company < companies[j].Company.Company
			}
			return companies[i].Company.Company > companies[j].Company.Company
		case "Contact":
			if ascending {
				return companies[i].Company.Contact < companies[j].Company.Contact
			}
			return companies[i].Company.Contact > companies[j].Company.Contact
		case "Country":
			if ascending {
				return companies[i].Company.Country < companies[j].Company.Country
			}
			return companies[i].Company.Country > companies[j].Company.Country
		case "Employees":
			if ascending {
				return companies[i].Company.Employees < companies[j].Company.Employees
			}
			return companies[i].Company.Employees > companies[j].Company.Employees
		default:
			if ascending {
				return companies[i].Company.ID < companies[j].Company.ID
			}
			return companies[i].Company.ID > companies[j].Company.ID
		}
	})
}

// getSortedOrdinalID returns the ordinal ID for a company based on the current sort order
func (instance CompaniesViewModel) getSortedOrdinalID(companies *models.Companies, id string, request *http.Request) int {
	// Get sort parameters from query string or use defaults
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	// Convert Company to CompanyViewModel
	companyViewModels := make([]CompanyViewModel, 0, len(companies.All()))
	for _, company := range companies.All() {
		companyViewModels = append(companyViewModels, CompanyViewModel{
			Company:    company,
			OrdinalID:  0, // Will be set after sorting
			SortColumn: sortColumn,
			SortDir:    sortDir,
		})
	}

	// Sort the companies based on the sort column and direction
	sortCompanies(companyViewModels, sortColumn, sortDir)

	// Assign ordinal IDs based on the current sort order
	for i := range companyViewModels {
		companyViewModels[i].OrdinalID = i + 1
	}

	// Find the ordinal ID for the specified company
	for _, company := range companyViewModels {
		if company.Company.ID == id {
			return company.OrdinalID
		}
	}

	return 0 // Not found
}

func (instance CompaniesViewModel) AddCompany(request *http.Request,
	simulatedDelay int) *web.Response {
	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	// Get sort parameters from query string or use defaults
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	// Create a view model with sort parameters
	viewModel := struct {
		SortColumn string
		SortDir    string
	}{
		SortColumn: sortColumn,
		SortDir:    sortDir,
	}

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"enterEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row-add.html", viewModel, headers)
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

	// Convert employees to int and validate
	employeesStr := request.Form.Get("employees")
	employees, err := strconv.Atoi(employeesStr)
	if err != nil {
		return web.GetEmptyResponse(http.StatusBadRequest, nil)
	}
	row.Employees = employees

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

	// Get sort parameters from query string or use defaults
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	// Get the ordinal ID based on the current sort order
	ordinalID := instance.getSortedOrdinalID(companies, row.ID, request)

	// Create a view model with company and sort parameters
	viewModel := CompanyViewModel{
		Company:    row,
		OrdinalID:  ordinalID,
		SortColumn: sortColumn,
		SortDir:    sortDir,
	}

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", viewModel, headers)
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

	// Get sort parameters from query string or use defaults
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	// Get the ordinal ID based on the current sort order
	ordinalID := instance.getSortedOrdinalID(companies, id, request)

	// Create a view model with company and sort parameters
	viewModel := CompanyViewModel{
		Company:    row,
		OrdinalID:  ordinalID,
		SortColumn: sortColumn,
		SortDir:    sortDir,
	}

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"enterEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row-edit.html", viewModel, headers)
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

	// Convert employees to int and validate
	employeesStr := request.Form.Get("employees")
	employees, err := strconv.Atoi(employeesStr)
	if err != nil {
		return web.GetEmptyResponse(http.StatusBadRequest, nil)
	}
	row.Employees = employees

	companies.Update(row)

	err = instance.jsonStorage.Write(companies)
	if err != nil {
		return web.GetEmptyResponse(http.StatusInternalServerError, nil)
	}

	time.Sleep(time.Duration(simulatedDelay) * time.Millisecond)

	// Get sort parameters from query string or use defaults
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	// Get the ordinal ID based on the current sort order
	ordinalID := instance.getSortedOrdinalID(companies, id, request)

	// Create a view model with company and sort parameters
	viewModel := CompanyViewModel{
		Company:    row,
		OrdinalID:  ordinalID,
		SortColumn: sortColumn,
		SortDir:    sortDir,
	}

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", viewModel, headers)
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

	// Get sort parameters from query string or use defaults
	sortColumn := request.URL.Query().Get("sort")
	sortDir := request.URL.Query().Get("dir")
	if sortColumn == "" {
		sortColumn = "ID" // Default sort column
	}
	if sortDir == "" {
		sortDir = "asc" // Default sort direction
	}

	// Get the ordinal ID based on the current sort order
	ordinalID := instance.getSortedOrdinalID(companies, id, request)

	// Create a view model with company and sort parameters
	viewModel := CompanyViewModel{
		Company:    row,
		OrdinalID:  ordinalID,
		SortColumn: sortColumn,
		SortDir:    sortDir,
	}

	headers := map[string]string{"HX-Trigger-After-Swap": "{\"exitEditMode\":\"\"}"}
	return web.RenderResponse(http.StatusOK, instance.templates, "row.html", viewModel, headers)
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
