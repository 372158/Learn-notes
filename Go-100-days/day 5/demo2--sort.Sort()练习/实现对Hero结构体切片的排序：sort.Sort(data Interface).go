package main

import (
	"fmt"
	"math/rand"
	"sort"
)

//1.声明Hero结构体
type Hero struct {
	Name string
	Age int
}

//2.声明一个Hero结构体切片类型
type HeroSlice []Hero

type StudentSlice []Student
//3.实现Interface接口

func (hs HeroSlice) Len() int {
	return len(hs)
}


//Less方法就是决定你使用什么标准进行排序
//1.按照Hero 的Age从小到大排序！！
func (hs HeroSlice) Less(i, j int) bool {
	return hs[i].Age < hs[j].Age

	//修改成对Name进行排序
	//return hs[i].Name < hs[j].Name
}



func (hs HeroSlice) Swap(i, j int) {
	//交换
	//temp := hs[i]
	//hs[i] = hs[j]
	//hs[j] = temp
	hs[i], hs[j] = hs[j], hs[i]
}


//1.声明Student结构体
type Student struct {
	Name string
	Age int
	Score float64
}

//将Student的切片，按Score从大到小排序！！

func (ss StudentSlice) Len() int {
	return len(ss)
}

func (ss StudentSlice) Less(i, j int) bool {
	return ss[i].Score > ss[j].Score // 从大到小排序
}

func (ss StudentSlice) Swap(i, j int) {
	//交换
	//temp := ss[i]
	//ss[i] = ss[j]
	//ss[j] = temp
	ss[i], ss[j] = ss[j], ss[i]
}
func main() {
	//先定义一个数组/切片
	var intSlice = []int{0, -1, 10, 7, 90}
	//要求对 intSlice 切片进行排序
	sort.Ints(intSlice)//系统方法
	fmt.Println(intSlice)

	// //对结构体切片进行排序
	// var heroes HeroSlice
	// for i:= 0; i < 10; i++ {
	// 	hero := Hero{
	// 		Name :fmt.Sprintf("英雄|%d",rand.Intn(100)),
	// 		Age : rand.Intn(100),
	// 	}
	// 	//将hero append 到heroes切片中
	// 	heroes = append(heroes,hero)
	// }
	// //看看排序前的顺序
	// for _, v := range heroes {
	// 	fmt.Println(v)
	// }

	// //调用sort.Sort()方法，传入我们实现了Interface接口的HeroSlice类型变量
	// sort.Sort(heroes)
	// fmt.Println("------排序后------")
	// //看看排序后的顺序
	// for _, v := range heroes {
	// 	fmt.Println(v)
	// }

	var students StudentSlice
	for i := 0; i < 10; i++ {
		student := Student{
			Name: fmt.Sprintf("学生|%d",rand.Intn(100)),
			Age: rand.Intn(100),
			Score: rand.Float64() * 100,
		}
		students = append(students, student)
	}
	//看看排序前的顺序
	for _, v := range students {
		fmt.Println(v)
	}
	//调用sort.Sort()方法，传入我们实现了Interface接口的StudentSlice类型变量
	sort.Sort(students)
	fmt.Println("------排序后------")
	//看看排序后的顺序
	for _, v := range students {
		fmt.Println(v)
	}


	}
	