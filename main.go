package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// Maintainer struct represents the maintainer information
type Maintainer struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// Application struct represents the application metadata
type Application struct {
	ID          string       `json:"id,omitempty" yaml:"id,omitempty"`
	Title       string       `json:"title,omitempty" yaml:"title,omitempty"`
	Version     string       `json:"version,omitempty" yaml:"version,omitempty"`
	Maintainers []Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"`
	Company     string       `json:"company,omitempty" yaml:"company,omitempty"`
	Website     string       `json:"website,omitempty" yaml:"website,omitempty"`
	Source      string       `json:"source,omitempty" yaml:"source,omitempty"`
	License     string       `json:"license,omitempty" yaml:"license,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
}

var (
	applications []Application
	mutex        sync.RWMutex
)

func main() {
	router := setupRouter()

	err := (http.ListenAndServe(":8000", router))
	if err != nil {
		log.Fatal(err)
	}
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/applications", getApplications).Methods("GET")
	router.HandleFunc("/applications/{id}", getApplication).Methods("GET")
	router.HandleFunc("/applications", createApplication).Methods("POST")
	router.HandleFunc("/applications/{id}", updateApplication).Methods("PUT")
	router.HandleFunc("/applications/{id}", deleteApplication).Methods("DELETE")

	return router
}

func getApplications(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	queryParams := r.URL.Query()

	if len(queryParams) == 0 {
		json.NewEncoder(w).Encode(applications)
	} else {
		filteredApplications := make([]Application, 0)

		for _, app := range applications {
			match := true

			for key, values := range queryParams {
				for _, value := range values {
					if !isMatch(key, value, app) {
						match = false
						break
					}
				}

				if !match {
					break
				}
			}

			if match {
				filteredApplications = append(filteredApplications, app)
			}
		}

		json.NewEncoder(w).Encode(filteredApplications)
	}
}

func getApplication(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	params := mux.Vars(r)
	for _, item := range applications {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func createApplication(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	contentType := r.Header.Get("Content-Type")
	var application Application

	switch {
	case strings.Contains(contentType, "application/json"):
		err := json.NewDecoder(r.Body).Decode(&application)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case strings.Contains(contentType, "application/x-yaml"):
		data, _ := io.ReadAll(r.Body)
		err := yaml.Unmarshal(data, &application)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	// Generate a unique ID
	newID := 0
	if len(applications) > 0 {
		lastID := applications[len(applications)-1].ID
		newID, _ = strconv.Atoi(lastID)
	}
	application.ID = strconv.Itoa(newID + 1)

	if err := validateApplication(&application, false); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	applications = append(applications, application)
	json.NewEncoder(w).Encode(applications)
}

func updateApplication(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	params := mux.Vars(r)
	for index, item := range applications {
		if item.ID == params["id"] {
			applications = append(applications[:index], applications[index+1:]...)
			contentType := r.Header.Get("Content-Type")
			var application Application

			switch {
			case strings.Contains(contentType, "application/json"):
				err := json.NewDecoder(r.Body).Decode(&application)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case strings.Contains(contentType, "application/x-yaml"):
				data, _ := io.ReadAll(r.Body)
				err := yaml.Unmarshal(data, &application)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			default:
				http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
				return
			}

			if err := validateApplication(&application, true); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			application.ID = params["id"]
			applications = append(applications, application)
			json.NewEncoder(w).Encode(applications)
			return
		}
	}
	json.NewEncoder(w).Encode(applications)
}

func deleteApplication(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	params := mux.Vars(r)
	for index, item := range applications {
		if item.ID == params["id"] {
			applications = append(applications[:index], applications[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(applications)
}

func isMatch(key, value string, app Application) bool {
	switch key {
	case "title":
		return strings.Contains(strings.ToLower(app.Title), strings.ToLower(value))
	case "version":
		return strings.Contains(strings.ToLower(app.Version), strings.ToLower(value))
	case "maintainer.name":
		for _, maintainer := range app.Maintainers {
			if strings.Contains(strings.ToLower(maintainer.Name), strings.ToLower(value)) {
				return true
			}
		}
		return false
	case "maintainer.email":
		for _, maintainer := range app.Maintainers {
			if strings.Contains(strings.ToLower(maintainer.Email), strings.ToLower(value)) {
				return true
			}
		}
		return false
	case "company":
		return strings.Contains(strings.ToLower(app.Company), strings.ToLower(value))
	case "website":
		return strings.Contains(strings.ToLower(app.Website), strings.ToLower(value))
	case "source":
		return strings.Contains(strings.ToLower(app.Source), strings.ToLower(value))
	case "license":
		return strings.Contains(strings.ToLower(app.License), strings.ToLower(value))
	case "description":
		return strings.Contains(strings.ToLower(app.Description), strings.ToLower(value))
	default:
		return false
	}
}

func validateApplication(app *Application, isUpdate bool) error {
	if isUpdate && app.ID != "" {
		return fmt.Errorf("ID cannot be patched")
	}

	if app.Title == "" {
		return fmt.Errorf("title is required")
	}

	if app.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(app.Maintainers) == 0 {
		return fmt.Errorf("at least one maintainer is required")
	}

	for _, maintainer := range app.Maintainers {
		if maintainer.Name == "" {
			return fmt.Errorf("maintainer name is required")
		}

		if maintainer.Email == "" {
			return fmt.Errorf("maintainer email is required")
		}

		if !isValidEmail(maintainer.Email) {
			return fmt.Errorf("invalid maintainer email")
		}
	}

	if app.Company == "" {
		return fmt.Errorf("company is required")
	}

	if app.Website == "" {
		return fmt.Errorf("website is required")
	}

	if app.Source == "" {
		return fmt.Errorf("source is required")
	}

	if app.License == "" {
		return fmt.Errorf("license is required")
	}

	if app.Description == "" {
		return fmt.Errorf("description is required")
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
