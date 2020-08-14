## Go from this...

```go
func (h handler) GetUsernameByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	username, err := h.repo.GetUsernameByID(userID)
	if err != nil {
		var errResponse errorResponse
        	var statusCode int

		if errors.Is(sql.ErrNoRows, err) {
			errResponse.Error = fmt.Sprintf("could not find username for user id '%s'", userID)
			statusCode = http.StatusBadRequest	
		} else {
			errResponse.Error = "request couldn't be processed, please try again later"
			statusCode = http.StatusInternalServerError
		}

		bytes, err := json.Marshal(errResponse)
		if err != nil {
		    w.WriteHeader(http.StatusInternalServerError)
		    return
		}

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(statusCode)
		w.Write(bytes)
		return
	}

	bytes, err := json.Marshal(usernameResponse{
		Username: username,
    	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
```

## TO THIS!

```go
func (h handler) GetUsernameByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	username, err := h.repo.GetUsernameByID(userID)
	if err != nil {
		hapi.RespondError(w, err)
		return
	}

	hapi.RespondOK(w, usernameResponse{
		Username: username,
	})
}
```

// Under construction...
