package main

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/bxffour/crest-countries/internal/data"
	"github.com/bxffour/crest-countries/internal/validator"
)

func (app *application) getCountriesByNameHandler(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	ctx, span := app.tracer.Start(reqCtx, r.URL.Path)
	defer span.End()

	name := app.readParam(r, "name")

	if !utf8.ValidString(name) {
		app.failedValidationResponse(w, r, map[string]string{"invalid input error": "only valid utf8 characeters are allowed"})
		return
	}

	var input data.Filters

	v := validator.New()

	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetByName(ctx, name, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.recordSpanError(ctx, err, "the requested resource cannot be found")
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

	code := strings.ToUpper(app.readParam(r, "code"))

	country, err := app.models.Countries.GetByCode(ctx, code)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

	params := app.readParam(r, "codes")

	codes := strings.Split(params, ",")

	countries, err := app.models.Countries.GetByCodes(ctx, codes)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

	params := app.readParam(r, "translation")

	country, err := app.models.Countries.GetByTranslation(ctx, params)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

	var input data.Filters

	v := validator.New()

	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1)
	input.PageSize = app.readInt(qs, "page_size", 20)

	if data.ValidateFilters(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	countries, metadata, err := app.models.Countries.GetAll(ctx, input)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

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

	countries, metadata, err := app.models.Countries.GetByCurrency(ctx, currency, input)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

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

	countries, metadata, err := app.models.Countries.GetByLanguage(ctx, language, input)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

	capital := strings.Title(app.readParam(r, "capital"))

	country, err := app.models.Countries.GetByCapital(ctx, capital)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

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

	countries, metadata, err := app.models.Countries.GetByRegion(ctx, region, input)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

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

	countries, metadata, err := app.models.Countries.GetBySubregion(ctx, region, input)
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
	ctx := r.Context()
	ctx, span := app.tracer.Start(ctx, r.URL.Path)
	defer span.End()

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

	countries, metadata, err := app.models.Countries.GetByDemonyms(ctx, demonyns, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"result": countries, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}