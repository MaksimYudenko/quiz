***

### Quiz game

***

This program will progress as follows read a json filepath and a time limit from flags prompt for a key press on key 
press, start the quiz as follows	

While time has not elapsed:
* print a random question to the screen
* prompt the user for an answer
* store the answer in a container
* normalize answers so they compare correctly
* output total questions answered correctly and how many questions there were
	
### Installation
	go get github.com/MaksimYudenko/quiz
	
### Keybindings

###### Available actions
* to start the quiz - choose 1
* to view top three gamers - choose 2
* to leave the quiz - choose 3
 * if you change your mind - press Ctrl+C
 
 
###### Run the quiz via terminal
	go run ./cmd/main.go

* you are able to tweak some parameters: use -h flag to view what and how

###### Run the quiz using docker
	make all
ps: at current stage app does not save your results during playing inside the docker.