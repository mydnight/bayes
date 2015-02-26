# bayes
Naive Bayes Classifier written in Go.
Accuracy is improved if you feed your data through a stemmer such as http://godoc.org/github.com/surge/porter2 before training
or classifying.
Sample usage:

func main() {
	fmt.Printf("Starting...\n")
	cat1Index := bayes.Index{}
	cat1.Name = "Test1"
	cat2Index := bayes.Index{}
	cat2Index.Name = "Test2"
	goodIndex = bayes.Train(goodIndex, strings.Split("Some chunk of text"," "))	
	badIndex = bayes.Train(badIndex, strings.Split("A different chunk of text"," "))	
	result := bayes.Classify(strings.Split("Text To Classify"," "), goodIndex, badIndex, 0.005)
	fmt.Printf(result)
	fmt.Sprintf("cat1Index.M = %v, cat2Index.M = %v", goodIndex.M, badIndex.M)
}
