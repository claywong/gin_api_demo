package models

import (
	"github.com/claywong/gin_api_demo/dbs"
	"github.com/golang/go/src/pkg/fmt"
	"github.com/lunny/log"
	"strconv"
)

type Member struct {
	Id        int    `json:"id" form:"id"`
	LoginName string `json:"login_name" form:"login_name"`
	Password  string `json:"password" form"password"`
}

func (m *Member) AddMember() (id int64, err error) {
	res, err := dbs.Conns.Exec("INSERT INTO ppgo_member(login_name, password) VALUES (?, ?)", m.LoginName, m.Password)
	if err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

func ListMember(page, pageSize int, filters ...interface{}) (lists []Member, count int64, err error) {
	lists = make([]Member, 0)
	where := "WHERE 1=1"
	if len(filters) > 0 {
		lengthOfFilters := len(filters)
		for k := 0; k < lengthOfFilters; k += 3 {
			where = where + " AND " + filters[k].(string) + filters[k+1].(string) + filters[k+2].(string)
		}
	}
	limit := strconv.Itoa((page-1)*pageSize) + "," + strconv.Itoa(pageSize)
	sql := "SELECT id, login_name, password FROM ppgo_member " + where + " LIMIT " + limit
	fmt.Println(sql)
	rows, err := dbs.Conns.Query(sql)
	defer rows.Close()

	if err != nil {
		return
	}
	count = 0
	for rows.Next() {
		var member Member
		rows.Scan(&member.Id, &member.LoginName, &member.Password)
		lists = append(lists, member)
		count++
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func OneMember(id int) (m Member, err error) {
	m.Id = 0
	m.LoginName = ""
	m.Password = ""
	err = dbs.Conns.QueryRow("SELECT id, login_name, password FROM ppgo_member WHERE id=? LIMIT 1", id).Scan(&m.Id, &m.LoginName, &m.Password)
	return
}

func (m *Member) UpdateMember(id int) (n int64, err error) {
	res, err := dbs.Conns.Prepare("UPDATE ppgo_member SET login_name=?,password=? WHERE id=?")
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	rs, err := res.Exec(m.LoginName, m.Password, m.Id)
	if err != nil {
		log.Fatal(err)
	}
	n, err = rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func DeleteMember(id int) (n int64, err error) {
	n = 0
	rs, err := dbs.Conns.Exec("DELETE FROM ppgo_member WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
		return
	}
	n, err = rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}
