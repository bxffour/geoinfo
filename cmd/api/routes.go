package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/countries/all", app.getAllCountriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/name/:name", app.getCountriesByNameHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/code/:code", app.getCountryByCodeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/codes/:codes", app.getCountriesByCodesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/currency/:currency", app.getCountriesByCurrencyHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/lang/:lang", app.getCountriesByLanguageHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/capital/:capital", app.getCountriesByCapitalHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/region/:region", app.getCountriesByRegionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/subregion/:region", app.getCountriesBySubregionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/demonym/:demonym", app.getCountriesByDemonymHandler)
	router.HandlerFunc(http.MethodGet, "/v1/countries/translation/:translation", app.getCountryByTranslationHandler)

	return app.otelhttp(app.rateLimit(app.recoverPanic(app.rateLimit(router))))
}
