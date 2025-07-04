package models

import (
	"github.com/google/uuid"
)

type Company struct {
	ID        string
	Company   string
	Contact   string
	Country   string
	Employees int
}

type Companies struct {
	Companies []Company `json:"companies"`
}

func NewCompanies() *Companies {
	return &Companies{}
}

func (c *Companies) All() []Company {
	return c.Companies
}

func (c *Companies) GetByID(id string) Company {
	var result Company
	for _, i := range c.Companies {
		if i.ID == id {
			result = i
			break
		}
	}
	return result
}

func (c *Companies) Update(company Company) {
	var result []Company
	for _, i := range c.Companies {
		if i.ID == company.ID {
			i.Company = company.Company
			i.Contact = company.Contact
			i.Country = company.Country
			i.Employees = company.Employees
		}
		result = append(result, i)
	}
	c.Companies = result
}

func (c *Companies) Add(company *Company) {
	id := uuid.New().String()

	company.ID = id
	c.Companies = append(c.Companies, Company{
		ID:        company.ID,
		Company:   company.Company,
		Contact:   company.Contact,
		Country:   company.Country,
		Employees: company.Employees,
	})
}

func (c *Companies) Delete(id string) {
	var result []Company
	for _, i := range c.Companies {
		if i.ID != id {
			result = append(result, i)
		}
	}
	c.Companies = result
}
