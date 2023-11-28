package module

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/billblis/billblis_be/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// CRUD
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func InsertManyDocsPemasukan(db *mongo.Database, col string, pemasukan []model.Pemasukan) (insertedIDs []primitive.ObjectID, err error) {
	var interfaces []interface{}
	for _, pemasukan := range pemasukan {
		interfaces = append(interfaces, pemasukan)
	}
	result, err := db.Collection(col).InsertMany(context.Background(), interfaces)
	if err != nil {
		return insertedIDs, fmt.Errorf("kesalahan server: insert")
	}
	for _, id := range result.InsertedIDs {
		insertedIDs = append(insertedIDs, id.(primitive.ObjectID))
	}
	return insertedIDs, nil
}

func InsertManyDocsPengeluaran(db *mongo.Database, col string, pengeluaran []model.Pengeluaran) (insertedIDs []primitive.ObjectID, err error) {
	var interfaces []interface{}
	for _, pengeluaran := range pengeluaran {
		interfaces = append(interfaces, pengeluaran)
	}
	result, err := db.Collection(col).InsertMany(context.Background(), interfaces)
	if err != nil {
		return insertedIDs, fmt.Errorf("kesalahan server: insert")
	}
	for _, id := range result.InsertedIDs {
		insertedIDs = append(insertedIDs, id.(primitive.ObjectID))
	}
	return insertedIDs, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// SIGN UP
func SignUp(db *mongo.Database, col string, insertedDoc model.User) error {
	objectId := primitive.NewObjectID()

	if insertedDoc.Name == "" || insertedDoc.Email == "" || insertedDoc.Password == "" || insertedDoc.MotherName == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}

	hash, _ := HashPassword(insertedDoc.Password)
	// insertedDoc.Password = hash
	user := bson.M{
		"_id":        objectId,
		"email":      insertedDoc.Email,
		"password":   hash,
		"name":       insertedDoc.Name,
		"mothername": insertedDoc.MotherName,
	}
	_, err := InsertOneDoc(db, col, user)
	if err != nil {
		return err
	}
	return nil
}

// SIGN IN
func SignIn(db *mongo.Database, col string, insertedDoc model.User) (user model.User, Status bool, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, false, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, false, fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	if !CheckPasswordHash(insertedDoc.Password, existsDoc.Password) {
		return user, false, fmt.Errorf("password salah")
	}

	return existsDoc, true, nil
}

// SUMBER

func InsertSumber(db *mongo.Database, col string, sumber model.Sumber) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), sumber)
	if err != nil {
		fmt.Printf("InsertSumber: %v\n", err)
		return primitive.NilObjectID, err
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func GetAllSumber(db *mongo.Database) (docs []model.Sumber, err error) {
	collection := db.Collection("sumber")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	err = cursor.All(context.Background(), &docs)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	return docs, nil
}

func GetSumberFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Sumber, err error) {
	collection := db.Collection("sumber")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// PEMASUKAN

func InsertPemasukan(db *mongo.Database, col string, tanggal_masuk string, jumlah_masuk int, id_sumber model.Sumber, deskripsi string, id_user model.User) (insertedID primitive.ObjectID, err error) {
	pemasukan := bson.M{
		"tanggal_masuk": tanggal_masuk,
		"jumlah_masuk":  jumlah_masuk,
		"id_sumber":     id_sumber,
		"deskripsi":     deskripsi,
		"id_user":       id_user,
	}
	result, err := db.Collection(col).InsertOne(context.Background(), pemasukan)
	if err != nil {
		fmt.Printf("InsertPemasukan: %v\n", err)
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func GetAllPemasukan(db *mongo.Database) (pemasukan []model.Pemasukan, err error) {
	collection := db.Collection("pemasukan")
	filter := bson.M{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pemasukan, fmt.Errorf("error GetAllPemasukan mongo: %s", err)
	}

	// Iterate through the cursor and decode each document
	for cursor.Next(context.Background()) {
		var p model.Pemasukan
		if err := cursor.Decode(&p); err != nil {
			return pemasukan, fmt.Errorf("error decoding document: %s", err)
		}
		pemasukan = append(pemasukan, p)
	}

	if err := cursor.Err(); err != nil {
		return pemasukan, fmt.Errorf("error during cursor iteration: %s", err)
	}

	return pemasukan, nil
}

func GetPemasukanFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Pemasukan, err error) {
	collection := db.Collection("pemasukan")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func UpdatePemasukan(db *mongo.Database, doc model.Pemasukan) (err error) {
	filter := bson.M{"_id": doc.ID}
	result, err := db.Collection("pemasukan").UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Printf("UpdatePemasukan: %v\n", err)
		return
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return nil
}

func DeletePemasukan(db *mongo.Database, doc model.Pemasukan) error {
	collection := db.Collection("pemasukan")
	filter := bson.M{"_id": doc.ID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", doc.ID, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", doc.ID)
	}

	return nil
}

// PENGELUARAN

func InsertPengeluaran(db *mongo.Database, col string, tanggal_keluar string, jumlah_keluar int, id_sumber model.Sumber, deskripsi string, id_user model.User) (insertedID primitive.ObjectID, err error) {
	pengeluaran := bson.M{
		"tanggal_keluar": tanggal_keluar,
		"jumlah_keluar":  jumlah_keluar,
		"id_sumber":      id_sumber,
		"deskripsi":      deskripsi,
		"id_user":        id_user,
	}
	result, err := db.Collection(col).InsertOne(context.Background(), pengeluaran)
	if err != nil {
		fmt.Printf("InsertPengeluaran: %v\n", err)
		return
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func GetAllPengeluaran(db *mongo.Database) (pengeluaran []model.Pengeluaran, err error) {
	collection := db.Collection("pengeluaran")
	filter := bson.M{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return pengeluaran, fmt.Errorf("error GetAllPengeluaran mongo: %s", err)
	}

	// Iterate through the cursor and decode each document
	for cursor.Next(context.Background()) {
		var p model.Pengeluaran
		if err := cursor.Decode(&p); err != nil {
			return pengeluaran, fmt.Errorf("error decoding document: %s", err)
		}
		pengeluaran = append(pengeluaran, p)
	}

	if err := cursor.Err(); err != nil {
		return pengeluaran, fmt.Errorf("error during cursor iteration: %s", err)
	}

	return pengeluaran, nil
}

func GetPengeluaranFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.Pengeluaran, err error) {
	collection := db.Collection("pengeluaran")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("_id tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func UpdatePengeluaran(db *mongo.Database, doc model.Pengeluaran) (err error) {
	filter := bson.M{"_id": doc.ID}
	result, err := db.Collection("pengeluaran").UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Printf("UpdatePengeluaran: %v\n", err)
		return
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return nil
}

func DeletePengeluaran(db *mongo.Database, doc model.Pengeluaran) error {
	collection := db.Collection("pengeluaran")
	filter := bson.M{"_id": doc.ID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", doc.ID, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", doc.ID)
	}

	return nil
}

// GET USER
func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetUserFromName(name string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"name": name}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("username tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}
