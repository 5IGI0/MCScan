package main

import (
	"bufio"
	"database/sql"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var runtime struct {
	db             *sqlx.DB
	serverTemplate *template.Template
	searchTemplate *template.Template
	defaultFavicon []byte
}

func main() {
	var err error

	templateBase, err := template.New("base").
		Funcs(template.FuncMap{
			"join":          strings.Join,
			"time_relative": timestamp_to_relative,
			"u32greater":    func(a uint32, b uint32) bool { return a > b }}).
		ParseFiles("templates/base.html")
	assert_err(err)
	runtime.serverTemplate, err = template.Must(templateBase.Clone()).ParseFiles("templates/server.html")
	assert_err(err)
	runtime.searchTemplate, err = template.Must(templateBase.Clone()).ParseFiles("templates/search.html")
	assert_err(err)
	runtime.defaultFavicon, err = ioutil.ReadFile("static/pack.png")
	assert_err(err)

	runtime.db = sqlx.MustOpen("mysql", os.Getenv("MYSQL_DB_URL"))
	runtime.db.SetMaxOpenConns(10)

	go scan_task()

	router := mux.NewRouter()
	router.HandleFunc("/v1/add_servers", add_handler).Methods("POST")
	router.HandleFunc("/v1/favicon/{id:[0-9]+}", get_favicon)
	router.HandleFunc("/v1/favicon/{id:[0-9]+}.png", get_favicon)
	router.HandleFunc("/v1/server/{id:[0-9]+}", ApiJsonDecorator(ApiGetServerById))
	router.HandleFunc("/v1/server/{addr}", ApiJsonDecorator(ApiGetServerByAddr))
	router.HandleFunc("/v1/search", ApiJsonDecorator(ApiSearch))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/server/{id:[0-9]+}", html_server)
	router.HandleFunc("/search", html_search)
	router.HandleFunc("/", html_search)

	if strings.HasPrefix(os.Getenv("LISTEN_ADDR"), "unix:") {
		ListenAndServeUnix(os.Getenv("LISTEN_ADDR")[5:], router)
	} else {
		http.ListenAndServe(os.Getenv("LISTEN_ADDR"), router)
	}
}

func ListenAndServeUnix(addr string, handler http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	ln, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func scan_task() {
	for {
		rows, err := runtime.db.Queryx("SELECT id, address FROM `servers` WHERE last_scan<? ORDER BY `last_scan` ASC LIMIT 100",
			get_timestamp()-SCAN_INTERVAL)
		assert_err(err)

		for rows.Next() {
			var tmp struct {
				Addr string `db:"address"`
				Id   int64  `db:"id"`
			}
			assert_err(rows.StructScan(&tmp))
			add_server(tmp.Addr, tmp.Id)
		}

		time.Sleep(time.Minute * 1)
	}
}

func html_search(w http.ResponseWriter, r *http.Request) {
	err, ret := ApiSearch(w, r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	err = runtime.searchTemplate.ExecuteTemplate(w, "base", ret)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
}

func html_server(w http.ResponseWriter, r *http.Request) {
	srv_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err, ret := InternalGetServerById(srv_id, true)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	err = runtime.serverTemplate.ExecuteTemplate(w, "base", ret)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
}

func get_favicon(w http.ResponseWriter, r *http.Request) {
	var favicon []byte
	err := runtime.db.Get(&favicon, "SELECT raw_favicon FROM favicons WHERE id=?", mux.Vars(r)["id"])

	if err == sql.ErrNoRows {
		w.Header().Add("Content-Type", "image/png")
		w.WriteHeader(404)
		w.Write(runtime.defaultFavicon)
		return
	}

	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.WriteHeader(200)
	w.Write(favicon)
}

func add_handler(w http.ResponseWriter, r *http.Request) {
	scanner := bufio.NewScanner(r.Body)

	scanner.Scan()
	if scanner.Text() != os.Getenv("ADD_SERVER_PASSWORD") {
		w.WriteHeader(501)
		w.Write([]byte("wrong password"))
		return
	}

	for scanner.Scan() {
		add_server(scanner.Text(), 0)
	}
}
