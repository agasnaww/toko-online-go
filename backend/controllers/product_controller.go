package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"toko-online-go/config"
	"toko-online-go/models"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		rows, err := config.DB.Query("SELECT id, nama_barang, harga, stok FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var hasil []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok); err != nil {
				log.Println("Error scan:", err)
				continue
			}
			hasil = append(hasil, p)
		}
		if hasil == nil {
			hasil = []models.Product{}
		}
		json.NewEncoder(w).Encode(hasil)

	case "POST":
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Data JSON tidak valid", http.StatusBadRequest)
			return
		}
		query := "INSERT INTO products (nama_barang, harga, stok) VALUES (?, ?, ?)"
		_, err := config.DB.Exec(query, p.Nama, p.Harga, p.Stok)
		if err != nil {
			http.Error(w, "Gagal menyimpan", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sukses tambah barang!"})

	case "PUT":
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID tidak dapat ditemukan!", http.StatusBadRequest)
			return // <--- JANGAN LUPA RETURN DI SINI
		}

		// PERBAIKAN 1: Tambahkan 'models.' di depan Product
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Data JSON rusak", http.StatusBadRequest)
			return
		}

		query := "UPDATE products SET nama_barang=?, harga=?, stok=? WHERE id=?"

		_, err := config.DB.Exec(query, p.Nama, p.Harga, p.Stok, id)
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

		_, err := config.DB.Exec("DELETE FROM products WHERE id=?", id)
		if err != nil {
			http.Error(w, "Gagal menghapus data", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Barang telah dihapus!"})

	default:
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
	}
}
