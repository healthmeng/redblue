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
	fmt.Println("-s/show show every ball hit info\n")
	fmt.Println("-g/get [param] get suggestions by different way\n")
}
/*
type Info struct{
    RedBalls [6] int
    BlueBall int
	Year int
	Term int
}*/

func ParseLine(buf string,info* dbop.Info,index *int) int{
    patred,_:=regexp.Compile("\\s*<span\\s+class=\\\"ball-list\\s+red\\\">(\\d+)</span>")
    patblue,_:=regexp.Compile("\\s*<span\\s+class=\\\"ball-list\\s+blue\\\">(\\d+)</span>")
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
    info:=dbop.Info{[6]int{0},0,year,term}
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
	lyear,lterm,err:=dbop.GetLastRecord(thisyear)
	if err!=nil{
		fmt.Println("Get last record error",err)
		return
	}
	fmt.Printf("Get last record: %d-%d\n", lyear,lterm)
	for y:=lyear;y<=thisyear;y++{
		for t:=lterm+1;t<160;t++{
			str:=fmt.Sprintf("https://kjh.55128.cn/ssq-kjjg-%d%03d.htm",y,t)
			info:=CheckoutUrl(str,y,t)
			if info!=nil{
				if i,_:=dbop.Lookup(info.Year,info.Term);i==nil{
					info.AddInfo()
					fmt.Printf("%d-%d updated\n",y,t)
				}
			}else{
				break
			}
		}
		lterm=0
	}
}

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
				fmt.Printf("Red: %d, %d, %d, %d, %d, %d; Blue :%d\n",info.RedBalls[0],info.RedBalls[1],
				info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall)
			}else{
				fmt.Println("no data")
			}
		}
	}
	fmt.Println("Update finished!")
}

func SimplePrint(info *dbop.Info){
	fmt.Printf("%d-%d: Red %d, %d, %d, %d, %d, %d; Blue %d\n",info.Year, info.Term,info.RedBalls[0],info.RedBalls[1],info.RedBalls[2],info.RedBalls[3],info.RedBalls[4],info.RedBalls[5],info.BlueBall)
}

func doShowAll(){
	dbop.EnumAll(2003,SimplePrint)
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
	dbop.EnumAll(from,func (info* dbop.Info){
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
		case "-s":
			fallthrough
		case "show":
			doShowAll()
		case "-g":
			fallthrough
		case "get":
				getSuggest()
		}
	}
}
