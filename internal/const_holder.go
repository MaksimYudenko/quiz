package internal

const (
	// paths to appropriate files
	FolderName    = "data"
	QuestionsName = "questions.json"
	ResultsName   = "results.json"
	Questions     = FolderName + "/" + QuestionsName
	Results       = FolderName + "/" + ResultsName
	// default values
	IsRandom = true
	Time     = 30
	// values for testing
	InitialQuestions = `{"2x2 = ":"4","2^10 = ":"1024","pi value (3 digits)":"3.14",
	"Best IT company":"EPAM","Belarus capital":"Minsk","Earth natural satellite":"Moon",
	"Best program language (full name)":"Golang",
	"GoLang is compiled language (1 = true, 0 = false)":"1",
	"Go was publicly announced in":"2009",
	"Are you human being? (yes/no)":"yes"}`
	InitialResult = `{"Alyosha" : 1 ,"Ann" : 4 ,"John" : 3 ,"Albert" : 10 ,
	"Camila" : 3 ,"Vanya" : 2}`
	// formatting
	Introduction = "Hello, let's GO!\nMake your choice, please:\n" +
		"< 1 > to start the quiz\n" +
		"< 2 > to view the rating list\n" +
		"< 3 > to quit"
	ResultHead = "\n----------Top 3 winners----------\n     name\t|\tscore" +
		"\n---------------------------------"
)
