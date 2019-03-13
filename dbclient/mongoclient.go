package dbclient

import (
	"UserService/model"
	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoClient struct {
	session *mgo.Session
}

func (mc *MongoClient) Connect() {

	var err error
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

	c := session.DB("store").C("users")

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

	c := session.DB("store").C("users")

	total, error := c.Find(bson.M{"username": user.UserName, "password": user.Password}).Count()

	if total > 0 {
		return true, nil
	}

	return false, error
}

func (mc *MongoClient) UserExists(user model.User) (bool, error) {

	session := mc.session.Copy()
	defer session.Close()

	c := session.DB("store").C("users")

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
	c := session.DB("store").C("users")

	err := c.Insert(user)

	return err
}


func (mc *MongoClient) GetAllUsers()([]model.User, error) {

	session := mc.session.Copy()

	c := session.DB("store").C("users")

	var users []model.User
	error := c.Find(bson.M{}).All(&users)

	return users, error

}
