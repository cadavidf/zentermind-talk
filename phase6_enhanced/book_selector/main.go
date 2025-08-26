package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Book struct {
	Title       string
	Author      string
	SourceCount int
}

type ContentConcept struct {
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	UniquenessScore  float64 `json:"uniqueness_score"`
	ViabilityScore   bool    `json:"viability_score"`
	CommercialScore  float64 `json:"commercial_score"`
	Status           string  `json:"status"`
	FailureReason    string  `json:"failure_reason"`
	ContentType      string  `json:"content_type"`
}

var ( 
	books []Book
)

func init() {
	books = generateBookList()
}

func generateBookList() []Book {
	bookMap := make(map[string]Book)

	// Initial hardcoded books
	initialBooks := []Book{
		{"The Women", "Kristin Hannah", 0},
		{"The Heaven & Earth Grocery Store", "James McBride", 0},
		{"Fourth Wing", "Rebecca Yarros", 0},
		{"Holly", "Stephen King", 0},
		{"The Ballad of Songbirds and Snakes", "Suzanne Collins", 0},
	}

	// Top selling books 2024
	top2024Books := []Book{
		{"The Women", "Kristin Hannah", 0},
		{"A Court of Thorns and Roses", "Sarah J. Maas", 0},
		{"The Scarlet Shedder (Dog Man #12)", "Dav Pilkey", 0},
		{"The Housemaid", "Freida McFadden", 0},
		{"It Ends with Us", "Colleen Hoover", 0},
		{"Iron Flame", "Rebecca Yarros", 0},
		{"Fourth Wing", "Rebecca Yarros", 0},
		{"Funny Story", "Emily Henry", 0},
		{"The Heaven & Earth Grocery Store", "James McBride", 0},
		{"Atomic Habits", "James Clear", 0},
		{"The Anxious Generation", "Jonathan Haidt", 0},
		{"Killers of the Flower Moon", "David Grann", 0},
		{"The Demon of Unrest", "Erik Larson", 0},
	}

	// Most anticipated books 2025
	ant2025Books := []Book{
		{"A new novel by Chimamanda Ngozi Adichie", "Chimamanda Ngozi Adichie", 0},
		{"Good Dirt", "Charmaine Wilkerson", 0},
		{"The Emperor of Gladness", "Ocean Vuong", 0},
		{"My Friends", "Fredrik Backman", 0},
		{"Dream Count: A Novel", "", 0},
		{"Everybody Says It's Everything: A Novel", "", 0},
		{"The Dissenters: A Novel", "", 0},
		{"Death of the Author: A Novel", "Nnedi Okorafor", 0},
		{"A new novel from Karen Russell", "Karen Russell", 0},
		{"Katabasis", "R.F. Kuang", 0},
		{"Bury Our Bones in the Midnight Soil", "V.E. Schwab", 0},
		{"New sci-fi from John Scalzi", "John Scalzi", 0},
		{"New horror from Grady Hendrix and Stephen Graham Jones", "Grady Hendrix and Stephen Graham Jones", 0},
		{"Birth of a Dynasty", "Chinaza Bado", 0},
		{"The Scrolls of Bishop Eubulus and Other Stories", "Rebecca Bradley", 0},
		{"All the Ash We Leave Behind", "R. Robert Cargill", 0},
		{"Not Quite Dead Yet", "Holly Jackson", 0},
		{"New titles from Lisa Jewell, Alex North, Alice Feeney, and Freida McFadden", "Lisa Jewell, Alex North, Alice Feeney, and Freida McFadden", 0},
		{"Savvy Summers and the Sweet Potato Crimes", "Sandra Jackson-Opoku", 0},
		{"The Man Who Died Seven Times", "Yasuhiko Nishizawa", 0},
		{"Glorious Rivals", "Jennifer Lynn Barnes", 0},
		{"The Unraveling of Julia", "", 0},
		{"The First Gentleman", "Clinton and Patterson", 0},
		{"We Are All Guilty Here", "", 0},
		{"It Was Her House First", "", 0},
		{"New series installments from Tracy Deonn (Legendborn) and Suzanne Collins (The Hunger Games)", "Tracy Deonn and Suzanne Collins", 0},
		{"Bones at the Crossroads", "LaDarrion Williams", 0},
		{"The Last Tiger", "Brad Riew and Julia Riew", 0},
		{"Soulmatch", "Rebecca Danzenbaker", 0},
		{"Outsider Kids", "", 0},
		{"New books from Emily Henry and Rebecca Yarros", "Emily Henry and Rebecca Yarros", 0},
		{"Mate", "Ali Hazelwood", 0},
		{"Rose in Chains", "Julie Soto", 0},
		{"Books on Mark Twain, Black history, and the JFK assassination plot", "", 0},
		{"New memoirs from Geraldine Brooks and Neko Case", "Geraldine Brooks and Neko Case", 0},
		{"We Can Do Hard Things", "Glennon Doyle, Abby Wambach, and Amanda Doyle", 0},
		{"Baking Across America", "B. Dylan Hollis", 0},
	}

	// Consolidate and count sources
	for _, book := range initialBooks {
		key := book.Title + "-" + book.Author
		if existingBook, ok := bookMap[key]; ok {
			existingBook.SourceCount++
			bookMap[key] = existingBook
		} else {
			book.SourceCount = 1
			bookMap[key] = book
		}
	}

	for _, book := range top2024Books {
		key := book.Title + "-" + book.Author
		if existingBook, ok := bookMap[key]; ok {
			existingBook.SourceCount++
			bookMap[key] = existingBook
		} else {
			book.SourceCount = 1
			bookMap[key] = book
		}
	}

	for _, book := range ant2025Books {
		key := book.Title + "-" + book.Author
		if existingBook, ok := bookMap[key]; ok {
			existingBook.SourceCount++
			bookMap[key] = existingBook
		} else {
			book.SourceCount = 1
			bookMap[key] = book
		}
	}

	// Convert map back to slice
	var consolidatedBooks []Book
	for _, book := range bookMap {
		consolidatedBooks = append(consolidatedBooks, book)
	}

	return consolidatedBooks
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", bookListHandler)
	http.HandleFunc("/select", selectBookHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func bookListHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func selectBookHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title", http.StatusBadRequest)
		return
	}

	var selectedBook Book
	for _, b := range books {
		if b.Title == title {
			selectedBook = b
			break
		}
	}

	if selectedBook.Title == "" {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	concept := ContentConcept{
		Title:            selectedBook.Title,
		Description:      fmt.Sprintf("A book by %s", selectedBook.Author),
		UniquenessScore:  8.5,
		ViabilityScore:   true,
		CommercialScore:  7.8,
		Status:           "approved",
		FailureReason:    "",
		ContentType:      "book",
	}

	jsonData, err := json.MarshalIndent(concept, "", "  ")
	if err != nil {
		http.Error(w, "Error creating selection file", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile("../selected_book.json", jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing selected_book.json: %v\n", err)
		http.Error(w, "Error saving selection", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Selected book: %s\n", title)
	fmt.Fprintf(w, "You selected: %s. It has been saved for the next step.", title)
}