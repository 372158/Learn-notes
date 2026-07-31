package main

//:name —— 匹配单个路径段。例如，/user/:name 匹配 /user/john，但不匹配 /user/ 或 /user。
//*action —— 匹配前缀之后的所有内容，包括斜杠。例如，/user/:name/*action 匹配 /user/john/send 和 /user/john/。捕获的值包含前导 /。

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//这段处理程序将匹配“/user/jion",但不会匹配“/user/” 或“/user”
	router.GET("/user/:name", func (c *gin.Context)  {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
		
	})

// 但是，这个正则表达式既能匹配“/user/john/”，也能匹配“/user/john/send” 
// 如果没有其他路由器与“/user/john”相匹配，它将重定向至“/user/john/”路径。
	router.GET("user/:name/*action", func (c *gin.Context)  {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is" + action
		c.String(http.StatusOK, message)
		
	})

	router.Run(":8080")
}
