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
		if err:=res.Scan(&info.Year,&info.Term,&info.RedBalls[0],&info.RedBalls[1],&info.RedBalls[2],&info.RedBalls[3],&info.RedBalls[4],&info.RedBalls[5],&info.BlueBall,&id);err==nil{
			return info,nil
		}else{
			fmt.Println("Scan err",err)
		}
	}
	return nil,nil
}

func EnumAll(proc func (info* Info)()){
	db:=GetDB()
	if res,err:=db.Query("select * from records");err!=nil{
		fmt.Println("slect in db error")
		return
	}else{
		info:=new(Info)
		var id int
		for res.Next(){
			if err:=res.Scan(&info.Year,&info.Term,&info.RedBalls[0],&info.RedBalls[1],&info.RedBalls[2],&info.RedBalls[3],&info.RedBalls[4],&info.RedBalls[5],&info.BlueBall,&id);err==nil{
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
	query:=fmt.Sprintf("insert into records (year,term,rb1,rb2,rb3,rb4,rb5,rb6,bb) values ('%d','%d','%d','%d','%d','%d','%d','%d','%d')",info.Year,info.Term,info.RedBalls[0],info.RedBalls[1], info.RedBalls[2], info.RedBalls[3], info.RedBalls[4], info.RedBalls[5], info.BlueBall) 
	if _,err:=db.Exec(query);err!=nil{
		fmt.Println("Insert into db error:",err)
		return err
	}
	return nil
}

func DelApp(id int64) error{
//	query:=fmt.Sprintf("delete from apps where id='%d'",id)
//	_,err:=db.Query(query)
//	return err
	return nil;
}
/*
func FindApp(id int64) (* AppInfo,error){
	db:=GetDB()
	query:=fmt.Sprintf("select * from apps where id='%d'",id)
	res,err:=db.Query(query)
	if err!=nil{
		log.Println("find store query error:",err)
		return nil,err
	}
	if res.Next(){
		info:=new(AppInfo)
		if err:=res.Scan(&info.ID,	&info.Name,
				&info.Version, &info.Vender, &info.Url, &info.Descr,
				&info.Icon,&info.Cost,&info.Sell, &info.Online);err!=nil{
			log.Println("Query error:",err)
			return nil,err
		}
		return info,nil
	}
	return nil,nil
}

func (info* TrackInfo)RegisterVisit() error{
	db:=GetDB()
    if _,err:=db.Exec(query);err!=nil{
        log.Println("Insert db error:",err)
        return err
    }
	return nil
}

func ViewTracks() ([]string,error){
	db:=GetDB()
	var vtime, name, app string
	ret:=make ([]string, 0, 100)
	query:="select tracks.visit,stores.name,apps.name from tracks,stores,apps where tracks.storeid=stores.id and tracks.appid=apps.id  order by tracks.visit desc";
	res,err:=db.Query(query)
	if err!=nil{
		log.Println("Query quick view of visit tracks error",err)
		return nil,err
	}
	for res.Next(){
		if err:=res.Scan(&vtime,&name,&app);err!=nil{
			log.Println ("Get object from result error:",err)
			return nil,err
		}else{
			ret=append(ret,fmt.Sprintf("%s   %s  %s",vtime,name,app))
		}
	}
	return ret,nil
}

func SearchMatch(from,to,store,app,desc,combine string)([]string,error){
	db:=GetDB()
	var rvtime, rname, rapp string
	ret:=make ([]string, 0, 100)
	prequery:="select tracks.visit,stores.name,apps.name from tracks,stores,apps where tracks.storeid=stores.id and tracks.appid=apps.id  %s %s %s %s order by tracks.visit %s";
	qfrom:=""
	qto:=""
	qstore:=""
	qapp:=""
	if from!=""{
		qfrom="and tracks.visit>="+from
	}
	if to!=""{
		qto=fmt.Sprintf("and tracks.visit<DATE_ADD(\"%s\",INTERVAL 1 DAY)",to)
	}
	if store!=""{
		qstore=fmt.Sprintf("and stores.name like '%%%s%%'",store)
	}
	if app!=""{
		qapp=fmt.Sprintf("and apps.name like '%%%s%%'",app)
	}
	query:=fmt.Sprintf(prequery,qfrom,qto,qstore,qapp,desc)
	res,err:=db.Query(query)
	if err!=nil{
		log.Println("Query quick view of visit tracks error",err)
		return nil,err
	}
	ltime,_:=time.Parse("2006-01-02 15:04:05","1970-01-01 01:00:00")
	lname:=""
	lapp:=""
	for res.Next(){
		if err:=res.Scan(&rvtime,&rname,&rapp);err!=nil{
			log.Println ("Get object from result error:",err)
			return nil,err
		}else{
			record:=true
			if combine=="combined"{
				curtime,_:=time.Parse("2006-01-02 15:04:05",rvtime)
				if lname==rname && lapp==rapp && math.Abs(curtime.Sub(ltime).Seconds())<30{
					record=false
				}
				ltime=curtime
				lname=rname
				lapp=rapp
			}
			if record{
				ret=append(ret,fmt.Sprintf("%s   %s  %s",rvtime,rname,rapp))
			}
		}
	}
	return ret,nil
}

func GetAllApps(storeid int64)([]*AppInfo,error){
	db:=GetDB()
	if storeid<1000{
		return nil,errors.New("Invalid storeid")
	}
	ret:=make([]*AppInfo,0,50)
	query:="select * from apps order by online desc, id desc"
	res,err:=db.Query(query)
	if err!=nil{
		log.Println("Query all apps error:",err)
		return nil,err
	}
	for res.Next(){
		info:=new(AppInfo)
		if err:=res.Scan(&info.ID,	&info.Name,
				&info.Version, &info.Vender,&info.Url,
				 &info.Descr,&info.Icon,
				&info.Cost,&info.Sell, &info.Online);err!=nil{
			log.Println("Get object from db result  error:",err)
			return nil,err
		}else{
			ret=append(ret,info)
		}
	}
	return ret,nil
}
*/
