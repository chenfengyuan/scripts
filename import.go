package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

const (
	QunFormat      = `create table "%s" (id int, qun_num bigint, mast_qq bigint, create_date date, title text, class bigint, qun_text text);`
	GroupFormat    = `create table "%s" (id int, qq_num bigint, nick text, age smallint, gender smallint, auth int, qun_num bigint);`
	QunIdxFormat   = `create index on "%s" (qun_num)`
	GroupIdxFormat = `create index on "%s" (qq_num, qun_num)`
)

func ifTableNeedCopy(db *sql.DB, table string) bool {
	c := 0
	err := db.QueryRow(`select count(1) from information_schema.tables where table_name = $1 `, table).Scan(&c)
	if err != nil {
		log.Fatal(err)
	}
	if c == 0 {
		return true
	}
	err = db.QueryRow(fmt.Sprintf(`select count(1) from "%s"`, table)).Scan(&c)
	if c == 0 {
		return true
	}
	return false
}

func processTable(db *sql.DB, fn string, tableFormat string, idxFormat string) {
	log.Println(fn)
	_, err := db.Exec(fmt.Sprintf(`drop table if exists "%s"`, fn))
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(fmt.Sprintf(tableFormat, fn))
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(fmt.Sprintf(idxFormat, fn))
	if err != nil {
		log.Fatal(err)
	}
	fin, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()
	reader := bufio.NewReader(fin)

	fout, err := os.OpenFile("/home/pi/temp.csv", os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = fout.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("generating temp.csv")
	reader.ReadLine()
	reader.ReadLine()
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		if isPrefix {
			log.Fatal("line too long")
		}
		if len(line) == 0 {
			break
		}
		_, err = fout.Write(line)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fout.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	fout.Close()
	log.Println("copying")
	_, err = db.Exec(fmt.Sprintf(`copy "%s" from '/home/pi/temp.csv' (format csv) `, fn))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")
	return
}

func main() {
	db, err := sql.Open("postgres", "user=pi dbname=qq host=/var/run/postgresql")
	fileInfos, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		fn := fileInfo.Name()
		if len(fn) >= 9 {
			if fn[:7] == "QunInfo" {
				if ifTableNeedCopy(db, fn) {
					processTable(db, fn, QunFormat, QunIdxFormat)
				}
				continue
			}
			if fn[:9] == "GroupData" {
				if ifTableNeedCopy(db, fn) {
					processTable(db, fn, GroupFormat, GroupIdxFormat)
				}
				continue
			}
		}
	}
}
