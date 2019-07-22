package main

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var globalSession *mgo.Session

// InitMgo ...
func InitMgo() error {
	var err error
	globalSession, err = mgo.Dial("127.0.0.1:27017")
	if err != nil {
		return errors.New("mgo dial fail")
	}
	globalSession.SetMode(mgo.Monotonic, true)
	globalSession.SetPoolLimit(100)
	return nil
}

// GetMgoC ...
func GetMgoC(sess *mgo.Session, db, c string) *mgo.Collection {
	return sess.DB(db).C(c)
}

// GetMgoS ...
func GetMgoS() *mgo.Session {
	return globalSession.Copy()
}

// CloseMgo ...
func CloseMgo() {
	globalSession.Close()
}
