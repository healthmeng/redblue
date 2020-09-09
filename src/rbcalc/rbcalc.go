package main

import (
"fmt"
"time"
"os"
"dbop"
"regexp"
"bufio"
"net/http"
"strconv"
"sort"
)


func printUsage(){
	fmt.Println("Usage:\n-u/update update latest data into database\n")
	fmt.Println("-a/all show every ball hit info\n")
	fmt.Println("-g/get [param] get suggestions by different way\n")
	fmt.Println("-s/set set current selected balls")
}
/*
type Info struct{
    RedBalls [6] int
    BlueBall int
	Year int
	Term int
	Date string
}*/

func ParseLine(buf string,info* dbop.Info,index *int) int{
    patred,_:=regexp.Compile("\\s*<span\\s+class=\\\"ball-list\\s+red\\\">(\\d+)</span>")
    patblue,_:=regexp.Compile("\\s*<span\\s+class=\\\"ball-list\\s+blue\\\">(\\d+)</span>")
    patdate,_:=regexp.Compile("\\s*开奖时间.+(\\d{4}-\\d+-\\d+)</span>")
    if date:=patdate.FindStringSubmatch(string(buf));date!=nil{
	info.Date=date[1]
	return -1
	}
    if num:=patred.FindStringSubmatch(string(buf));num!=nil{
        info.RedBalls[*index],_=strconv.Atoi(num[1])
        (*index)++
        return (*index)
    }else if num:=patblue.FindStringSubmatch(string(buf));num!=nil{
        info.BlueBall,_=strconv.Atoi(num[1])
        return 0
    }
    return -1
}

func CheckoutUrl(url string,year, term int) *dbop.Info{
    res,err:=http.Get(url);
    if err!=nil || res.StatusCode!=200{
        return nil
    }
    defer res.Body.Close()
    index:=0
    info:=dbop.Info{[6]int{0},0,year,term,"19000101"}
    rb:=bufio.NewReader(res.Body)
    for{
        line,_,err:=rb.ReadLine()
        if err!=nil || ParseLine(string(line),&info,&index)==0{ // when get blue ,return 0
            break
        }
    }
	if index<6{
		return nil
	}else{
		return  &info
	}
}

func doUpdateOpt(){
	thisyear:=time.Now().Year()
	lyear,lterm,date,err:=dbop.GetLastRecord(thisyear)
	if err!=nil{
		fmt.Println("Get last record error",err)
		return
	}
	fmt.Printf("Get last record: %d-%d(%s)\n", lyear,lterm,date)
	for y:=lyear;y<=thisyear;y++{
		for t:=lterm+1;t<160;t++{
			str:=fmt.Sprintf("https://kjh.55128.cn/ssq-kjjg-%d%03d.htm",y,t)
			info:=CheckoutUrl(str,y,t)
			if info!=nil{
				if i,_:=dbop.Lookup(info.Year,info.Term);i==nil{
					info.AddInfo()

					fmt.Printf("%d-%d(%s) updated: Red: %d, %d, %d, %d, %d, %d; Blue :%d\n",y,t,info.Date,info.RedBalls[0],info.RedBalls[1],info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall)

				checkMatch(info)
				}
			}else{
				break
			}
		}
		lterm=0
	}
}


func checkMatch(info *dbop.Info) int{
	money:=0
	dst:=make(map[int] int)
	for _,v:=range info.RedBalls{
		dst[v]=1
	}
	s:=dbop.GetSelected()	
	for _,mysel:=range s{
		redhit:=0
		for _,v:=range mysel.RedBalls{
			if dst[v]==1{
				redhit++
			}
		}	
	bluehit:=false
	money=0
	if mysel.BlueBall==info.BlueBall{
		bluehit=true		
	}
	if bluehit{
		switch redhit{
		case 0,1,2:
			money=5
		case 3:
			money=10
		case 4:
			money=200
		case 5:
			money=3000
		case 6:
			money=5000000
		}
	}else{
		switch redhit{
		case 4:
			money=10
		case 5:
			money=200
		case 6:
			money=50000
		}
	}
	if money>0{
		fmt.Println("Hit!!! ", mysel.RedBalls,"[",mysel.BlueBall,"]  Get bonus:￥",money)
	}
	}
	return money
}

/* 
func doUpdate(){
	year:=time.Now().Year()
StopFind:
	for y:=year;y>=2003;y--{
		for term:=160;term>0;term--{
		if f,err:=dbop.Lookup(y,term);err!=nil || f!=nil{
				break StopFind
			}
			str:=fmt.Sprintf("https://kjh.55128.cn/ssq-kjjg-%d%03d.htm",y,term)
			fmt.Printf("Searching %d-%d ...",y,term)
			info:=CheckoutUrl(str,y,term)
			if info!=nil{
				info.AddInfo()
				fmt.Printf("Red: %d, %d, %d, %d, %d, %d; Blue :%d.\n",info.RedBalls[0],info.RedBalls[1],
				info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall)
				checkMatch(info)
			}else{
				fmt.Println("no data")
			}
		}
	}
	fmt.Println("Update finished!")
}*/

func SimplePrint(info *dbop.Info){
	fmt.Printf("%d-%d(%s): Red %d, %d, %d, %d, %d, %d; Blue %d\n",info.Year, info.Term,info.Date,info.RedBalls[0],info.RedBalls[1],info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall)
}

func CheckBonus(info *dbop.Info){
	fmt.Printf("%d-%d(%s): Red %d, %d, %d, %d, %d, %d; Blue %d.   ",info.Year, info.Term,info.Date,info.RedBalls[0],info.RedBalls[1],info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall)
	if checkMatch(info)==0{
		fmt.Println("")
	}
}
func doShowAll(limit int64){
	if limit<0{
		dbop.EnumAll(2003,limit,SimplePrint)
	}else{
		dbop.EnumAll(2003,limit,CheckBonus)
	}
}

type backet struct{
    redbk [33]int
    bluebk [16]int
}
type numdata struct{
    index int
    times int
}

type SortType struct{
    data []numdata
}

func (s *SortType)Len() int{
    return len(s.data)
}

func (s *SortType)Less(i,j int) bool{
    return s.data[i].times<s.data[j].times
}

func (s *SortType)Swap(i,j int){
    tmp:=s.data[i]
    s.data[i]=s.data[j]
    s.data[j]=tmp
}


func GetLeast(from int){
	bks:=backet{[33]int{0},[16]int{0}}
	dbop.EnumAll(from,-1,func (info* dbop.Info){
		for i:=0;i<6;i++{
			bks.redbk[info.RedBalls[i]-1]++
		}
		bks.bluebk[info.BlueBall-1]++
	})
    var st SortType
    st.data=make([]numdata,33,33)
    for i:=0;i<33;i++{
        st.data[i].index=i+1
        st.data[i].times=bks.redbk[i] }
    sort.Sort(&st)
    fmt.Println("Red info:")
    for i:=0;i<33;i++{
            fmt.Println(st.data[i].index,"--",st.data[i].times)
    }
    st.data=make([]numdata,16,16)
    for i:=0;i<16;i++{
        st.data[i].index=i+1
        st.data[i].times=bks.bluebk[i]
    }
    sort.Sort(&st)
    fmt.Println("Blue info:")
    for i:=0;i<16;i++{
            fmt.Println(st.data[i].index,"--",st.data[i].times)
    }

}

func getSuggest(){
// switch predictway{
// case SimpleStatis:
	year:=2003
	if len(os.Args)>=3{
		y,err:=strconv.Atoi(os.Args[2])
		if err==nil{
			year=y
		}
	}
	GetLeast(year)
}

func prtCurSel(){
	infos:=dbop.GetSelected()
	for _,i:=range infos{
		fmt.Println(i.Id,": Red:",i.RedBalls,"Blue:[",i.BlueBall,"]; Selected in:",i.Date)
	}
}

func addSel(){
	info:=new (dbop.MySelInfo)
	fmt.Println("Six red balls:")
	for i:=0;i<6;i++{
		fmt.Scanf("%d",&info.RedBalls[i])
	}
	fmt.Println("Blue ball:")
	fmt.Scanf("%d",&info.BlueBall)
	curtm:=time.Now()
	info.Date=fmt.Sprintf("%04d-%02d-%02d",curtm.Year(),curtm.Month(),curtm.Day())
	dbop.InsertSel(info)
}

func delSel(){
	var sel int
	fmt.Print("Input ID to be deleted:")
	fmt.Scanf("%d",&sel)
	if !dbop.DelSel(sel){
		fmt.Println("Can't find id ",sel)
	}else{
		fmt.Println("Delelted ok")
	}
}

func editSel(){
	var sel int
	fmt.Print("Input ID to be edited:")
	fmt.Scanln(&sel)
	info:=new (dbop.MySelInfo)
	fmt.Print("Six red balls:")
	for i:=0;i<6;i++{
		fmt.Scanln(&info.RedBalls[i])
	}
	fmt.Print("Blue ball:")
	fmt.Scanln(&info.BlueBall)
	curtm:=time.Now()
	info.Date=fmt.Sprintf("%04d-%02d-%02d",curtm.Year(),curtm.Month(),curtm.Day())
	info.Id=sel
	if ok,err:=info.UpdateInfo();ok==true{
		fmt.Println("Edit ok")
	}else if err!=nil{
		fmt.Println("Edit error:",err.Error())
	}
}

func setSelect(){
oversel:
	for{
		prtCurSel()
		var sel string
		fmt.Println("(A) Add... (D) Del. (E) Edit... (Q) Quit")
		fmt.Scanln(&sel)
		switch sel{
			case "A":
				fallthrough
			case "a":
				addSel()
			case "D":
				fallthrough
			case "d":
				delSel()
			case "E":
				fallthrough
			case "e":
				editSel()
			case "q":
				fallthrough
			case "Q":
				break oversel
		}
		
	}
}

func main(){
	argc:=len(os.Args)
	if argc <2 {
		printUsage()
	}else{
		switch os.Args[1]{
		case "-u":
			fallthrough
		case "update":
			doUpdateOpt()
		case "-a":
			fallthrough
		case "all":
			var recent int64=-1
			if argc>=3{
				if rec,err:=strconv.ParseInt(os.Args[2],10,64);err==nil{
					recent=rec
				}
			}
			doShowAll(recent)
		case "-g":
			fallthrough
		case "get":
			getSuggest()
		case "-s":
			fallthrough
		case "set":
			setSelect()
		}
	}
}
