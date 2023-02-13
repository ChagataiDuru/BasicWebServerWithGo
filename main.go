package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Employee struct {
	ID        string   `json:"id"`
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Email     string   `json:"email"`
	Manager   *Manager `json:"manager"`
}

type Manager struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

var employees []Employee

func getEmployees(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getEmployees")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteEmployee")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range employees {
		if item.ID == params["id"] {
			employees = append(employees[:index], employees[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(employees)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getEmployee")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range employees {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Employee{})
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createEmployee")
	w.Header().Set("Content-Type", "application/json")
	var employee Employee
	_ = json.NewDecoder(r.Body).Decode(&employee)
	employee.ID = strconv.Itoa(len(employees) + 1)
	employees = append(employees, employee)
	json.NewEncoder(w).Encode(employee)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateEmployee")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range employees {
		if item.ID == params["id"] {
			employees = append(employees[:index], employees[index+1:]...)
			var employee Employee
			_ = json.NewDecoder(r.Body).Decode(&employee)
			employee.ID = params["id"]
			employees = append(employees, employee)
			json.NewEncoder(w).Encode(employee)
			return
		}
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.Handle("/", http.FileServer(http.Dir("./staticfiles")))
	r.HandleFunc("/main", mainHandler)
	r.HandleFunc("/form", formHandler)

	employees = append(employees, Employee{ID: "1", FirstName: "lorem", LastName: "ipsum", Email: "lorem.ipsum@gmail.com", Manager: &Manager{ID: "1", FirstName: "admin", LastName: "mamdin", Email: "admin@admin.com"}})
	employees = append(employees, Employee{ID: "2", FirstName: "dolor", LastName: "sit", Email: "dolor.sit@gmail.com", Manager: &Manager{ID: "2", FirstName: "admin2", LastName: "admin2", Email: "admin@gmail.com"}})
	r.HandleFunc("/employees", getEmployees).Methods("GET")
	r.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	r.HandleFunc("/employees", createEmployee).Methods("POST")
	r.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	http.Handle("/", r)

	fmt.Print("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
