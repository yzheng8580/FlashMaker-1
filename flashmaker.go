package flashmaker

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"appengine"
	"appengine/urlfetch"
)

var (
	flashmaker []byte
	foundDefs  = template.Must(template.ParseFiles("result.html"))
	fdictURL   = "http://dictionary.reference.com/browse/"
	edictURL   = "?s=t"
)

func init() {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	flashmaker = index
	http.HandleFunc("/", root)
	http.HandleFunc("/define", define)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write(flashmaker)
}

func define(w http.ResponseWriter, r *http.Request) {
	words := strings.Split(r.FormValue("content"), "\n")

	err := foundDefs.Execute(w, "search result")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func search(w http.ResponseWriter, r *http.Request, word string) string {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get(fdictURL + word + edictURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "Error"
	}
	var buffer bytes.Buffer
	var num int
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	text := strings.Split(string(body), "def-number")
	var content, class, div []string
	for i := 1; i < 3; i++ {
		content = strings.Split(text[i], "<div class=\"def-content\">")
		for j := 0; j < len(content); j++ {
			class = strings.Split(content[j], ": <div class")
			for k := 0; k < len(class); k++ {
				div = strings.Split(class[k], "</div>")
				for x := 0; x < len(div); x++ {
					s := div[x]
					s = strings.TrimSpace(s)
					if len(s) > 0 {
						var first = s[0]
						var upper = first >= 65 && first <= 90
						var lower = first >= 97 && first <= 122
						var parenth = first == 40
						if upper || lower || parenth {
							num++
							buffer.WriteString(strconv.Itoa(num) + ". " + s + "\n")
						}
					}
				}
			}

		}
	}
	return buffer.String()
}
