//Package bayes implements a naive Bayes classifier.  Trained models are
//stored in a SQLite database.
package bayes

/*Takes a text sample (stemmed state doesn't matter and is only used for
 *improving accuracy) and either trains the system based on it, or classifies
 *it based on previously trained data.*/

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

//Tolerance := 0.05 //Tolerance -- i.e. if the result is 0.5 +- tolerance value, return "undetermined" category

type index struct {
	sync.RWMutex
	m map[string]float64
}

//Classify takes a text sample and classifies it as one of two categories.
func Classify(text []string, cat1Index string, cat2Index string, tolerance float64) int {
	myCat1Index := loadIndex(cat1Index)
	myCat2Index := loadIndex(cat2Index)
	myCat1Index.RLock()
	myCat2Index.RLock()
	//iterate over goodindex and badindex, and sum the values.
	cat1total := 0.00
	cat2total := 0.00
	for _, value := range myCat1Index.m {
		cat1total += value
	}
	for _, value := range myCat2Index.m {
		cat2total += value
	}
	i := 0.00
	invi := 0.00
	for _, token := range text {
		cat1token := myCat1Index.m[token]
		cat2token := myCat2Index.m[token]
		prob := calcProbability(cat1token, cat1total, cat2token, cat2total)
		if i != 0.00 {
			i = i * prob
		} else {
			i = prob
		}
		if invi != 0.00 {
			invi = invi * (1.00 - prob)
		} else {
			invi = 1.00 - prob
		}
	}
	myCat1Index.RUnlock()
	myCat2Index.RUnlock()
	totprob := i / (i + invi)
	if totprob > (0.5 + tolerance) {
		return 1
	} else if totprob < (0.5 - tolerance) {
		return -1
	}
	return 0
}

//calcProbability calculates the probability any given word is part of either category
func calcProbability(cat1count float64, cat1total float64, cat2count float64, cat2total float64) float64 {
	bw := cat1count / cat1total
	gw := cat2count / cat2total
	pw := ((bw) / ((bw) + (gw)))
	s := 1.00
	x := 0.5
	n := cat1count + cat2count
	fw := ((s * x) + (n * pw)) / (s + n)
	return fw

}

//Train trains a category with a text sample.
func Train(trainedIndex string, text []string) {
	myIndex := loadIndex(trainedIndex)
	myIndex.Lock()
	for _, token := range text {
		myIndex.m[token]++
	}
	saveIndex(trainedIndex, myIndex)
	myIndex.Unlock()
}

//loadIndex loads an index from a database.
func loadIndex(indexName string) index {

	myIndex := index{}
	database, err := sql.Open("sqlite3", "./bayes.db")
	_, err = database.Exec(
		"CREATE TABLE IF NOT EXISTS Bayes ( id varchar(255), token varchar(255), count REAL, PRIMARY KEY (id, token))")
	if err != nil {
		log.Fatal(err)
	}
	sth, err := database.Query("SELECT token, count FROM Bayes WHERE id=?", indexName)
	myIndex.m = make(map[string]float64)
	for sth.Next() {
		var token sql.NullString
		var count sql.NullFloat64
		if err := sth.Scan(&token, &count); err != nil {
			log.Fatal(err)
		}
		myCount := count.Float64
		myToken := token.String
		myIndex.m[myToken] = myCount
	}
	database.Close()
	return myIndex
	//load index from database
}

//saveIndex saves the index to the database.
func saveIndex(indexName string, myIndex index) {
	//save index to database
	database, err := sql.Open("sqlite3", "./rss2.db")
	if err != nil {
		log.Fatal(err)
	}
	for token, count := range myIndex.m {
		database.Exec("INSERT OR REPLACE INTO Bayes(id, token, count) VALUES(?,?,?)", indexName, token, count)
	}
	database.Close()
}
