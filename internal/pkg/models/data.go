package models

import (
	"strconv"
)

type Company struct {
	ID      string
	Company string
	Contact string
	Country string
}

type Companies struct {
	companies []Company
}

func NewCompanies() *Companies {
	data := &Companies{
		companies: []Company{
			{
				ID:      "1",
				Company: "Amazon",
				Contact: "Jeff Bezos",
				Country: "United States",
			},
			{
				ID:      "2",
				Company: "Apple",
				Contact: "Tim Cook",
				Country: "United States",
			},
			{
				ID:      "3",
				Company: "Microsoft",
				Contact: "Satya Nadella",
				Country: "United States",
			},
		},
	}
	return data
}

func (c *Companies) Companies() []Company {
	return c.companies
}

func (c *Companies) GetByID(id string) Company {
	var result Company
	for _, i := range c.companies {
		if i.ID == id {
			result = i
			break
		}
	}
	return result
}

func (c *Companies) Update(company Company) {
	var result []Company
	for _, i := range c.companies {
		if i.ID == company.ID {
			i.Company = company.Company
			i.Contact = company.Contact
			i.Country = company.Country
		}
		result = append(result, i)
	}
	c.companies = result
}

func (c *Companies) Add(company Company) {
	max := 0
	for _, i := range c.companies {
		n, _ := strconv.Atoi(i.ID)
		if n > max {
			max = n
		}
	}
	max++
	id := strconv.Itoa(max)

	c.companies = append(c.companies, Company{
		ID:      id,
		Company: company.Company,
		Contact: company.Contact,
		Country: company.Country,
	})
}

func (c *Companies) Delete(id string) {
	var result []Company
	for _, i := range c.companies {
		if i.ID != id {
			result = append(result, i)
		}
	}
	c.companies = result
}
