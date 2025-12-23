package models

type Product struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama_barang"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}
