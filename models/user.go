package models

import (
	// mgo "gopkg.in/mgo.v2"
	"log"
	"time"

	"github.com/btshopng/btshopng/config"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	FBPicture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	Link                 string `json:"link"`
	DateCreated          time.Time
	FormattedDateCreated string
	Password             []byte
}

func (user User) Upsert(c *config.Conf) error {

	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.USERCOLLECTION).With(mgoSession)

	_, err := collection.Upsert(bson.M{
		"$or": []bson.M{
			bson.M{
				"id": user.ID,
			},
			bson.M{
				"email": user.Email,
			},
		},
	}, user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (user User) Get(c *config.Conf) (User, error) {

	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.USERCOLLECTION).With(mgoSession)

	result := User{}
	err := collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{
				"id": user.ID,
			},
			bson.M{
				"email": user.Email,
			},
		},
	}).One(&result)

	if err != nil {
		log.Println(err)
		return user, err
	}
	return user, nil
}

// CheckUser checks if a user exists in the database
func (user User) CheckUser(c *config.Conf) (User, error) {

	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.USERCOLLECTION).With(mgoSession)

	result := User{}
	// log.Println("user.Email: ", user.Email)

	err := collection.Find(bson.M{"email": user.Email}).One(&result)
	if err != nil {
		log.Println("User not found:", err)
		// log.Println("error User: ", user)

		// return result, errors.New("Username or password is incorrect (no user)")
		return result, err
	}

	// log.Println("USER: ", user, "result:", result)
	// if result.Name == "" {
	// 	return result, errors.New("Username or Password is incorrect(no user)")
	// }

	return result, nil
}

// Insert : inserts user data to the DB
func (user User) Insert(c *config.Conf) error {

	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.USERCOLLECTION).With(mgoSession)

	err := collection.Insert(user)
	if err != nil {
		log.Println("Could not insert into the DB")
		return err
	}

	return nil
}
