package dbop

import (
"database/sql"
"os"
"fmt"
_"MySQL"
"time"
)

var curdb *sql.DB

type Info struct{
    RedBalls [6] int
    BlueBall int
	Year int
	Term int
	Date string
}

func init(){
	ConnDB()
}

func ConnDB(){
	var err error
	curdb,err=sql.Open("mysql","work:Work4All;@tcp(123.206.55.31:3306)/rbdata")
	if err!=nil{
		fmt.Println("Open database error:",err)
		os.Exit(1)
	}
	curdb.SetConnMaxLifetime(time.Second*500)
}

func GetDB() *sql.DB{
	if err:=curdb.Ping();err!=nil{
		curdb.Close()
		ConnDB()
	}
	return curdb
}

func GetLastRecord(thisyear int)(int,int,error){
	db:=GetDB()
	year:=2003
	term:=0
	for y:=thisyear;y>=2003;y--{
		query:=fmt.Sprintf("select year,term from records where year='%d' order by term desc",y)
		res,err:=db.Query(query)
		if err!=nil{
			return year,term,err
		}
		if res.Next(){
			if err:=res.Scan(&year,&term);err!=nil{
				fmt.Println("Get last record scan error ", err)
				return year,term,err
			}else{
				break
			}
		}
	}
	return year,term,nil
}

func Lookup(year, term int) (*Info,error){
	db:=GetDB()
	var id int
	query:=fmt.Sprintf("select * from records where year='%d' and term='%d'",year, term)
	res,err:=db.Query(query)
	if err!=nil{
		fmt.Println("Lookup in database error:",err)
		return nil,err
	}
	if res.Next(){
		info:=new(Info)
		if err:=res.Scan(&info.Year,&info.Term,&info.RedBalls[0],&info.RedBalls[1],&info.RedBalls[2],&info.RedBalls[3],&info.RedBalls[4],&info.RedBalls[5],&info.BlueBall,&info.Date,&id);err==nil{
			return info,nil
		}else{
			fmt.Println("Scan err",err)
		}
	}
	return nil,nil
}

func EnumAll(startyear int, proc func (info* Info)()){
	db:=GetDB()
	query:=fmt.Sprintf("select * from records where year>=%d",startyear)
	if res,err:=db.Query(query);err!=nil{
		fmt.Println("slect in db error")
		return
	}else{
		info:=new(Info)
		var id int
		for res.Next(){
			if err:=res.Scan(&info.Year,&info.Term,&info.RedBalls[0],&info.RedBalls[1],&info.RedBalls[2],&info.RedBalls[3],&info.RedBalls[4],&info.RedBalls[5],&info.BlueBall,&info.Date,&id);err==nil{
			proc(info)
			}else{
				fmt.Println("Scan query data error",err)
				break
			}
		}
	}
}

func (info* Info)AddInfo() error{
	db:=GetDB()
	query:=fmt.Sprintf("insert into records (year,term,rb1,rb2,rb3,rb4,rb5,rb6,bb,runtime) values ('%d','%d','%d','%d','%d','%d','%d','%d','%d','%s')",info.Year,info.Term,info.RedBalls[0],info.RedBalls[1], info.RedBalls[2], info.RedBalls[3], info.RedBalls[4], info.RedBalls[5], info.BlueBall,info.Date) 
	if _,err:=db.Exec(query);err!=nil{
		fmt.Println("Insert into db error:",err)
		return err
	}
	return nil
}

