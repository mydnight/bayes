package main
import (
  "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"fmt"
	"strings"
    "github.com/mydnight/bayes"
	)
//loadIndex loads an index from a database.
func loadIndex(indexName string) bayes.Index {

	myIndex := bayes.Index{}
	database, err := sql.Open("sqlite3", "./bayes.db")
	_, err = database.Exec(
		"CREATE TABLE IF NOT EXISTS Bayes ( id varchar(255), token varchar(255), count REAL, PRIMARY KEY (id, token))")
	if err != nil {
		log.Fatal(err)
	}
	sth, err := database.Query("SELECT token, count FROM Bayes WHERE id=?", indexName)
	if err != nil {
		return myIndex;
	}
	myIndex.M = make(map[string]float64)
	myIndex.Name = indexName
	for sth.Next() {
		var token sql.NullString
		var count sql.NullFloat64
		if err := sth.Scan(&token, &count); err != nil {
			log.Fatal(err)
		}
		myCount := count.Float64
		myToken := token.String
		myIndex.M[myToken] = myCount
		myIndex.Total += myCount
	}
	database.Close()
	return myIndex
	//load index from database
}
func main() {
	fmt.Printf("Starting...\n")
	cat1Index := loadIndex("test1")
	cat2Index := loadIndex("test2")
	cat1Index = bayes.Train(goodIndex, strings.Split("Some Random Text"," "))
	cat2Index = bayes.Train(badIndex, strings.Split("Different Random Text"," "))
	result := bayes.Classify(strings.Split("Text to be classified"," "), cat1Index, cat2Index, 0.005)
	saveIndex("test1",cat1Index);
	saveIndex("test2",cat2Index);
	fmt.Printf(result)
}
//saveIndex saves the index to the database.
func saveIndex(indexName string, myIndex bayes.Index) {
	//save index to database
	database, err := sql.Open("sqlite3", "./bayes.db")
	if err != nil {
		log.Fatal(err)
	}
	for token, count := range myIndex.M {
		database.Exec("INSERT OR REPLACE INTO Bayes(id, token, count) VALUES(?,?,?)", indexName, token, count)
	}
	database.Close()
}
