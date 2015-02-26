# bayes
Naive Bayes Classifier written in Go.
Accuracy is improved if you feed your data through a stemmer such as http://godoc.org/github.com/surge/porter2 before training
or classifying.  A sample driver program with the ability to load and save models from a SQLite database is in the samples directory.
