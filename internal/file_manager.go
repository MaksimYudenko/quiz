package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Create(path string) error {
	name := filepath.Base(path)

	if _, er := os.Stat(path); er != nil {
		if os.IsNotExist(er) {
			if name == FolderName {
				er = os.MkdirAll(path, 0777)
				if er != nil {
					return fmt.Errorf("error in Create():\n%v", er)
				}
				_, er = fmt.Printf("directory [%s] has been created\n", name)
			} else {
				file, er := os.Create(path)
				if er != nil {
					return fmt.Errorf("error in Create():\n%v", er)
				}
				if name == QuestionsName {
					er = Write(path, []byte(InitialQuestions))
					if er != nil {
						return fmt.Errorf("error in Create() while writing %s:\n%v",
							QuestionsName, er)
					}
				}

				if name == ResultsName {
					er = Write(path, []byte(InitialResult))
				}
				if er != nil {
					return fmt.Errorf("error in Create() while writing %s:\n%v",
						ResultsName, er)
				}
				_, er = fmt.Printf("file [%s] has been created\n", name)
				if er != nil {
					return fmt.Errorf("error in Create():\n%v", er)
				}
				defer file.Close()
			}
		} else {
			return er
		}
	}
	return nil
}

func Write(path string, data []byte) error {
	// open file using READ & WRITE permission
	file, er := os.OpenFile(path, os.O_RDWR, 0666)
	if er != nil {
		return fmt.Errorf("error in Write() while opening:\n%v", er)
	}
	defer file.Close()
	// write some text line-by-line to file
	er = ioutil.WriteFile(path, data, 0666)
	if er != nil {
		return fmt.Errorf("error in Write() while writing:\n%v", er)
	}
	// save changes
	er = file.Sync()
	if er != nil {
		return fmt.Errorf("error in Write() while syncing:\n%v", er)
	}
	return nil
}

func Read(filePath string) ([]byte, error) {
	file, er := os.OpenFile(filePath, os.O_RDWR, 0666)
	if er != nil {
		return nil, fmt.Errorf("error in Read() while opening the file:\n%v", er)
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func Del(path string) error {
	er := os.Remove(path)
	if er != nil {
		return fmt.Errorf("error in del() while removing:\n%v", er)
	}
	_, _ = fmt.Printf("File [%s] has been deleted.", filepath.Base(path))
	return nil
}
