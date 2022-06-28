package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Nana-Seyram/crest-countries/internal/data"
	"github.com/Nana-Seyram/crest-countries/internal/validator"
)

func (app *application) createCountryHandler(w http.ResponseWriter, r *http.Request) {
	var countries []data.Country

	err := app.readJSON(w, r, &countries)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	for _, country := range countries {

		item := &data.Item{
			Country: country,
		}

		err = app.models.Countries.Insert(item)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.writeJSON(w, http.StatusCreated, envelope{"message": "operation was successful"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) showCountriesByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := app.readParam(r, "name")

	countries, err := app.models.Countries.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}

}

func (app *application) showCountriesCodeHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) showCountriesHandler(w http.ResponseWriter, r *http.Request) {
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
