package main

import (
	"database/sql"
	"encoding/json"
	"muxwebapp/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		utils.ErrorResp(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResp(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.ErrorResp(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		utils.ErrorResp(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResp(w, http.StatusCreated, p)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResp(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.ErrorResp(w, http.StatusNotFound, "Product not found")
		default:
			utils.ErrorResp(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONResp(w, http.StatusOK, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResp(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.ErrorResp(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		utils.ErrorResp(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResp(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResp(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		utils.ErrorResp(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResp(w, http.StatusOK, map[string]string{"result": "success"})
}
