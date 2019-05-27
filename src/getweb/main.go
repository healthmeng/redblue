package main

import (
"fmt"
"net/http"
"bufio"
"regexp"
"strconv"
)

type Info struct{
	RedBalls [6] int
	BlueBall int
}

func ParseLine(buf string,info* Info,index *int) int{
/*

                <span class="ball-list red">26</span>
                <span class="ball-list blue">11</span>
*/
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

func CheckoutUrl(url string) *Info{
	res,err:=http.Get(url);
//	res,err:=http.Get("https://kjh.55128.cn/ssq-kjjg-2016057.htm");
	if err!=nil || res.StatusCode!=200{
//		fmt.Println(err)
		return nil
	}
//	fmt.Println(res.StatusCode)
	defer res.Body.Close()
	index:=0
	info:=Info{[6]int{0},0}
	rb:=bufio.NewReader(res.Body)
	for{
		line,_,err:=rb.ReadLine()
		if err!=nil || ParseLine(string(line),&info,&index)==0{ // when get blue ,return 0
			break
		}
	}
	return &info
}

func main(){
	for y:=2003;y<=2019;y++{
		for q:=1;q<160;q++{
			str:=fmt.Sprintf("https://kjh.55128.cn/ssq-kjjg-%d%03d.htm",y,q)
			if info:=CheckoutUrl(str);info!=nil{
					fmt.Println(y,q,":",info.RedBalls,"--",info.BlueBall)
		}
	}
	}
}
