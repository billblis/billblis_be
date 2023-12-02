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

var (
	Responsed           model.Credential
	pemasukanResponse   model.PemasukanResponse
	pengeluaranResponse model.PengeluaranResponse
	datauser            model.User
	pemasukan           model.Pemasukan
	pengeluaran         model.Pengeluaran
)

func GCFHandlerSignup(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Responsed)
	}
	err = SignUp(conn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}
	Responsed.Status = true
	Responsed.Message = "Halo " + datauser.Username
	return GCFReturnStruct(Responsed)
}

func GCFHandlerSignin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Responsed model.Credential
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Responsed)
	}
	user, status1, err := SignIn(conn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}
	Responsed.Status = true
	tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Responsed.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Responsed.Message = "Selamat Datang " + user.Username + " di Billblis" + strconv.FormatBool(status1)
		Responsed.Token = tokenstring
	}
	return GCFReturnStruct(Responsed)
}

// USER

func GCFHandlerGetAllUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	userlist, err := GetAllUser(mconn, collectionname)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Get User Success"
	Responsed.Data = userlist

	return GCFReturnStruct(Responsed)
}

func GCFHandlerGetUserFromUsername(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	user, err := GetUserFromUsername(mconn, collectionname, username)
	if err != nil {
		Responsed.Message = "Error retrieving user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Hello user"
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
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
	pemasukanResponse.Status = false
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	err = json.NewDecoder(r.Body).Decode(&pemasukan)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err = InsertPemasukan(mconn, collectionname, pemasukan.Tanggal_masuk, pemasukan.Jumlah_masuk, pemasukan.Sumber, pemasukan.Deskripsi)
	if err != nil {
		pemasukanResponse.Message = "error inserting Pemasukan: " + err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Insert Pemasukan success"
	pemasukanResponse.Data = []model.Pemasukan{pemasukan}
	return GCFReturnStruct(pemasukanResponse)
}

func GCFHandlerGetPemasukanFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pemasukanResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		pemasukanResponse.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(pemasukanResponse)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		pemasukanResponse.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukan, err := GetPemasukanFromID(mconn, collectionname, ID)
	if err != nil {
		pemasukanResponse.Message = err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Get todo success"
	pemasukanResponse.Data = []model.Pemasukan{pemasukan}

	return GCFReturnStruct(pemasukanResponse)
}

// func GCFHandlerGetPemasukanFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	var Response model.PemasukanResponse
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
// 	pemasukan, err := GetPemasukanFromID(dataUser.ID, conn)
// 	if err != nil {
// 		Response.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(Response)
// 	}
// 	Response.Status = true
// 	Response.Message = "Selamat Datang " + dataUser.Email
// 	Response.Data = []model.Pemasukan{pemasukan}
// 	return GCFReturnStruct(Response)
// }

func GCFHandlerGetAllPemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pemasukanResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukan, err := GetAllPemasukan(mconn, collectionname)
	if err != nil {
		pemasukanResponse.Message = err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Get pemasukan success"
	pemasukanResponse.Data = pemasukan

	return GCFReturnStruct(pemasukanResponse)
}

// func GCFHandlerGetAllPemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)

// 	pemasukanResponse.Status = false
// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		pemasukanResponse.Message = "error parsing application/json1:"
// 		return GCFReturnStruct(pemasukanResponse)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		pemasukanResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(pemasukanResponse)
// 	}
// 	pemasukan, err := GetAllPemasukan(conn)
// 	if err != nil {
// 		pemasukanResponse.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(pemasukanResponse)
// 	}
// 	pemasukanResponse.Status = true
// 	pemasukanResponse.Message = "Berhasil mendapatkan semua pemasukan"
// 	pemasukanResponse.Data = pemasukan
// 	return GCFReturnStruct(pemasukanResponse)
// }

func GCFHandlerUpdatePemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)

	pemasukanResponse.Status = false

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		pemasukanResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	err := json.NewDecoder(r.Body).Decode(&pemasukan)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}
	err = UpdatePemasukan(conn, pemasukan)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}
	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Pemasukan berhasil diupdate"
	pemasukanResponse.Data = []model.Pemasukan{pemasukan}
	return GCFReturnStruct(pemasukanResponse)
}

func GCFHandlerDeletePemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pemasukanResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		pemasukanResponse.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(pemasukanResponse)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		pemasukanResponse.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err = DeletePemasukan(mconn, collectionname, ID)
	if err != nil {
		pemasukanResponse.Message = err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Delete todo success"

	return GCFReturnStruct(pemasukanResponse)
}

// func GCFHandlerDeletePemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)

// 	pemasukanResponse.Status = false

// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		pemasukanResponse.Message = "error parsing application/json1:"
// 		return GCFReturnStruct(pemasukanResponse)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		pemasukanResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(pemasukanResponse)
// 	}
// 	err := DeletePemasukan(conn, pemasukan)
// 	if err != nil {
// 		pemasukanResponse.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(pemasukanResponse)
// 	}
// 	pemasukanResponse.Status = true
// 	pemasukanResponse.Message = "Pemasukan berhasil dihapus"
// 	return GCFReturnStruct(pemasukanResponse)
// }

// PENGELUARAN

func GCFHandlerInsertPengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	err = json.NewDecoder(r.Body).Decode(&pengeluaran)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	_, err = InsertPengeluaran(mconn, collectionname, pengeluaran.Tanggal_keluar, pengeluaran.Jumlah_keluar, pengeluaran.Sumber, pengeluaran.Deskripsi)
	if err != nil {
		pengeluaranResponse.Message = "error inserting Pengeluaran: " + err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Insert Pengeluaran success"
	pengeluaranResponse.Data = []model.Pengeluaran{pengeluaran}
	return GCFReturnStruct(pengeluaranResponse)
}

func GCFHandlerGetPengeluaranFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pengeluaranResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		pengeluaranResponse.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(pengeluaranResponse)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		pengeluaranResponse.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaran, err := GetPengeluaranFromID(mconn, collectionname, ID)
	if err != nil {
		pengeluaranResponse.Message = err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Get todo success"
	pengeluaranResponse.Data = []model.Pengeluaran{pengeluaran}

	return GCFReturnStruct(pengeluaranResponse)
}

// func GCFHandlerGetPengeluaranFromID(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	pengeluaranResponse.Status = false

// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		pengeluaranResponse.Message = "error parsing application/json1:" + token
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		pengeluaranResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}
// 	pengeluaran, err := GetPengeluaranFromID(datauser.ID, conn)
// 	if err != nil {
// 		pengeluaranResponse.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}
// 	pengeluaranResponse.Status = true
// 	pengeluaranResponse.Message = "Selamat Datang " + datauser.Email
// 	pengeluaranResponse.Data = []model.Pengeluaran{pengeluaran}
// 	return GCFReturnStruct(pengeluaranResponse)
// }

func GCFHandlerGetAllPengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pengeluaranResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaran, err := GetAllPengeluaran(mconn, collectionname)
	if err != nil {
		pengeluaranResponse.Message = err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Get pengeluaran success"
	pengeluaranResponse.Data = pengeluaran

	return GCFReturnStruct(pengeluaranResponse)
}

// func GCFHandlerGetAllPengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	pengeluaranResponse.Status = false
// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		pengeluaranResponse.Message = "error parsing application/json1:"
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		pengeluaranResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}
// 	pengeluaran, err := GetAllPengeluaran(conn)
// 	if err != nil {
// 		pengeluaranResponse.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}
// 	pengeluaranResponse.Status = true
// 	pengeluaranResponse.Message = "Berhasil mendapatkan semua pengeluaran"
// 	pengeluaranResponse.Data = pengeluaran
// 	return GCFReturnStruct(pengeluaranResponse)
// }

func GCFHandlerUpdatePengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pengeluaranResponse.Status = false

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

	if err1 != nil {
		pengeluaranResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	err := json.NewDecoder(r.Body).Decode(&pengeluaran)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}
	err = UpdatePengeluaran(conn, pengeluaran)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}
	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Pengeluaran berhasil diupdate"
	pengeluaranResponse.Data = []model.Pengeluaran{pengeluaran}
	return GCFReturnStruct(pengeluaranResponse)
}

func GCFHandlerDeletePengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pengeluaranResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		pengeluaranResponse.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(pengeluaranResponse)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		pengeluaranResponse.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(pengeluaranResponse)
	}
	_, err = DeletePengeluaran(mconn, collectionname, ID)
	if err != nil {
		pengeluaranResponse.Message = err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Delete todo success"

	return GCFReturnStruct(pengeluaranResponse)
}

// func GCFHandlerDeletePengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
// 	pengeluaranResponse.Status = false

// 	// get token from header
// 	token := r.Header.Get("Authorization")
// 	token = strings.TrimPrefix(token, "Bearer ")
// 	if token == "" {
// 		pengeluaranResponse.Message = "error parsing application/json1:"
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}

// 	// decode token
// 	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)

// 	if err1 != nil {
// 		pengeluaranResponse.Message = "error parsing application/json2: " + err1.Error() + ";" + token
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}
// 	err := DeletePengeluaran(conn, pengeluaran)
// 	if err != nil {
// 		pengeluaranResponse.Message = "error parsing application/json4: " + err.Error()
// 		return GCFReturnStruct(pengeluaranResponse)
// 	}
// 	pengeluaranResponse.Status = true
// 	pengeluaranResponse.Message = "Pengeluaran berhasil dihapus"
// 	return GCFReturnStruct(pengeluaranResponse)
// }

// return
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// get id
func GetID(r *http.Request) string {
	return r.URL.Query().Get("id")
}
