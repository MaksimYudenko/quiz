package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	in "github.com/MaksimYudenko/quiz/internal"
)

var (
	wg                 sync.WaitGroup
	QuestionsPath      string
	ResultsPath        string
	IsRandomize        bool
	Duration           int
	qNum               int
	answers, responses map[int]string
)

type result struct {
	name  string
	score int
}

func init() {
	flag.StringVar(&QuestionsPath, "questions_file", in.Questions, "path/to/questions_file")
	flag.StringVar(&ResultsPath, "results_file", in.Results, "path/to/results_file")
	flag.BoolVar(&IsRandomize, "random", in.IsRandom, "randomize order of questions")
	flag.IntVar(&Duration, "time", in.Time, "test duration, sec")
	flag.Parse()
	if er := PopulateTestData(); er != nil {
		log.Fatal(er)
	}
}

func main() {
	// ---------------------------------obtaining an initial data--------------------------------
	data, er := GetFileData(QuestionsPath)
	if er != nil {
		log.Fatalln(fmt.Errorf("something went wrong, look at the stack:\n"+
			"error in Start() while performing GetFileData():\n%v", er))
	}

	qNum = len(data)
	questions := make(map[int]string, qNum)
	answers = make(map[int]string, qNum)
	responses = make(map[int]string, qNum)

	i := 0
	for q, a := range data {
		questions[i] = q
		answers[i] = a.(string)
		i++
	}
	// -------------------------------------starting the quiz--------------------------------------
	er = makeChoice()
	if er != nil {
		log.Fatalln(fmt.Errorf("something went wrong, look at the stack:\n"+
			"error in Start() while performing makeChoice():\n%v", er))
	}

	respondTo := make(chan string)
	if IsRandomize {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	randPool := rand.Perm(qNum)

	wg.Add(1)
	timeUp := time.After(time.Second * time.Duration(Duration))
	go func() {
		defer wg.Done()
	label:
		for i := 0; i < qNum; i++ {
			j := randPool[i]
			go ask(os.Stdout, os.Stdin, questions[j], respondTo)
			select {
			case <-timeUp:
				_, _ = fmt.Fprintln(os.Stderr, "\n\nTime is up, thank you.")
				break label
			case ans, ok := <-respondTo:
				if ok {
					responses[j] = ans
				} else {
					break label
				}
			}
		}
	}()
	wg.Wait()
	// -------------------------finishing the quiz and saving results -------------------------
	score := finish()
	er = retainUser(score)

	if er != nil {
		log.Fatalln(fmt.Errorf("something went wrong, look at the stack:\n"+
			"error in Start() while performing retainUser():\n%v", er))
	}
}

func PopulateTestData() (er error) {
	er = in.Create(in.FolderName)
	if er != nil {
		return fmt.Errorf("error in PopulateTestData() while creating %s:\n%v",
			in.FolderName, er)
	}
	er = in.Create(in.Questions)
	if er != nil {
		return fmt.Errorf("error in PopulateTestData() while creating %s:\n%v",
			in.Questions, er)
	}
	er = in.Create(in.Results)
	if er != nil {
		return fmt.Errorf("error in PopulateTestData() while creating %s:\n%v",
			in.Results, er)
	}
	return
}

func GetFileData(filePath string) (map[string]interface{}, error) {
	bytes, er := in.Read(filePath)
	if er != nil {
		return nil, fmt.Errorf("error in GetFileData() while reading the file:\n%v", er)
	}
	data := make(map[string]interface{})
	er = json.Unmarshal(bytes, &data)
	if er != nil {
		return nil, fmt.Errorf("error in GetFileData() while unmarshalling:\n%v", er)
	}
	return data, er
}

func makeChoice() error {
	fmt.Println(in.Introduction)
	var choice int
	_, er := fmt.Fscan(os.Stdin, &choice)
	if er != nil {
		return fmt.Errorf("error in makeChoice() while scanning user choice:\n%v", er)
	}

	switch choice {
	case 1:
		return nil
	case 2:
		data, er := GetFileData(ResultsPath)
		if er != nil {
			return fmt.Errorf("error in makeChoice() while fetching data:\n%v", er)
		}
		res := getSortedResults(data)
		fmt.Println(in.ResultHead)
		for i := 0; i < len(res); i++ {
			fmt.Println("   ", res[i].name, "\t|\t", res[i].score)
		}
		os.Exit(0)
	case 3:
		fmt.Println("Good bye.")
		os.Exit(0)
	default:
		fmt.Println("Check your choice, please.")
		er = makeChoice()
	}
	return er
}

func ask(w io.Writer, r io.Reader, question string, replyTo chan string) {
	reader := bufio.NewReader(r)
	_, _ = fmt.Fprintln(w, "Question:\t"+question)
	_, _ = fmt.Fprint(w, "Answer:\t\t")

	answer, er := reader.ReadString('\n')
	if er != nil {
		close(replyTo)
		log.Fatalln(fmt.Errorf("error while reading user input:\n%v", er))
	}
	replyTo <- answer
}

func finish() int {
	correct := 0
	for i := 0; i < qNum; i++ {
		if strings.EqualFold(
			strings.TrimSpace(answers[i]), strings.TrimSpace(responses[i])) {
			correct++
		}
	}
	_, _ = fmt.Fprintf(os.Stdout,
		"\nYou've answered %d/%d questions correctly.\n\n",
		correct, qNum)
	return correct
}

func retainUser(score int) error {
	data, er := GetFileData(ResultsPath)
	if er != nil {
		return fmt.Errorf("error in retainUser() while performing GetFileData():\n%v", er)
	}
	results := getSortedResults(data)
	if score >= results[len(results)-1].score {
		fmt.Println("To store your progress - type your name:")
		var name string
		_, er := fmt.Fscan(os.Stdin, &name)
		if er != nil {
			return fmt.Errorf("error in retainUser() while scanning name:\n%v", er)
		}
		er = addUser(name, score)
		if er != nil {
			return fmt.Errorf("error in retainUser() while performing addUser():\n%v", er)
		}
		fmt.Println("You've reached at the top 3, congratulations!")
	}
	return nil
}

func getSortedResults(m map[string]interface{}) (res []result) {
	var ss []result
	for key, value := range m {
		num, ok := value.(float64)
		if ok {
			ss = append(ss, result{key, int(num)})
		} else {
			er := errors.New("sorting error: map value has unknown number format")
			log.Fatalln("error in getSortedResults() while sorting:\n", er)
		}
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].score > ss[j].score
	})

	var c, max int
	for i := 0; i < len(ss); i++ {
		max = ss[i].score
		res = append(res, ss[i])
		if i+1 < len(ss) && max == ss[i+1].score {
			continue
		}
		c++
		if c == 3 {
			break
		}
	}
	return
}

func addUser(name string, score int) error {
	m, er := GetFileData(ResultsPath)
	if er != nil {
		return fmt.Errorf("error in addUser() while performing GetFileData():\n%v", er)
	}
	m[name] = score
	data, er := json.Marshal(m)
	if er != nil {
		return fmt.Errorf("error in addUser() while marshalling the file:\n%v", er)
	}
	er = in.Write(ResultsPath, data)
	if er != nil {
		return fmt.Errorf("error in addUser() while writing into the file:\n%v", er)
	}
	return nil
}
