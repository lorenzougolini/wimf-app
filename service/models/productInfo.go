package models

type ProductInfo struct {
	Barcode string `json:"_id"`
	Name    string `json:"product_name"`
	NameIT  string `json:"product_name_it"`
	NameEN  string `json:"product_name_en"`
	Brand   string `json:"brands"`
}
