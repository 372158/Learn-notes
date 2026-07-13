# DAY 1 收获与感受 

今天把Go的基础语法看过了一遍，掌握的还行，下面说一下个人感受

首先，Go程序的必须品得有：

~~~go
package main

import "fmt" 

func main() {
    
    ```
}
~~~

这个段代码必包含的，就好比C语言里的 头文件+main主函数一样：

~~~c
#include<stdio.h>
int main() {

}
~~~

然后，就是Go的输入和输出，相比较于C语言还是有点像的，就是需要导入包 “fmt”：

~~~go
package main 

import "fmt"
var a,b int
func main() {
    fmt.Scan(&a,&d) //输入
    fmt.Pirntf("%d %d",a,b) //输出
}
~~~

今天的代码练习主要还是，熟悉基本的输入输出和循环，像指针、结构体、map、方法都还没有具体练习，不过和C语言还是有很大程度的相似的。



