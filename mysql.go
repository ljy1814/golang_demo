package main

/*
CREATE TABLE `userdetail` (         `uid` INT(10) NOT NULL DEFAULT '0',         `intro` TEXT NULL,         `profile` TEXT NULL,         PRIMARY KEY (`uid`)     )ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `userinfo` (         `uid` INT(10) NOT NULL AUTO_INCREMENT,         `username` VARCHAR(64) NULL DEFAULT NULL,         `departname` VARCHAR(64) NULL DEFAULT NULL,         `created` DATE NULL DEFAULT NULL,         PRIMARY KEY (`uid`)     )ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/beego?charset=utf8")
	fmt.Println(db)
	fmt.Println(err)
	checkErr(err)

	stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	fmt.Println(stmt)
	res, err := stmt.Exec("astaxie", "研发部门", "2017-02-17")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)

	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)

	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string

		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}
	db.Close()

	fmt.Println("OK")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
