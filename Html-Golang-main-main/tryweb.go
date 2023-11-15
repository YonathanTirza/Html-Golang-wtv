package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	// "log"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

var total float64 = 0.00
var itemName string
var itemPrice float64
var templates = template.Must(template.ParseFiles("templates/home.html"))

type FoodItem struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	Items []FoodItem `json:"items"`
}

var (
	order      = Order{}
	orderMutex sync.Mutex
	db         *sql.DB
)

func init() {
	// Open the MySQL database
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/foodiemoodie")
	if err != nil {
		fmt.Println("Error opening database:", err)
	}

	// Check if the connection is established
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
}

func main() {
	router := mux.NewRouter()

	// Define routes for different pages
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/add-to-order", addToOrderHandler).Methods("POST")
	router.HandleFunc("/order-details", orderDetailsHandler).Methods("GET")

	// Serve static files (e.g., CSS, JS) from the "static" directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on :8085")
	http.Handle("/", router)
	http.ListenAndServe(":8085", nil)
}

func renderHTML(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "home")
}

func addToOrderHandler(w http.ResponseWriter, r *http.Request) {
	var foodItem FoodItem
	err := json.NewDecoder(r.Body).Decode(&foodItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderMutex.Lock()
	defer orderMutex.Unlock()

	// Insert the food item into the database
	_, err = db.Exec("INSERT INTO orders (Makanan, Harga) VALUES (?, ?)", foodItem.Name, foodItem.Price)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order.Items = append(order.Items, foodItem)

	w.WriteHeader(http.StatusOK)
}

func orderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	orderMutex.Lock()
	defer orderMutex.Unlock()

	// Retrieve orders from the database
	rows, err := db.Query("SELECT Makanan, Harga FROM orders")
	if err != nil {
		fmt.Println("Error querying database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and populate the order
	order.Items = nil
	for rows.Next() {
		err := rows.Scan(&itemName, &itemPrice)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		order.Items = append(order.Items, FoodItem{Name: itemName, Price: itemPrice})
	}

	// Return the order as JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Error marshalling order to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(orderJSON)

	// _, err = db.Exec("INSERT INTO orders (Makanan,Harga) VALUES (?, ?)", itemName,itemPrice)
    // if err != nil {
    //     log.Fatal(err)
    // }
}
