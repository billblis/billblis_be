package billblis

import (
	"fmt"
	"testing"

	model "github.com/billblis/billblis_be/model"
	module "github.com/billblis/billblis_be/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "billblis")

// TEST SIGN UP
func TestSignUp(t *testing.T) {
	var doc model.User
	doc.Username = "Marlina Lubis"
	doc.Email = "marlina@gmail.com"
	doc.Password = "marlinalubis12"

	err := module.SignUp(db, "user", doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", doc.Username)
	}
}

// TEST SIGN IN
func TestSignIn(t *testing.T) {
	var doc model.User
	doc.Username = "Marlina Lubis"
	doc.Password = "marlinalubis12"
	user, Status, err := module.SignIn(db, "user", doc)
	fmt.Println("Status :", Status)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Welcome back:", user)
	}
}

// SUMBER

func TestInsertSumber(t *testing.T) {
	var doc model.Sumber
	doc.Nama_sumber = "pen"

	_id, err := module.InsertSumber(db, "sumber", doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil ditambah dengan id :", _id)
	}
}

func TestGetAllSumber(t *testing.T) {
	var docs []model.Sumber
	docs, err := module.GetAllSumber(db)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", docs)
	}
	fmt.Println(docs)
}

func TestGetSumberFromID(t *testing.T) {
	id := "65657dc24c1690d49d426f44"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		user, err := module.GetSumberFromID(objectId, db)
		if err != nil {
			t.Errorf("Error getting document: %v", err)
		} else {
			fmt.Println(user)
		}
	}
}

// PEMASUKAN

func TestInsertPemasukan(t *testing.T) {
	var doc model.Pemasukan
	doc.Tanggal_masuk = "26/02/2023"
	doc.Jumlah_masuk = 50000
	doc.Sumber = "Gaji"
	doc.Deskripsi = "dari kantor"

	username := "Fedhira Syaila"

	id, err := module.InsertPemasukan(db, "pemasukan", doc, username)
	if err != nil {
		t.Errorf("Error inserting pemasukan: %v", err)
	}
	fmt.Println(id)
}

// func TestGetAllPemasukan(t *testing.T) {
// 	doc, err := module.GetAllPemasukan(db, "pemasukan")
// 	if err != nil {
// 		t.Errorf("Error getting pemasukan: %v", err)
// 		return
// 	}
// 	fmt.Println(doc)
// }

func TestGetPemasukanFromUser(t *testing.T) {
	user := "Huang Renjun"
	doc, err := module.GetPemasukanFromUser(db, "pemasukan", user)
	if err != nil {
		t.Errorf("Error getting pemasukan: %v", err)
		return
	}
	fmt.Println(doc)
}

func TestGetPemasukanFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("65755a9c5f92c5e6fb964960")
	doc, err := module.GetPemasukanFromID(db, "pemasukan", id)
	if err != nil {
		t.Errorf("Error getting pemasukan: %v", err)
		return
	}
	fmt.Println(doc)
}

func TestUpdatePemasukan(t *testing.T) {
	var doc model.Pemasukan
	doc.Tanggal_masuk = "22/02/2023"
	doc.Jumlah_masuk = 230000
	doc.Sumber = "Freelance"
	doc.Deskripsi = "dari joki ngoding"

	id := "65756ad1aa89d76ea9193564"

	ID, err := primitive.ObjectIDFromHex(id)
	doc.ID = ID
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := module.UpdatePemasukan(db, "pemasukan", doc)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updating pemasukan with id: %v", err)
			return
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
		fmt.Println(doc)
	}
}

// func TestUpdatePemasukan(t *testing.T) {
// 	var doc model.Pemasukan
// 	doc.Tanggal_masuk = "22/02/2023"
// 	doc.Jumlah_masuk = 230000
// 	doc.Sumber = "Freelance"
// 	doc.Deskripsi = "dari joki ngoding"
// 	// doc.User.Username = "Fedhira"
// 	id, err := primitive.ObjectIDFromHex("656b136035e729467ff7874e")
// 	doc.ID = id
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil diubah dengan id")
// 	} else {

// 		err = module.UpdatePemasukan(db, doc)
// 		if err != nil {
// 			t.Errorf("Error updateting document: %v", err)
// 		} else {
// 			fmt.Println("Data berhasil diubah dengan id :", doc.ID)
// 		}
// 	}

// }

func TestDeletePemasukan(t *testing.T) {
	id := "65755a9c5f92c5e6fb964960"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	} else {

		status, err := module.DeletePemasukan(db, "pemasukan", ID)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error deleting document: %v", err)
			return
		} else {
			fmt.Println("Delete success")
		}
	}
}

// func TestDeletePemasukan(t *testing.T) {
// 	var doc model.Pemasukan
// 	id, err := primitive.ObjectIDFromHex("6564639a6e6e2f66eee84ddd")
// 	doc.ID = id
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil dihapus dengan id")
// 	} else {

// 		err = module.DeletePemasukan(db, doc)
// 		if err != nil {
// 			t.Errorf("Error updateting document: %v", err)
// 		} else {
// 			fmt.Println("Data berhasil dihapus dengan id :", doc.ID)
// 		}
// 	}
// }

// PENGELUARAN

func TestInsertPengeluaran(t *testing.T) {
	var doc model.Pengeluaran
	doc.Tanggal_keluar = "02/12/2023"
	doc.Jumlah_keluar = 50000
	doc.Sumber = "Konsumsi"
	doc.Deskripsi = "makan ayam"

	username := "Fedhira Syaila"

	id, err := module.InsertPengeluaran(db, "pengeluaran", doc, username)
	if err != nil {
		t.Errorf("Error inserting pemasukan: %v", err)
	}
	fmt.Println(id)
}

func TestGetPengeluaranFromUser(t *testing.T) {
	user := "Fedhira Syaila"
	doc, err := module.GetPengeluaranFromUser(db, "pengeluaran", user)
	if err != nil {
		t.Errorf("Error getting pengeluaran: %v", err)
		return
	}
	fmt.Println(doc)
}

// func TestGetAllPengeluaran(t *testing.T) {
// 	doc, err := module.GetAllPengeluaran(db, "pengeluaran")
// 	if err != nil {
// 		t.Errorf("Error getting pengeluaran: %v", err)
// 		return
// 	}
// 	fmt.Println(doc)
// }

func TestGetPengeluaranFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("657563cca63753461e33bade")
	doc, err := module.GetPengeluaranFromID(db, "pengeluaran", id)
	if err != nil {
		t.Errorf("Error getting pengeluaran: %v", err)
		return
	}
	fmt.Println(doc)
}

// func TestUpdatePemasukan(t *testing.T) {
// 	var doc model.Pemasukan
// 	doc.Tanggal_masuk = "22/02/2023"
// 	doc.Jumlah_masuk = 230000
// 	doc.ID_sumber.ID, _ = primitive.ObjectIDFromHex("65643f6242ef94271017c94a")
// 	doc.Deskripsi = "dari joki ngoding"
// 	doc.ID_user.ID, _ = primitive.ObjectIDFromHex("65631b4de009209dea4dc55e")
// 	id, err := primitive.ObjectIDFromHex("6565676bb3e79ceef0540910")
// 	doc.ID = id
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil diubah dengan id")
// 	} else {
// 		err = module.UpdatePemasukan(db, doc)
// 		if err != nil {
// 			t.Errorf("Error updateting document: %v", err)
// 		} else {
// 			fmt.Println("Data berhasil diubah dengan id :", doc.ID)
// 		}
// 	}
// }

func TestUpdatePengeluaran(t *testing.T) {
	var doc model.Pengeluaran
	doc.Tanggal_keluar = "22/02/2023"
	doc.Jumlah_keluar = 230000
	doc.Sumber = "Kesehatan"
	doc.Deskripsi = "ke rs"

	id := "65756b69063b69c7d78d5ab5"

	ID, err := primitive.ObjectIDFromHex(id)
	doc.ID = ID
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := module.UpdatePengeluaran(db, "pengeluaran", doc)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updating pengeluaran with id: %v", err)
			return
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
		fmt.Println(doc)
	}
}

// func TestUpdatePengeluaran(t *testing.T) {
// 	var doc model.Pengeluaran
// 	doc.Tanggal_keluar = "22/02/2023"
// 	doc.Jumlah_keluar = 230000
// 	doc.Sumber = "Kesehatan"
// 	doc.Deskripsi = "ke rs"
// 	// doc.User.Username = "Fedhira"
// 	id, err := primitive.ObjectIDFromHex("656b13fde033611d0d6e691a")
// 	doc.ID = id
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil diubah dengan id")
// 	} else {

// 		err = module.UpdatePengeluaran(db, doc)
// 		if err != nil {
// 			t.Errorf("Error updateting document: %v", err)
// 		} else {
// 			fmt.Println("Data berhasil diubah dengan id :", doc.ID)
// 		}
// 	}

// }

func TestDeletePengeluaran(t *testing.T) {
	id := "657563cca63753461e33bade"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	} else {

		status, err := module.DeletePengeluaran(db, "pengeluaran", ID)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error deleting document: %v", err)
			return
		} else {
			fmt.Println("Delete success")
		}
	}
}

// func TestDeletePengeluaran(t *testing.T) {
// 	var doc model.Pengeluaran
// 	id, err := primitive.ObjectIDFromHex("6565774763f64428805965ef")
// 	doc.ID = id
// 	if err != nil {
// 		fmt.Printf("Data tidak berhasil dihapus dengan id")
// 	} else {

// 		err = module.DeletePengeluaran(db, doc)
// 		if err != nil {
// 			t.Errorf("Error updateting document: %v", err)
// 		} else {
// 			fmt.Println("Data berhasil dihapus dengan id :", doc.ID)
// 		}
// 	}
// }

// TEST GET USER

func TestGetUserFromID(t *testing.T) {
	id := "65631b4de009209dea4dc55e"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	}

	doc, err := module.GetUserFromID(db, "user", ID)
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(doc)
}

func TestGetUserFromEmail(t *testing.T) {
	doc, _ := module.GetUserFromEmail(db, "user", "yellow12@gmail.com")
	fmt.Println(doc)
}

func TestGetUserFromUsername(t *testing.T) {
	doc, err := module.GetUserFromUsername(db, "user", "Fedhira Syaila")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(doc)
}

func TestGetAllUser(t *testing.T) {
	doc, err := module.GetAllUser(db, "user")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(doc)
}
