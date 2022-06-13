// Database connection

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
	//timer()
	autoCalculateSquad()
}

func timer() {
	msc, _ := time.LoadLocation("Europe/Moscow")
	c := cron.New(cron.WithLocation(msc))

	c.AddFunc("@every 1m", func() {
		deleteFromParentsTimer()
	})

	c.Start()

	for {
		time.Sleep(time.Second * 1)
	}
}

func deleteFromParentsTimer() {
	connStr := "user=postgres password=6858 dbname=test_1 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type dateReg struct {
		id        int
		dateR     time.Time
		statusPay bool
	}

	rows, err := db.Query("select id, date_reg, status_pay from parents")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dateRegs := []dateReg{}

	for rows.Next() {
		d := dateReg{}
		err := rows.Scan(&d.id, &d.dateR, &d.statusPay)
		if err != nil {
			log.Fatal(err)
		}
		dateRegs = append(dateRegs, d)
	}

	for _, d := range dateRegs {
		if !d.statusPay {
			differDate := time.Since(d.dateR)
			fmt.Println(differDate)

			if differDate.Hours() > 720 {
				result, err := db.Exec("delete from parents where id = $1", d.id)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(result.RowsAffected())
			} else {
				fmt.Println("payment deadline has not expired")
			}
		}
	}
}

func autoCalculateSquad() {
	connStr := "user=postgres password=6858 dbname=test_1 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type childCount struct {
		count int
	}

	row := db.QueryRow("select * from child_count")
	childC := childCount{}
	err = row.Scan(&childC.count)
	if err != nil {
		panic(err)
	}

	totalPlace := childC.count
	squadPlace := 25
	i := 1

	if totalPlace%squadPlace != 0 {
		for ; i <= totalPlace/squadPlace; i++ {
			_, err := db.Exec("insert into groups (group_year, group_number, count) values (0, $1, $2)", i, squadPlace)
			if err != nil {
				panic(err)
			}
		}

		_, err := db.Exec("insert into groups (group_year, group_number, count) values (0, $1, $2)", i, totalPlace%squadPlace)
		if err != nil {
			panic(err)
		}
		i = 0
	} else {
		for ; i <= totalPlace/squadPlace; i++ {
			_, err := db.Exec("insert into groups (group_year, group_number, count) values (0, $1, $2)", i, squadPlace)
			if err != nil {
				panic(err)
			}
		}
		i = 0
	}
}
