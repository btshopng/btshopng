package main

import (
	"fmt"
	"html/template"
	"net/http"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// template of signup/signin page to be served
var tmpl = template.Must(template.New("signin_signup.html").ParseFiles("../signin_signup.html"))
var profile = template.Must(template.New("profile_notifications.html").ParseFiles("../profile_notifications.html"))

// User interface to contain user information
type User struct {
	Name     string
	Email    string
	Password string
}

// Data contains data to be passed to templates
type Data struct {
	SignupMessage string
	LoginMessage  string
}

func main() {

	// Serve the signup page.
	http.HandleFunc("/signup", signupHandler)
	// process login form
	http.HandleFunc("/login", loginHandler)
	// Serve profile page
	http.HandleFunc("/profile", profileHandler)

	// Serve static files.
	f := http.FileServer(http.Dir("../"))
	http.Handle("/public/", http.StripPrefix("/public/", f))

	// Start up the server
	port := "3000"
	fmt.Printf("Server started on %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

//////-- HANDLERS --//////

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// instantiate data
		data := Data{}

		// create mongo session
		session, err := mgo.Dial("mongodb://spankie:506dad@ds163738.mlab.com:63738/btshopng")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// close the session
		defer session.Close()

		session.SetSafe(&mgo.Safe{})

		// get the form values
		r.ParseForm()

		// TODO: Validate form data before checking the db.(ifempty, matches the required format)
		email := r.PostFormValue("email")
		// TODO: password would be hashed to match the DB result
		passwd := r.PostFormValue("password")

		// select Collection
		c := session.DB("btshopng").C("Users")

		// result struct
		result := bson.D{}

		// Check if email and password matches any in the DB
		err = c.Find(bson.M{"Email": email, "Password": passwd}).One(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// debug purposes
		fmt.Println("result:", result.Map()["Email"])

		if len(result) == 0 {

			data.LoginMessage = "Username or Password is incorrect"
			tmpl.Execute(w, data)

		} else {

			// TODO: redirect to the the users profile page with session
			http.Redirect(w, r, "/profile", http.StatusFound)

		}
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {

	// instantiate data
	data := Data{}

	// check if the request is a post request.
	if r.Method == "POST" {
		// If it is a post request, process the request.

		// connect to mongodb server
		session, err := mgo.Dial("mongodb://spankie:506dad@ds163738.mlab.com:63738/btshopng")
		if err != nil {
			// reply with internal server error
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// close the connection lastly.
		defer session.Close()

		// set session to safe.
		session.SetSafe(&mgo.Safe{})

		// get the post form values
		r.ParseForm()
		// new struct containing the form values
		newUser := User{
			Name:  r.PostFormValue("name"),
			Email: r.PostFormValue("email"),
			// TODO: Password would be hashed before storage
			Password: r.PostFormValue("passwd"),
		}

		// create a db connection
		c := session.DB("btshopng").C("Users")

		// check if email has already been used by querying the db
		var count int
		count, err = c.Find(bson.M{"Email": string(newUser.Email)}).Count()
		// Select(bson.M{"email": 0}).
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// compare the result from the db to an empty struct
		if count <= 0 {

			// if result is empty, then the email can be used.
			fmt.Println("Email is available")

			// insert the User data to db
			err = c.Insert(newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			// Execute template
			data.SignupMessage = "A verification email has been sent to you"
			tmpl.Execute(w, data)

		} else {

			// if result is not empty, then the email has already been used before.
			fmt.Println("email already exists")
			data.SignupMessage = "Email already exists"
			tmpl.Execute(w, data)
		}

	} else if r.Method == "GET" {
		// if the request is not a post request, just Serve the page
		tmpl.Execute(w, data)
	}

}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: establish session at login or here.

	data := struct {
		Username string
		Time     time.Time
	}{
		Username: "Dummy Name",
		Time:     time.Now(),
	}

	profile.Execute(w, data)
}