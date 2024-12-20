package main

type applicationConfig struct {
	Host           string `config_default:"localhost" config_description:"Server host interface"`
	Port           int    `config_default:"8080" config_description:"Server port"`
	SimulatedDelay int    `config_default:"0" config_description:"Simulated delay in milliseconds for HTMX interactions"`
}
