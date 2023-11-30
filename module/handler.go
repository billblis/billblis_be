package module

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "github.com/billblis/billblis_be/model"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GCFHandlerSignup(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUp(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + dataUser.Name
	return GCFReturnStruct(Response)
}

func GCFHandlerSignin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	user, status1, err := SignIn(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	// Response.Message = "Halo " + dataUser.Name
	tokenstring, err := watoken.Encode(dataUser.Email, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Email + " di Billblis" + strconv.FormatBool(status1)
		Response.Token = tokenstring
	}
	return GCFReturnStruct(Response)
}

func GCFGetUserFromID(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV, dbname, collectionname string, r *http.Request) string {
	var Response model.Credential
	Response.Status = false
	db := MongoConnect(MONGOCONNSTRINGENV, dbname)

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	profile, err := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)
	if err != nil {
		Response.Message = "Error decoding token: " + err.Error()
		return GCFReturnStruct(Response)
	}

	userID, err := primitive.ObjectIDFromHex(profile.Id)
	if err != nil {
		Response.Message = "Invalid user ID in token"
		return GCFReturnStruct(Response)
	}

	user, err := GetUserFromID(userID, db)
	if err != nil {
		Response.Message = "Error fetching user: " + err.Error()
		return GCFReturnStruct(Response)
	}

	return GCFReturnStruct(user)
}

// func GCFHandlerGetUserFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	db := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Credential
// 	Response.Status = false
// 	var dataUser model.User

// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		Response.Message = "error parsing application/json1:" + token
// 		return GCFReturnStruct(Response)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(Response)
// 	}
// 	user, err := GetUserFromID(dataUser.ID, db)
// 	if err != nil {
// 		Response.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	Response.Status = true
// 	Response.Message = "Hello user"
// 	Response.Data = []model.User{user}
// 	return GCFReturnStruct(Response)
// }

// func GCFHandlerGetUserFromName(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.Credential
// 	Response.Status = false
// 	var dataUser model.User

// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		Response.Message = "error parsing application/json1:" + token
// 		return GCFReturnStruct(Response)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(Response)
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&dataUser)
// 	if err != nil {
// 		Response.Message = "error parsing application/json3: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	user, err := GetUserFromName(dataUser.Name, conn)
// 	if err != nil {
// 		Response.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	Response.Status = true
// 	Response.Message = "Hello user"
// 	Response.Data = []model.User{user}
// 	return GCFReturnStruct(Response)
// }

// SUMBER
func GCFHandlerInsertSumber(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.SumberResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var dataSumber model.Sumber

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&dataSumber)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = InsertSumber(mconn, collectionname, dataSumber)
	if err != nil {
		Response.Message = "error inserting Sumber: " + err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Insert Sumber success"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetSumberFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.SumberResponse
	Response.Status = false
	var dataUser model.User

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	sumber, err := GetSumberFromID(dataUser.ID, conn)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Selamat Datang " + dataUser.Email
	Response.Data = []model.Sumber{sumber}
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllSumber(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.SumberResponse
	Response.Status = false
	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	sumber, err := GetAllSumber(conn)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil mendapatkan semua sumber"
	Response.Data = sumber
	return GCFReturnStruct(Response)
}

// PEMASUKAN

func GCFHandlerInsertPemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.PemasukanResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var dataPemasukan model.Pemasukan

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&dataPemasukan)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = InsertPemasukan(mconn, collectionname, dataPemasukan.Tanggal_masuk, dataPemasukan.Jumlah_masuk, dataPemasukan.ID_sumber, dataPemasukan.Deskripsi, dataPemasukan.ID_user)
	if err != nil {
		Response.Message = "error inserting Pemasukan: " + err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Insert Pemasukan success"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetPemasukanFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PemasukanResponse
	Response.Status = false
	var dataUser model.User

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	pemasukan, err := GetPemasukanFromID(dataUser.ID, conn)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Selamat Datang " + dataUser.Email
	Response.Data = []model.Pemasukan{pemasukan}
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllPemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PemasukanResponse
	Response.Status = false
	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	pemasukan, err := GetAllPemasukan(conn)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil mendapatkan semua pemasukan"
	Response.Data = pemasukan
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdatePemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PemasukanResponse
	Response.Status = false
	var dataPemasukan model.Pemasukan

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataPemasukan)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdatePemasukan(conn, dataPemasukan)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Pemasukan berhasil diupdate"
	Response.Data = []model.Pemasukan{dataPemasukan}
	return GCFReturnStruct(Response)
}

func GCFHandlerDeletePemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PemasukanResponse
	Response.Status = false
	var dataPemasukan model.Pemasukan

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	err := DeletePemasukan(conn, dataPemasukan)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Pemasukan berhasil dihapus"
	return GCFReturnStruct(Response)
}

// PENGELUARAN

func GCFHandlerInsertPengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response model.PengeluaranResponse
	Response.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var dataPengeluaran model.Pengeluaran

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&dataPengeluaran)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = InsertPengeluaran(mconn, collectionname, dataPengeluaran.Tanggal_keluar, dataPengeluaran.Jumlah_keluar, dataPengeluaran.ID_sumber, dataPengeluaran.Deskripsi, dataPengeluaran.ID_user)
	if err != nil {
		Response.Message = "error inserting Pengeluaran: " + err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Insert Pengeluaran success"
	return GCFReturnStruct(Response)
}

func GCFHandlerGetPengeluaranFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PengeluaranResponse
	Response.Status = false
	var dataUser model.User

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	pengeluaran, err := GetPengeluaranFromID(dataUser.ID, conn)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Selamat Datang " + dataUser.Email
	Response.Data = []model.Pengeluaran{pengeluaran}
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllPengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PengeluaranResponse
	Response.Status = false
	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	pengeluaran, err := GetAllPengeluaran(conn)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Berhasil mendapatkan semua pengeluaran"
	Response.Data = pengeluaran
	return GCFReturnStruct(Response)
}

func GCFHandlerUpdatePengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PengeluaranResponse
	Response.Status = false
	var dataPengeluaran model.Pengeluaran

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataPengeluaran)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = UpdatePengeluaran(conn, dataPengeluaran)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Pengeluaran berhasil diupdate"
	Response.Data = []model.Pengeluaran{dataPengeluaran}
	return GCFReturnStruct(Response)
}

func GCFHandlerDeletePengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.PengeluaranResponse
	Response.Status = false
	var dataPengeluaran model.Pengeluaran

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}
	err := DeletePengeluaran(conn, dataPengeluaran)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Pengeluaran berhasil dihapus"
	return GCFReturnStruct(Response)
}

// return
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// get id
func GetID(r *http.Request) string {
	return r.URL.Query().Get("id")
}
