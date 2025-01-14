package api

import (
	"book/internal/delivery/http/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

// todo возможно стоит создать новую структуру для router

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	MethodFunc(method string, pattern string, handlerFn http.HandlerFunc)
	URLParam(r *http.Request, key string) string
}

type router struct {
	mux *chi.Mux
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *router) MethodFunc(method string, pattern string, handlerFn http.HandlerFunc) {
	r.mux.MethodFunc(method, pattern, handlerFn)
}

// todo delete
func (r *router) URLParam(req *http.Request, key string) string {
	return chi.URLParam(req, key)
}

func NewHTTPRouter(handlers handlers.Handlers) Router {
	r := &router{mux: chi.NewRouter()}
	r.initRoutes(handlers)
	// todo delete it
	r.mux.Mount("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))
	r.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	return r
}

func (r *router) initRoutes(handlers handlers.Handlers) {
	r.MethodFunc("GET", "/notFound", handlers.NotFound)

	r.MethodFunc("GET", "/book/{uuid}", handlers.Book)
	r.MethodFunc("GET", "/books/", handlers.Books)

	r.MethodFunc("POST", "/books", handlers.CreateBook)
	r.MethodFunc("PUT", "/books", handlers.UpdateBook)
	r.MethodFunc("DELETE", "/books/{title}", handlers.DeleteBook)

	r.MethodFunc("GET", "/read/{bookUUID}", handlers.ReadBook)
	r.MethodFunc("GET", "/download/{bookUUID}", handlers.DownloadBook)

	r.MethodFunc("GET", "/authors/{name}/{surname}/{patronymic}", handlers.Author)
	r.MethodFunc("POST", "/authors", handlers.CreateAuthor)
	r.MethodFunc("DELETE", "/authors/{name}/{surname}/{patronymic}", handlers.DeleteAuthor)
	//todo refactor if work
	//r.MethodFunc("GET", "/web/static/css/{name}", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("/web/static/css/" + r.URL.Path[len("/web/static/css/"):])
	//	http.ServeFile(w, r, "/web/static/css/"+r.URL.Path[len("/web/static/css/"):])
	//})
	//router.MethodFunc(w, r)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//router.MethodFunc("GET", "/", NotFound)
}

//func NotFound(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(http.StatusNotFound)
//	w.Write([]byte("Not found hih"))
//}
