package bebillblis

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "time"
)

type User struct {
	ID           primitive.ObjectID 	`bson:"_id,omitempty" json:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	// Email		 string             	`bson:"email,omitempty" json:"email,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// type Pemasukan struct {
// 	ID         			primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
// 	Tanggal_masuk     	time.Time      			`bson:"tanggal_masuk,omitempty" json:"tanggal_masuk,omitempty"`
// 	Jumlah_masuk  		int   			`bson:"jumlah_masuk,omitempty" json:"jumlah_masuk,omitempty"`
// 	ID_sumber 			Sumber  			`bson:"id_sumber,omitempty" json:"id_sumber,omitempty"`
// 	Deskripsi      		string      			`bson:"deskripsi,omitempty" json:"deskripsi,omitempty"`
// 	ID_user       		User 				`bson:"id_user,omitempty" json:"id_user,omitempty"`
// }

// type Pengeluaran struct {
// 	ID         			primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
// 	Tanggal_keluar     	time.Time      			`bson:"tanggal_keluar,omitempty" json:"tanggal_keluar,omitempty"`
// 	Jumlah_keluar  		int  			`bson:"jumlah_keluar,omitempty" json:"jumlah_keluar,omitempty"`
// 	ID_sumber 			Sumber  			`bson:"id_sumber,omitempty" json:"id_sumber,omitempty"`
// 	Deskripsi      		string     			`bson:"deskripsi,omitempty" json:"deskripsi,omitempty"`
// 	ID_user       		User 				`bson:"id_user,omitempty" json:"id_user,omitempty"`
// }

// type Sumber struct {
// 	ID      			 primitive.ObjectID 	`bson:"_id,omitempty" json:"_id,omitempty"`
// 	Nama_sumber   		 string            		`bson:"nama_sumber,omitempty" json:"nama_sumber,omitempty"`
	
// }


