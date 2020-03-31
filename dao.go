package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type FieldInfo struct {
	name     string
	datatype string
}

func dbInsert(m StringList, tb string) {
	fis := getFields(tb)
	s := "insert into " + tb + " ("
	var key, val string
	keys := make([]string, len(fis))
	values := make([]string, len(fis))
	for i, v := range fis {
		key += v.name + ","
		val += "?,"
		keys[i] = v.name
	}
	s = s + key
	s = s[:len(s)-1]
	s = s + ") values ("
	val = val[:len(val)-1] + ")"
	s = s + val
	dsn := c.Dbinfo.User + ":" + c.Dbinfo.Pwd + "@tcp(" + c.Dbinfo.Uri + ":" + c.Dbinfo.Port + ")/tias?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open database err in doSwdjInsert method")
	}
	defer db.Close()
	//stmt,err:=db.Prepare(s)
	//if err != nil {
	//	log.Println(err.Error())
	//	log.Fatalf("prepare failure %s",s)
	//}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("create tx fail...,%v", err)
	} else {
		for _, top := range m {
			it := make([]interface{}, 0)
			for j := 0; j < len(keys); j++ {
				values[j] = top[keys[j]]
				it = append(it, values[j])
			}

			_, err := tx.Exec(s, it...)
			if err != nil {
				log.Fatalf("Exec insert failure: %v", err)
				tx.Rollback()
			}
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("commit err,%v", err)
		}
	}
	//stmt.Close()
}
func getFields(tb string) []FieldInfo {
	djInfo := make([]FieldInfo, 0)
	dsn := c.Dbinfo.User + ":" + c.Dbinfo.Pwd + "@tcp(" + c.Dbinfo.Uri + ":" + c.Dbinfo.Port + ")/tias?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open database err in doSwdjInsert method")
	}
	defer db.Close()
	rows, err := db.Query(`select cs.column_name,cs.data_type from information_schema.columns cs where table_schema ='tias' and table_name = '` + tb + `' and extra = "" and column_default is null;`)
	if err != nil {
		log.Fatalf("query table gt3_dj_nsrxx structure err:%v", err)
	}
	for rows.Next() {
		var v FieldInfo
		err := rows.Scan(&v.name, &v.datatype)
		if err != nil {
			log.Fatalf("Rows Scan %v.", err)
		}
		djInfo = append(djInfo, v)
	}
	return djInfo
}
