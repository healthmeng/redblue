package dbop

import (
"database/sql"
"os"
"fmt"
"errors"
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

type MySelInfo struct{
	RedBalls[6] int
	BlueBall int
	Date string
	Id int
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

func GetLastRecord(thisyear int)(int,int,string,error){
	db:=GetDB()
	year:=2003
	term:=0
	date:=""
	for y:=thisyear;y>=2003;y--{
		query:=fmt.Sprintf("select year,term,runtime from records where year='%d' order by term desc limit 1",y)
		res,err:=db.Query(query)
		if err!=nil{
			return year,term,date,err
		}
		if res.Next(){
			if err:=res.Scan(&year,&term,&date);err!=nil{
				fmt.Println("Get last record scan error ", err)
				return year,term,date,err
			}else{
				break
			}
		}
	}
	return year,term,date,nil
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

func DelSel(id int) bool{
	db:=GetDB()
	query:=fmt.Sprintf("delete from mysel where id=%d",id)
	if res,err:=db.Exec(query);err!=nil{
		fmt.Println("db exec error:",err.Error())
	}else{
		if row,_:=res.RowsAffected();row>0{
			return true
		}
	}
	return false
}

func findLeastID() (int,error){
	db:=GetDB()
	for i:=1;;i++{
		query:=fmt.Sprintf("select bb from mysel where id=%d",i)
		if res,err:=db.Query(query);err!=nil{
			return 0,err
		}else{
			if res.Next(){
				continue
			}else{
				return i,nil
			}
		}
	}
	return 0,errors.New("Impossible error")
}

func (info* MySelInfo)UpdateInfo()(bool,error){
	db:=GetDB()
	ret:=false
	query:=fmt.Sprintf("update mysel set rb1=%d,rb2=%d,rb3=%d,rb4=%d,rb5=%d,rb6=%d,bb=%d,date='%s' where id=%d",info.RedBalls[0],info.RedBalls[1],info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall,info.Date,info.Id) 
	
	if res,err:=db.Exec(query);err!=nil{
		return false,err
	}else if rows,_:=res.RowsAffected();rows>0{
		ret=true
	}
	return ret,nil
}

func InsertSel(info* MySelInfo){
	db:=GetDB()
	var err error
	if info.Id,err=findLeastID();err!=nil{
		fmt.Println("Find id error:",err.Error())
		return
	}
	query:=fmt.Sprintf("insert into mysel(id,rb1,rb2,rb3,rb4,rb5,rb6,bb,date) values ('%d','%d','%d','%d','%d','%d','%d','%d','%s')",info.Id,info.RedBalls[0],info.RedBalls[1],info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall,info.Date)
	if _,err:=db.Exec(query);err!=nil{
		fmt.Println("insert failed:",err.Error())
	}
}

func GetSelected() []*MySelInfo{
	db:=GetDB()
	query:=fmt.Sprintf("Select * from mysel order by id asc")
	if res,err:=db.Query(query);err!=nil{
		fmt.Println("select in db error",err.Error())
		return nil
	}else{
	list:=make([]*MySelInfo,0,100) // should <10
	for res.Next(){
		info:=new(MySelInfo)
		if err:=res.Scan(&info.Id,&info.RedBalls[0],&info.RedBalls[1],&info.RedBalls[2],&info.RedBalls[3],&info.RedBalls[4],&info.RedBalls[5],&info.BlueBall,&info.Date);err!=nil{
			fmt.Println("Scan query error in GetSeled:",err.Error())
			return nil
		}else{
			list=append(list,info)
		}
		}	
		return list
	}
}

func EnumAll(startyear int, limit int64,proc func (info* Info)()){
	db:=GetDB()
	query:=fmt.Sprintf("select count(*) as value from records")
	var rows int64
    	if err := db.QueryRow(query).Scan(&rows); err != nil {
        	fmt.Println("Query rows error")
        	return 
    	}
	if rows<limit || limit <0{
		limit=rows
	}
	
	query=fmt.Sprintf("select * from records where year>=%d limit %d,%d",startyear,rows-limit,limit)
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

