package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var dir string
var base = make(map[string]string)
var currentfile = "./urler"

func ConnectFile(a string) error {
	dir = a
	f, err := os.Open(currentfile)
	if err != nil {
		log.Printf("ERR: FILE %s, Cant open file %s", dir, currentfile)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		log.Println(s)
		base[s[0]] = s[1]
	}
	if len(base) == 0 {

		err = CreateNewFile("./urler")
		if err != nil {
			return err
		}
		currentfile = "./urler"
	}
	return nil
}

func PostLocal(newrow *Baserow) error {
	for k, v := range base {
		if v == newrow.Lurl {
			newrow.Surl = k
			return nil
		}
	}

	err := PostStringFile(newrow)
	if err != nil {
		newrow.Surl = ""
		return err
	}
	base[newrow.Surl] = newrow.Lurl
	return nil
}

func GetLocal(newrow *Baserow) error {
	for k, v := range base {
		if k == newrow.Surl {
			newrow.Lurl = v
			return nil
		}
	}
	newrow.Lurl = ""
	return nil
}

func PostStringFile(newrow *Baserow) error {
	f, err := os.OpenFile(currentfile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Printf("ERR: Cant open file ", currentfile)
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
	stri := fmt.Sprintf("%s %s\n", newrow.Surl, newrow.Lurl)
	_, err = f.WriteString(stri)
	if err != nil {
		log.Printf("ERR: cant write new row %s, %s", newrow, err)
		return err
	}
	return nil
}

func CreateNewFile(currentfile string) error {
	f, err := os.OpenFile(currentfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERR: Cant create file ", currentfile)
		return err
	}
	defer f.Close()
	// _, err = f.WriteString("SURL LURL\n")
	// if err != nil {
	// 	return err
	// }
	log.Printf("file created %s", currentfile)
	return nil
}
