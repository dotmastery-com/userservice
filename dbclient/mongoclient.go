package dbclient

import (
	"UserService/model"

	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoDBHosts = "mongodb"
	AuthDatabase = "sampledb"
	AuthUserName = "user3G3"
	AuthPassword = "skxIfr8Ocn2QxU07"
	//TestDatabase = "sampledb"
	TestDatabase = "store"
)

type Datastore interface {
	Connect()
	AuthenticateUser(user model.User) (bool, error)
	UserExists(user model.User) (bool, error)
	SaveUser(user model.User) error
	GetAllUsers() ([]model.User, error)
}

type MongoClient struct {
	session *mgo.Session
}

// Connect to the Mongo database
func (mc *MongoClient) Connect() {

	var err error

	// We need this object to establish a session to our MongoDB.
	//	mongoDBDialInfo := &mgo.DialInfo{
	//		Addrs:    []string{MongoDBHosts},
	//		Timeout:  60 * time.Second,
	//		Database: AuthDatabase,
	//		Username: AuthUserName,
	//		Password: AuthPassword,
	//	}

	//mc.session, err = mgo.DialWithInfo(mongoDBDialInfo)

	mc.session, err = mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	mc.ensureIndex()

	mc.session.SetMode(mgo.Monotonic, true)

}

func (mc *MongoClient) ensureIndex() {
	session := mc.session.Copy()
	defer session.Close()

	c := session.DB(TestDatabase).C("users")

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func (mc *MongoClient) AuthenticateUser(user model.User) (bool, error) {

	session := mc.session.Copy()
	defer session.Close()

	c := session.DB(TestDatabase).C("users")

	total, error := c.Find(bson.M{"username": user.UserName, "password": user.Password}).Count()

	if total > 0 {
		return true, nil
	}

	return false, error
}

func (mc *MongoClient) UserExists(user model.User) (bool, error) {

	session := mc.session.Copy()
	defer session.Close()

	c := session.DB(TestDatabase).C("users")

	total, error := c.Find(bson.M{"username": user.UserName}).Count()

	if total > 0 {
		return true, nil
	}

	return false, error
}

func (mc *MongoClient) SaveUser(user model.User) error {

	session := mc.session.Copy()
	defer session.Close()

	user.Id = xid.New().String()
	c := session.DB(TestDatabase).C("users")

	err := c.Insert(user)

	return err
}

func (mc *MongoClient) GetAllUsers() ([]model.User, error) {

	session := mc.session.Copy()

	c := session.DB(TestDatabase).C("users")

	var users []model.User
	error := c.Find(bson.M{}).All(&users)

	return users, error

}
