package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/web"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

// Router struct would carry the httprouter instance, so its methods could be verwritten and replaced with methds with wraphandler
type Router struct {
	*httprouter.Router
}

// Get is an endpoint to only accept requests of method GET
func (r *Router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

// Post is an endpoint to only accept requests of method POST
func (r *Router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

// Put is an endpoint to only accept requests of method PUT
func (r *Router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

// Delete is an endpoint to only accept requests of method DELETE
func (r *Router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

// NewRouter is a wrapper that makes the httprouter struct a child of the router struct
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	web.TemplateInit()
	config.Init()

	commonHandlers := alice.New(web.LoggingHandler)
	//web.RecoverHandler, context.ClearHandler,
	router := NewRouter()
	log.Println(commonHandlers)

	router.Get("/", commonHandlers.Append(web.AuthCheckerMiddleware).ThenFunc(web.HomeHandler))

	router.Get("/fb_oauth_redirect", commonHandlers.ThenFunc(web.FBOauthRedirectHandler))
	router.Post("/login", commonHandlers.ThenFunc(web.LoginHandler))
	router.Post("/signup", commonHandlers.ThenFunc(web.SignupHandler))
	router.Get("/signup", commonHandlers.Append(web.AuthCheckerMiddleware).ThenFunc(web.SignupPageHandler))

	router.Get("/archive", commonHandlers.Append(web.FrontAuthHandler).ThenFunc(web.ItemsArchiveHandler))
	router.Get("/notifications", commonHandlers.Append(web.FrontAuthHandler).ThenFunc(web.NotificationsHandler))
	router.Get("/search", commonHandlers.Append(web.AuthCheckerMiddleware).ThenFunc(web.SearchHandler))
	// router.Get("/google_oauth_redirect", commonHandlers.ThenFunc(web.OauthRedirectHandler))

	router.Get("/profile", commonHandlers.Append(web.FrontAuthHandler).ThenFunc(web.ProfileHandler))
	router.Post("/profile", commonHandlers.Append(web.FrontAuthHandler).ThenFunc(web.SaveNewItemHandler))

	//router.ServeFiles("/static/*filepath", http.Dir("./build/static"))

	fileServer := http.FileServer(http.Dir("./web/templates/assets"))
	router.GET("/public/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		fileServer.ServeHTTP(w, r)
	})

	// router.NotFound = commonHandlers.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./build/index.html")
	// })

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No Global port has been defined, using default")
		PORT = "8080"
	}

	handler := cors.New(cors.Options{
		//		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	}).Handler(router)
	log.Println("serving ")
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
