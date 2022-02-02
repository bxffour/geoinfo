package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/countries", app.createCountryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/all", app.showCountriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/name/:name", app.showCountriesByNameHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/code/:code", app.showCountriesCodeHandler)

	return router
}
