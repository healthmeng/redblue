package main

import (
"fmt"
"bufio"
"os"
"sort"
)

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

func main(){
	bks:=backet{[33]int{0}, [16]int{0}}
	fd,err:=os.Open("/tmp/rbresult.txt")
	if err!=nil{
		fmt.Println("Open file error:",err)
		return
	}
	defer fd.Close()
	redb:=[6]int{0}
	blueb:=0
	rb:=bufio.NewReader(fd)
	for{
		line,_,err:=rb.ReadLine()
		if err==nil{
			n,err:=fmt.Sscanf(string(line),"%d%d%d%d%d%d%d",&redb[0],&redb[1],&redb[2],&redb[3],&redb[4],&redb[5],&blueb)
			if err==nil && n==7{
				fmt.Println(redb,blueb)
				for j:=0;j<6;j++{
					bks.redbk[redb[j]-1]++
				}
				bks.bluebk[blueb-1]++
			}else{
				fmt.Println("parse data error:",string(line))
			}
		}else{
			break
		}
	}
	var st SortType
	st.data=make([]numdata,33,33)
	for i:=0;i<33;i++{
		st.data[i].index=i+1
		st.data[i].times=bks.redbk[i]
	}
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
