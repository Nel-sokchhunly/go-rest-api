package book

import (
	"encoding/json"
	"fmt"
	"net/http"

	errResp "github.com/Nel-sokchhunly/go-rest-api/cmd/api/resource/common/err"
	validatorUtil "github.com/Nel-sokchhunly/go-rest-api/util/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type API struct {
	respository *Repository
	validator   *validator.Validate
}

func New(db *gorm.DB, v *validator.Validate) *API {
	return &API{
		respository: NewRepository(db),
		validator:   v,
	}
}

// List godoc
//
// @summary List Books
// @description get all books
// @tags books
// @accept  json
// @produce  json
// @success 200 {object} DTO
// @failure 400 {object} err.Error
// @router /books [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	books, err := a.respository.List()
	if err != nil {
		errResp.ServerError(w, errResp.RespDBDataAccessFailure)
		return
	}

	if len(books) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		errResp.ServerError(w, errResp.RespJSONEncodeFailure)
		return
	}
}

// Create godoc
//
// @summary Create Book
// @description create a new book
// @tags books
// @accept  json
// @produce  json
// @param body body Form true "Book form"
// @success 201 {object} DTO
// @failure 400 {object} err.Error
// @failure 422 {object} err.Errors
// @failure 500 {object} err.Error
// @router /books [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errResp.ServerError(w, errResp.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			errResp.ServerError(w, errResp.RespJSONEncodeFailure)
			return
		}

		errResp.ValidationError(w, respBody)
		return
	}

	newBook := form.ToModel()
	newBook.ID = uuid.New()

	book, err := a.respository.Create(newBook)
	if err != nil {
		errResp.ServerError(w, errResp.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		errResp.ServerError(w, errResp.RespJSONEncodeFailure)
		return
	}
}

// Read godoc
//
// @summary        Read book
// @description    Read book
// @tags           books
// @accept         json
// @produce        json
// @param          id     path        string  true    "Book ID"
// @success        200 {object}    DTO
// @failure        400 {object}    err.Error
// @failure        404
// @failure        500 {object}    err.Error
// @router         /books/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	id, err := retriveID(r)
	if err != nil {
		errResp.BadRequest(w, errResp.RespInvalidURLParamID)
		return
	}

	book, err := a.respository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		errResp.ServerError(w, errResp.RespDBDataAccessFailure)
		return
	}

	dto := book.ToDTO()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		errResp.ServerError(w, errResp.RespJSONEncodeFailure)
		return
	}
}

// Update godoc
//
// @summary        Update book
// @description    Update book
// @tags           books
// @accept         json
// @produce        json
// @param          id      path    string  true    "Book ID"
// @param          body    body    Form    true    "Book form"
// @success        200
// @failure        400 {object}    err.Error
// @failure        404
// @failure        422 {object}    err.Errors
// @failure        500 {object}    err.Error
// @router         /books/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := retriveID(r)
	if err != nil {
		errResp.BadRequest(w, errResp.RespInvalidURLParamID)
		return
	}

	form := &Form{} // form is kinda like DTO in this case
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		errResp.ServerError(w, errResp.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			errResp.ServerError(w, errResp.RespJSONEncodeFailure)
			return
		}

		errResp.ValidationError(w, respBody)
		return
	}

	book := form.ToModel() // convert form to model
	book.ID = id

	rows, err := a.respository.Update(book)
	if err != nil {
		errResp.ServerError(w, errResp.RespDBDataUpdateFailure)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Delete godoc
//
// @summary        Delete book
// @description    Delete book
// @tags           books
// @accept         json
// @produce        json
// @param          id  path    string  true    "Book ID"
// @success        200
// @failure        400 {object}    err.Error
// @failure        404
// @failure        500 {object}    err.Error
// @router         /books/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := retriveID(r)
	if err != nil {
		errResp.BadRequest(w, errResp.RespInvalidURLParamID)
		return
	}

	rows, err := a.respository.Delete(id)
	if err != nil {
		errResp.ServerError(w, errResp.RespDBDataRemoveFailure)
		return
	}

	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// util
func retriveID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
