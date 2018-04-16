package handler

import (
    "net/http"
    "html/template"
    "regexp"
    "github.com/andgarland/url_shortener/database"
)

//Function that creates the request multiplexer. Matches URLs to patterns.
func Handlers() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", index)
    mux.HandleFunc("/short.url/", forward)
    return mux
}

//Function that renders the URL shortener homepage
func index(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        render(w, "templates/index.html", nil)
    } else {
        submit(w, r)
    }
}

//Function that identifies invalid submissions and accepts valid URLs
func submit(w http.ResponseWriter, r *http.Request) {
    longURL := r.FormValue("url")

    if validate(longURL) == false {
        render(w, "templates/index.html", true)
    } else {

        //Stores original URLs in database and generates short URL
        shortURL, err := database.GetShortURL(longURL)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        //Renders webpage displaying the short URL
        render(w, "templates/shortened.html", shortURL)
    }
}

//Function that redirects the short URL to the actual URL
func forward(w http.ResponseWriter, r *http.Request) {
    shortURL := r.URL.Path[len("/short.url/"):]

    //Fetches the original URL from the database
    longURL, err := database.GetLongURL(shortURL)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found"))
		return
    }

    http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

//Function that checks for valid URL submissions via regular expressions
func validate(url string) bool {
    re := regexp.MustCompile("htt(p|ps)://[[:word:]]+\\..+")
    matched := re.Match([]byte(url))
    return matched
}

//Function that renders the webpage templates and passes on the relevant data to them
func render(w http.ResponseWriter, filename string, data interface{}) {
    tmpl, err := template.ParseFiles(filename)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}