package billblis

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "yellow12"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("private key")
	fmt.Println(privateKey)
	fmt.Println("public key")
	fmt.Println(publicKey)
	fmt.Println("hasil")
	hasil, err := watoken.Encode("Billblis", privateKey)
	fmt.Println(hasil, err)
}

func TestValidateToken(t *testing.T) {
	tokenstring := "Bearer v4.public.eyJleHAiOiIyMDIzLTExLTAyVDA3OjA4OjMzWiIsImlhdCI6IjIwMjMtMTEtMDJUMDU6MDg6MzNaIiwiaWQiOiJlcmRpdG9AZ21haWwuY29tIiwibmJmIjoiMjAyMy0xMS0wMlQwNTowODozM1oifTzDGPfgi7dYoY_MhpsfndX6mR7srxjJK98J3S_PzeCMAWjDmSngxhxMFpCCrq8zdIu2xP1ziBGQ34oPUo04KAw" // Gantilah dengan token PASETO yang sesuai
	tokenstring = strings.TrimPrefix(tokenstring, "Bearer ")
	publicKey := "f48bd58cb3b3972d05bb9303b15ce9b83f4fcb9c871d1b05906f2fec20620ea0"
	payload, _err := watoken.Decode(publicKey, tokenstring)
	if _err != nil {
		fmt.Println("expire token", _err)
	} else {
		fmt.Println(payload)
		// fmt.Println(payload.Nbf)
		// fmt.Println(payload.Iat)
		fmt.Println(payload.Exp)
	}

}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "billblis")

	var userdata User
	userdata.Username = "renjun"
	userdata.Password = "yellow12"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "billblis")
	var userdata User
	userdata.Username = "renjun"
	userdata.Password = "yellow12"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "billblis")
	var userdata User
	userdata.Username = "renjun2"
	userdata.Password = "yellow12"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestGCFPostHandler(t *testing.T) {

	// Membuat body request sebagai string
	requestBody := `{"username": "renjun", "password": "yellow12"}`

	// Membuat objek http.Request
	r := httptest.NewRequest("POST", "https://contoh.com/path", strings.NewReader(requestBody))
	r.Header.Set("Content-Type", "application/json")

	resp := GCFPostHandler("PASETOPRIVATEKEY", "MONGOSTRING", "billblis", "user", r)
	fmt.Println(resp)
}
