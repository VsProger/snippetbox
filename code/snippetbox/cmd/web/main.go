package main

import (
	"database/sql"
	"flag"
	"github.com/VsProger/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() { // 60 44
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
