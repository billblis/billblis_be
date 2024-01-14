package module

import (
	"encoding/json"
	"net/http"
	"os"

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
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := SignIn(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	tokenstring, err := watoken.Encode(user.Username, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Responsed.Message = "Gagal Encode Token :" + err.Error()

	} else {
		Responsed.Message = "Selamat Datang " + user.Username
		Responsed.Token = tokenstring
		Responsed.Data = []model.User{user}
	}

	return GCFReturnStruct(Responsed)
}

// USER

func GCFHandlerDeleteUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	_, err = DeleteUser(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error deleting user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Delete user " + datauser.Username + " success"

	return GCFReturnStruct(Responsed)
}

func GCFHandlerChangePassword(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := ChangePassword(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error changing password: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Password change success for user " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

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

func GCFHandlerGetUserFromToken(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	// Konversi string ke primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(userInfo.Id)
	if err != nil {
		Responsed.Message = "error converting userID:" + err.Error()
		return GCFReturnStruct(Responsed)
	}

	user, err := GetUserFromToken(mconn, collectionname, userID)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Hello user"
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

// PEMASUKAN

func GCFHandlerInsertPemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pemasukanResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	err = json.NewDecoder(r.Body).Decode(&pemasukan)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	_, err = InsertPemasukan(mconn, collectionname, pemasukan, userInfo.Id)
	if err != nil {
		pemasukanResponse.Message = err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Insert pemasukan success"
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
	pemasukanResponse.Message = "Get pemasukan success"
	pemasukanResponse.Data = []model.Pemasukan{pemasukan}

	return GCFReturnStruct(pemasukanResponse)
}

func GCFHandlerGetPemasukanFromUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pemasukanResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pemasukanResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pemasukanResponse)
	}

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukan, err := GetPemasukanFromUser(mconn, collectionname, userInfo.Id)
	if err != nil {
		pemasukanResponse.Message = err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Get pemasukan success"
	pemasukanResponse.Data = pemasukan

	return GCFReturnStruct(pemasukanResponse)
}

func GCFHandlerUpdatePemasukan(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
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
	pemasukan.ID = ID

	err = json.NewDecoder(r.Body).Decode(&pemasukan)
	if err != nil {
		pemasukanResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukan, _, err := UpdatePemasukan(mconn, collectionname, pemasukan)
	if err != nil {
		pemasukanResponse.Message = err.Error()
		return GCFReturnStruct(pemasukanResponse)
	}

	pemasukanResponse.Status = true
	pemasukanResponse.Message = "Update pemasukan success"
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
	pemasukanResponse.Message = "Delete pemasukan success"

	return GCFReturnStruct(pemasukanResponse)
}

// PENGELUARAN

func GCFHandlerInsertPengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pengeluaranResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	err = json.NewDecoder(r.Body).Decode(&pengeluaran)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	_, err = InsertPengeluaran(mconn, collectionname, pengeluaran, userInfo.Id)
	if err != nil {
		pengeluaranResponse.Message = err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Insert pengeluaran success"
	pengeluaranResponse.Data = []model.Pengeluaran{pengeluaran}

	return GCFReturnStruct(pengeluaranResponse)
}

func GCFHandlerGetPengeluaranFromUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	pengeluaranResponse.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		pengeluaranResponse.Message = "error parsing application/json1:"
		return GCFReturnStruct(pengeluaranResponse)
	}

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaran, err := GetPengeluaranFromUser(mconn, collectionname, userInfo.Id)
	if err != nil {
		pengeluaranResponse.Message = err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Get pengeluaran success"
	pengeluaranResponse.Data = pengeluaran

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
	pengeluaranResponse.Message = "Get pengeluaran success"
	pengeluaranResponse.Data = []model.Pengeluaran{pengeluaran}

	return GCFReturnStruct(pengeluaranResponse)
}

func GCFHandlerUpdatePengeluaran(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
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
	pengeluaran.ID = ID

	err = json.NewDecoder(r.Body).Decode(&pengeluaran)
	if err != nil {
		pengeluaranResponse.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaran, _, err := UpdatePengeluaran(mconn, collectionname, pengeluaran)
	if err != nil {
		pengeluaranResponse.Message = err.Error()
		return GCFReturnStruct(pengeluaranResponse)
	}

	pengeluaranResponse.Status = true
	pengeluaranResponse.Message = "Update pengeluaran success"
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
	pengeluaranResponse.Message = "Delete pengeluaran success"

	return GCFReturnStruct(pengeluaranResponse)
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
