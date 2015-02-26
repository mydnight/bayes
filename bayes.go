//Package bayes implements a naive Bayes classifier.  Trained models are
//stored in a SQLite database.
package bayes

/*Takes a text sample (stemmed state doesn't matter and is only used for
 *improving accuracy) and either trains the system based on it, or classifies
 *it based on previously trained data.*/

import (
	"sync"
	"math"
)

type Index struct {
	sync.RWMutex
	M map[string]float64
	Name string
	Total float64
}

//Classify takes a text sample and classifies it as one of two categories.
func Classify(text []string, cat1Index Index, cat2Index Index, tolerance float64) string {
	cat1Index.RLock()
	cat2Index.RLock()
	i := 0.00
	invi := 0.00
	for _, token := range text {
		prob := calcProbability(cat1Index.M[token], cat1Index.Total, cat2Index.M[token], cat2Index.Total)
		if i > 0.00  && prob > 0.00{
			i = i * prob
		} else if i == 0.00 {
			i = prob
		}
		if invi != 0.00 {
			invi = invi * (1.00 - prob)
		} else {
			invi = 1.00 - prob
		}
	}
	cat1Index.RUnlock()
	cat2Index.RUnlock()
	totprob := i / (i + invi)
	if totprob > (0.5 + tolerance) {
		return cat1Index.Name 
	} else if totprob < (0.5 - tolerance) {
		return cat2Index.Name
	}
	return "Unknown"
}

//calcProbability calculates the probability any given word is part of either category
func calcProbability(cat1count float64, cat1total float64, cat2count float64, cat2total float64) float64 {
	cat1Percent:= math.Log(cat1count / cat1total)
	cat2Percent:= math.Log(cat2count / cat2total)
	probWord := cat1Percent / (cat1Percent + cat2Percent)
	if math.IsNaN(probWord) { //Word doesn't exist in either of the indices
		probWord = 0.0
	}
	return probWord
}

//Train trains a category with a text sample.
func Train(myIndex Index, text []string) Index {
	//myIndex := loadIndex(trainedIndex)
	myIndex.Lock()
	for _, token := range text {
		myIndex.M[token]++
		myIndex.Total++
	}
	myIndex.Unlock()
	return myIndex
}
