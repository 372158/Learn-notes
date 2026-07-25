package main

import (
	"fmt"
)

//编写一个学生考试系统

type Student struct {
	Name string
	Age int
	Score float64
}

//将Pupil 和 Graduate共有方法也绑定到 *Student
func (s *Student) ShowInfo() {
	fmt.Printf("学生姓名： %v,年龄： %v,成绩： %v\n", s.Name, s.Age, s.Score)

}

func (s *Student) SetScore(score float64) {
	//业务判断
	s.Score = score
}

//小学生
type Pupil struct {
	Student //嵌入了 Student匿名结构体
}

//显示他的成绩


//Pupil结构体特有的方法，保留
func (p *Pupil) testing() {
	fmt.Println("小学生考试中……")
}

//大学，研究生



//大学生
type Graduate struct {
	Student// 嵌入了Student 匿名结构体
}

//显示他的成绩

//这是Graduate结构体特有的方法，保留

func (p *Graduate) testing() {
	fmt.Println("大学生正在考试……")
}


func main() {
	//当我们对结构体潜入了匿名结构体使用方法会发生变化
	pupil := &Pupil{}
	pupil.Student.Name = "Tom"
	pupil.Student.Age = 8
	pupil.testing()
	pupil.Student.SetScore(70)
	pupil.Student.ShowInfo()

	graduate :=&Graduate{}
	graduate.Student.Name = "Mary"
	graduate.Student.Age = 18
	graduate.testing()
	graduate.Student.SetScore(90.5)
	graduate.Student.ShowInfo()


}