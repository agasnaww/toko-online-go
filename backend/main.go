package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Product struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama_barang"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var db *sql.DB

func main() {
	// Koneksi ke db
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	fmt.Println("Sukses!")

	http.HandleFunc("/products", productsHandler)

	fmt.Println("Server berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, nama_barang, harga, stok FROM products")
		if err != nil {
			http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
			log.Println("Error query:", err)
			return
		}
		defer rows.Close()

		var hasil []Product

		for rows.Next() {
			var p Product

			err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok)
			if err != nil {
				log.Println("Error scan:", err)
			}

			hasil = append(hasil, p)
		}

		if hasil == nil {
			hasil = []Product{}
		}

		json.NewEncoder(w).Encode(hasil)

	case "POST":
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Data JSON tidak valid!", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO products (nama_barang, harga, stok) VALUES (?, ? ,?)"

		_, err := db.Exec(query, p.Nama, p.Harga, p.Stok)

		if err != nil {
			http.Error(w, "Gagal menyimpan ke database", http.StatusInternalServerError)
			log.Println("Error Insert:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sukses menambah barang!"})

	case "PUT":
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID tidak dapat ditemukan!", http.StatusBadRequest)
		}

		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Data JSON rusak", http.StatusBadRequest)
			return
		}

		query := "UPDATE products SET nama_barang=?, harga=?, stok=? WHERE id=?"
		_, err := db.Exec(query, p.Nama, p.Harga, p.Stok, id)
		if err != nil {
			http.Error(w, "Gagal update data", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Barang berhasil diupdate!"})

	case "DELETE":
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID tidak dapat ditemukan!", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("DELETE FROM products WHERE id=?", id)
		if err != nil {
			http.Error(w, "Gagal menghapus data", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Barang telah dihapus!"})

	default:
		http.Error(w, "Method tidak diizinkan!", http.StatusMethodNotAllowed)
	}

}
