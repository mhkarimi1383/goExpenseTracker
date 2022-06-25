package httpHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/mhkarimi1383/goExpenseTracker/database"
	"github.com/mhkarimi1383/goExpenseTracker/languageParser"
	"github.com/mhkarimi1383/goExpenseTracker/logger"
	"github.com/mhkarimi1383/goExpenseTracker/types"
)

func index(w http.ResponseWriter, r *http.Request) {
	userDataCookie, err := r.Cookie("user_data")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	userDataJson, err := base64Decode(userDataCookie.Value)
	if err != nil {
		logger.Warnf(true, "error while decoding cookie: %v", err)
		resp := http.StatusText(http.StatusInternalServerError)
		responseWriter(w, &resp, http.StatusInternalServerError)
		return
	}
	userData := make(map[string]any)
	err = json.Unmarshal([]byte(userDataJson), &userData)
	if err != nil {
		logger.Warnf(true, "error while opening cookie: %v", err)
		resp := http.StatusText(http.StatusInternalServerError)
		responseWriter(w, &resp, http.StatusInternalServerError)
		return
	}
	username := fmt.Sprintf("%v", userData["preferred_username"])
	t, err := template.ParseFiles("templates/index.html.gotmpl")
	if err != nil {
		logger.Warnf(true, "error while parsing template: %v", err)
		resp := http.StatusText(http.StatusInternalServerError)
		responseWriter(w, &resp, http.StatusInternalServerError)
	}
	list, err := database.ListItems(username)
	if err != nil {
		logger.Warnf(true, "error while getting items from database: %v", err)
		resp := http.StatusText(http.StatusInternalServerError)
		responseWriter(w, &resp, http.StatusInternalServerError)
		return
	}
	amount := uint(0)
	lastId := uint(0)
	for _, item := range list {
		if item.Operator == "+" {
			amount += item.Amount
		} else {
			amount -= item.Amount
		}
		if item.Id > lastId {
			lastId = item.Id
		}
	}
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		t.Execute(w, &types.IndexPage{
			Title:    information.Title,
			Amount:   amount,
			Items:    list,
			Language: languageParser.GetSelectedLanguage(),
		})
		return
	} else if r.Method == http.MethodPost {
		if r.FormValue("action") == "CREATE" {
			amountStr := r.FormValue("amount")
			amount, err := strconv.Atoi(amountStr)
			if err != nil {
				logger.Warnf(true, "Error converting amount to number: %v", err)
				resp := http.StatusText(http.StatusBadRequest)
				responseWriter(w, &resp, http.StatusBadRequest)
				return
			}
			description := r.FormValue("description")
			operatorStr := r.FormValue("operator")
			operator := operatorCheckboxTranslator(operatorStr)
			_, err = database.InsertItem(username, types.Item{
				Description: description,
				Operator:    operator,
				Amount:      uint(amount),
				Id:          lastId + 1,
			})
			if err != nil {
				logger.Warnf(true, "error creating new item: %v", err)
				resp := http.StatusText(http.StatusInternalServerError)
				responseWriter(w, &resp, http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else if r.FormValue("action") == "DELETE" {
			id, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				logger.Warnf(true, "error while removing item: %v", err)
				resp := http.StatusText(http.StatusInternalServerError)
				responseWriter(w, &resp, http.StatusInternalServerError)
				return
			}
			database.DeleteItem(username, uint(id))
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
}

func indexHandler() http.Handler {
	return http.HandlerFunc(index)
}
