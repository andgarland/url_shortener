package database

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/andgarland/url_shortener/encode"
    "github.com/andgarland/url_shortener/config"
)

//Global database variable
var (
    DB *sql.DB
)

//Function to open a connection to the database and if necessary create the required table
func GetDB(user, password, database string) (*sql.DB, error) {

    dataSource := fmt.Sprintf("%s%s%s%s%s%s%s", user, ":", password, "@tcp(", config.Settings.Host, ")/", database)

    DB, err := sql.Open("mysql", dataSource)

    if err != nil {
		return nil, err
	}

    err = DB.Ping()
    if err != nil {
        return nil, err
    }

    tableQuery := "CREATE TABLE IF NOT EXISTS urls (id serial NOT NULL, longURL VARCHAR (10000));"

	_, err = DB.Exec(tableQuery)
	if err != nil {
		return nil, err
	}

    return DB, err
}

//Function that stores submitted URLs in the database and fetches the respective short URLs
func GetShortURL(longURL string) (string, error) {

    res, err := DB.Exec("INSERT INTO urls(longURL) VALUES(?)", longURL)
    if err != nil {
		return "", err
	}

    id, err := res.LastInsertId()
    if err != nil {
		return "", err
	}

    shortURL := encode.Encode(id)

    return shortURL, nil
}

//Function that fetches the original URLs from the database when given its short equivalent
func GetLongURL(shortURL string) (string, error) {
    var longURL string

    id := encode.Decode(shortURL)

    err := DB.QueryRow("SELECT longURL FROM urls WHERE id=?", id).Scan(&longURL)
    if err != nil {
		return "", err
	}

    return longURL, nil
}