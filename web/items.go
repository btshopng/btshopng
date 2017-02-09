package web

import (
	"log"
	"net/http"

	"time"

	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/models"
	uuid "github.com/satori/go.uuid"
)

type Data struct {
	User    models.User
	Barters []models.Barter
}

func NewItemHandler(w http.ResponseWriter, r *http.Request) {

	user, err := Userget(r)
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}
	//log.Println(user)

	user.FormattedDateCreated = user.DateCreated.Format("Mon, 02 Jan 2006")

	result, err := user.Get(config.GetConf())
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}
	data := Data{User: result}

	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, data)
}

func SaveNewItemHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post data from the request.
	r.ParseForm()

	user, err := Userget(r)
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}

	result, err := user.Get(config.GetConf())
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}

	have := r.FormValue("have")
	haveCat := r.FormValue("haveCat")
	need := r.FormValue("need")
	needCat := r.FormValue("needCat")
	location := r.FormValue("location")

	if have == "" || haveCat == "" || need == "" || needCat == "" || location == "" {
		http.Redirect(w, r, "/newitem?newerror=Fill+out+all+fields", 301)
		return
	}

	uniqueID := uuid.NewV1().String()
	// create a barter model....
	barter := models.Barter{
		ID:           uniqueID,
		UserID:       result.ID,
		Have:         have,
		HaveCategory: haveCat,
		Need:         need,
		NeedCategory: needCat,
		Location:     location,
		DateCreated:  time.Now(),
		Status:       true,
		Images:       []string{"", "", ""},
	}

	err = barter.Upsert(config.GetConf())
	if err != nil {
		http.Redirect(w, r, "/newitem?error=Could+not+save+your+barter", 301)
		return
	}
	log.Println("New barter added")
	// send a notification to the user that the barter has been added.
	http.Redirect(w, r, "/newitem", 301)
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {

	user, err := Userget(r)
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}
	//log.Println(user)

	user.FormattedDateCreated = user.DateCreated.Format("Mon, 02 Jan 2006")

	result, err := user.Get(config.GetConf())
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}

	// Supply UserID to be used for retrieving all barters.
	barter := models.Barter{UserID: result.ID}

	data := Data{User: result}

	// data := struct {
	// 	User    models.User
	// 	Barters []models.Barter
	// }{
	// 	User: user,
	// }

	data.Barters, err = barter.GetAll(config.GetConf())
	if err != nil {
		log.Println("No barter for this user.")
	}

	tmp := GetTemplates().Lookup("profile_barter_archive.html")
	tmp.Execute(w, data)
}