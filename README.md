# gofly-orm ORM框架

**一套go语言的orm框架**

### TO DO
- [x] 对象表结构映射
- [x] 表的创建、删除、存在性判断
- [x] 记录新增查询
- [ ] 链式操作与更新删除
- [ ] 实现钩子(Hooks)
- [ ] 支持事务(Transaction)
- [ ] 支持模板引擎 html/templates
- [ ] ....


### Use：

```go
package main

import (
	"fmt"
	"github.com/QXQZX/gofly-orm/gform"
	"log"

	// 导入时会注册 sqlite3 的驱动
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `gform:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {

	engine, _ := gform.NewEngine("sqlite3", "gofly.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})

	err1 := s.DropTable()
	err2 := s.CreateTable()

	rows, err3 := s.Insert(user1, user2, user3)

	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("failed init test records")
	}
	fmt.Println("插入成功:", rows)

	var users []User
	if err := s.Find(&users); err != nil {
		log.Fatal("failed to query all")
	}

	fmt.Println(users)
}
```

<hr>
仅学习使用