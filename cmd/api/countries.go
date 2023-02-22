package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bxffour/crest-countries/internal/data"
	"github.com/bxffour/crest-countries/internal/validator"
)

func (app *application) getCountriesByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := app.readParam(r, "name")

	var input data.Filters

	v := validator.New()

	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetByName(name, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}

func (app *application) getCountryByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(app.readParam(r, "code"))

	country, err := app.models.Countries.GetByCode(code)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": country}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesByCodesHandler(w http.ResponseWriter, r *http.Request) {
	params := app.readParam(r, "codes")

	codes := strings.Split(params, ",")

	countries, err := app.models.Countries.GetByCodes(codes)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountryByTranslationHandler(w http.ResponseWriter, r *http.Request) {
	params := app.readParam(r, "translation")

	country, err := app.models.Countries.GetByTranslation(params)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": country}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllCountriesHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Filters

	v := validator.New()

	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetAll(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesByCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currency := strings.ToUpper(app.readParam(r, "currency"))
	var input data.Filters

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetByCurrency(currency, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesByLanguageHandler(w http.ResponseWriter, r *http.Request) {
	language := strings.ToLower(app.readParam(r, "lang"))

	var input data.Filters

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetByLanguage(language, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesByCapitalHandler(w http.ResponseWriter, r *http.Request) {
	capital := strings.Title(app.readParam(r, "capital"))

	country, err := app.models.Countries.GetByCapital(capital)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": country}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesByRegionHandler(w http.ResponseWriter, r *http.Request) {
	region := app.readParam(r, "region")

	var input data.Filters

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetByRegion(region, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesBySubregionHandler(w http.ResponseWriter, r *http.Request) {
	region := app.readParam(r, "region")

	var input data.Filters

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetBySubregion(region, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCountriesByDemonymHandler(w http.ResponseWriter, r *http.Request) {
	demonyns := app.readParam(r, "demonym")

	var input data.Filters

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetByDemonyms(demonyns, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}