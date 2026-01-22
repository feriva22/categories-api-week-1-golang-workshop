package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var categories = []Category{}

// Create init function to create db file using json data
var dbFileName string = "db-categories.json"
var dataInit = []byte(`[
	{
		"id": 1,
		"name": "Electronics"
	},
	{
		"id": 2,
		"name": "Books"
	}
]`)

func initDBFile(buffer *[]byte) []byte {
	data := &dataInit
	// You can also write it to a file as a whole.
	err := os.WriteFile(dbFileName, *data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Success create db file " + dbFileName)

	buffer = data
	return *buffer
}

func writeCurrentDatatoDBFile() {
	//masukkan data ke json file
	dataBufferNew, err := json.Marshal(&categories)
	err = os.WriteFile(dbFileName, dataBufferNew, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	//open db file
	data, err := os.ReadFile(dbFileName)
	if err != nil {
		// Init DB file if error
		log.Print(err)
		log.Print("Creating new db file....")
		data = initDBFile(&data)
	}

	fmt.Println("Reading from existing file")
	// 3. Unmarshal the JSON buffer into the slice
	err = json.Unmarshal(data, &categories)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed %d categories successfully.\n", len(categories))
}

// GET /api/category/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category Belum Ada", http.StatusBadRequest)
}

// PUT /api/category/{id}
func updateCategoryID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	//loop produk cari yang id nya sama
	for i := range categories {
		if categories[i].ID == id {
			categories[i] = updateCategory

			writeCurrentDatatoDBFile()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	http.Error(w, "Category Belum Ada", http.StatusBadRequest)
}

// DELETE /api/category/{id}
func deleteCategoryID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	//loop categories cari yang id nya sama
	for i, c := range categories {
		if c.ID == id {

			//buat slice list produk sebelum dihapus produknya dan setelahnya
			categories = append(categories[:i], categories[i+1:]...)

			writeCurrentDatatoDBFile()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Berhasil delete id",
			})
			return
		}
	}

	http.Error(w, "Category Belum Ada", http.StatusBadRequest)
}

func main() {

	// get detail produk
	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			getCategoryByID(w, r)
		}

		if r.Method == "PUT" {
			updateCategoryID(w, r)
		}

		if r.Method == "DELETE" {
			updateCategoryID(w, r)
		}
	})

	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			//baca dari request
			var categoryNew Category
			err := json.NewDecoder(r.Body).Decode(&categoryNew)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			//masukkan data ke var
			categoryNew.ID = len(categories) + 1
			categories = append(categories, categoryNew)

			writeCurrentDatatoDBFile() //call to write into db

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) //201

			json.NewEncoder(w).Encode(categories)
		}
	})

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
