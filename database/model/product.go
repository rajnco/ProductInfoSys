package model


// Product - Product Table struct
type Product struct {
        Id          int     `json:"id"          gorm:"primaryKey;type:int;autoIncrement"`
        Name        string  `json:"name"        gorm:"type:varcha(64)"`
        Description string  `json:"description" gorm:"type:varchar(256)"`
        Price       float32 `json:"price"       gorm:"type:numeric"`                    //UnitPrice
        Quantity    int     `json:"quantity"    gorm:"type:int"`
        Discount    int     `json:"discount"    gorm:"type:int"`                        //MaxDiscountPercent
        Country     string  `json:"country"     gorm:"type:varchar(64)"`
        Region      string  `json:"region"      gorm:"type:varchar(8)"`
}

func (Product) TableName() string{
	return "products"
}

type CreateProduct struct {
        Name        string  `json:"name"` 
        Description string  `json:"description"`
        Price       float32 `json:"price"`                    //UnitPrice
        Quantity    int     `json:"quantity"`
        Discount    int     `json:"discount"`                 //MaxDiscountPercent
        Country     string  `json:"country"`
        Region      string  `json:"region"`
}

