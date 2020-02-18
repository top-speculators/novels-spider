package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main_test() {
	// 日志输出到控制台的同时，记录到文件
	// 必须写在 gin.Defualt 之前
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Printf("get form err: %s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 如果只记录到文件
	// gin.DefaultWriter = io.MultiWriter(f)

	// gin.DebugPrintRouteFunc 可以用于定义路由日志的格式

	r := gin.Default() // 这个方法默认使用了 Logger 和 Recovery 中间件，如果想要不使用默认中间件，则使用 gin.New()

	// 静态文件服务，将 url path 映射到服务器的文件夹路径
	r.Static("/assets", "./assets")
	r.StaticFS("/more_static", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 返回 json 数据
	r.GET("/c.JSON", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"lang": "GO语言",
			"tag":  "<br>",
		})
	})

	// 返回完全被转义的 json 数据，包括中文都会被转义
	r.GET("/c.AsciiJSON", func(c *gin.Context) {
		c.AsciiJSON(200, gin.H{
			"lang": "GO语言",
			"tag":  "<br>",
		})
	})

	// HTML 渲染
	//r.LoadHTMLGlob("html/*.tmpl")
	r.GET("/c.HTML", func(c *gin.Context) {

		// 也可以在 main 方法下直接使用 r.LoadHTMLGlob("html/*") 的形式，将所有 html 文件都加载进来
		//r.LoadHTMLFiles("html/index.tmpl")
		c.HTML(200, "index.tmpl", gin.H{
			"title": "hello world",
		})

		//head := r.HTMLRender.Instance("head.tmpl", gin.H{
		//	"title": "hello world",
		//})
		//
		//footer := r.HTMLRender.Instance("footer.tmpl", gin.H{
		//	"d": "66",
		//})
		//
		//c.HTML(200, "index.tmpl", gin.H{
		//	"head":   head,
		//	"footer": footer,
		//
		//	"content": "6666",
		//})

		// 还可以自定义模板渲染器、自定义页面渲染的标签分隔符
	})

	// Multipart/Urlencoded 绑定，根据 Content-Type 将请求内容解析出来，并绑定到某一结构体上
	r.POST("/c.ShouldBind", func(c *gin.Context) {
		type LoginForm struct {
			User     string `form:"user" binding:"required"`
			Password string `form:"password" binding:"required"`
		}
		var form LoginForm
		if c.ShouldBind(&form) == nil { // 绑定成功
			c.JSON(200, gin.H{
				"content": "ok!",
			})
		} else { // 绑定失败
			c.JSON(http.StatusBadRequest, gin.H{
				"content": "error!",
			})
		}
	})

	// Multipart/Urlencoded 表单取值
	r.POST("/c.PostForm", func(c *gin.Context) {
		message := c.PostForm("message")               // 从 post 表单中得到一个 key 为 message 的值，若不存在，则返回一个空字符串
		nick := c.DefaultPostForm("nick", "anonymous") // 可以指定默认值的 PostForm

		c.JSON(200, gin.H{
			"message": message,
			"nick":    nick,
		})
	})

	// PureJSON 返回未被转义的字面字符，如 <b> 中的 <> 符号，用 JSON 方法会被转义成 unicode 实体，而使用 PureJSON 则原样输出
	r.GET("/c.PureJSON", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>",
		})
	})

	// Query 方法获取 url 上的参数
	r.GET("/c.Query", func(c *gin.Context) {
		id := c.Query("id") // 如果不存在，则为空字符串

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})

	r.POST("/c.Query", func(c *gin.Context) {
		id := c.Query("id") // 如果不存在，则为空字符串

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})

	// SecureJSON 方法会为返回的 json 加上前缀，以防止 json 劫持

	// c.XML / c.YAML / c.ProtoBuf 可以分别渲染不同格式的数据，就跟 c.JSON 一样

	// 使用 gin.H 是 map[string]interface{} 的快捷方式，也可以直接使用结构体来返回

	// c.FormFile 可以处理单文件上传，请求必须是 multipart/form-data
	r.POST("/upload", func(c *gin.Context) {
		// 默认允许 multipart forms 的内存限制是 32 mib ；可以使用以下语句来修改该限制
		// r.MaxMultipartMemory = 8 << 20 // 8 mib，如果 8 << 20 << 20 就是 8 gib

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil { // 保存文件
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
	})

	// 多文件上传
	r.POST("/upload2", func(c *gin.Context) {
		form, err := c.MultipartForm() // 获取表单结构体
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		files := form.File["files"] // 获取文件数组（其实是切片）

		for _, file := range files { // 用 _ 丢弃掉了下标，循环保存
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
		}

		c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files with fields.", len(files)))
	})

	// 使用 reader 读取数据
	r.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://cdn4.buysellads.net/uu/1/46074/1559075223-slack-carbon-green_2x.png") // GET 得到一张图片
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable) // 如果 err 有值，或 response.StatusCode 不是 200
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		// 扩展头部
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		// 将指定的 reader 写入主体流并更新 HTTP 代码
		// 只要传递的 reader 实现了 io.Reader 接口
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	// 路由组使用 gin.BasicAuth() 中间件
	// 可以使用 Group 方法来定义路由组，如下，该路由组，只要以 admin 为 url 前缀的，就会用到 BasicAuth 中间件
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		// BasicAuth 方法是以基础的 http auth 作为权限验证
		// gin.Accounts 是 map[string]string 的一种快捷方式，形象意义上，包含了账号、密码
		// BasicAuth 方法需要传递的参数，也是一组用户名和密码，用于校验登陆
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// 路由：/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string) // 获取到请求中字段为 user 的值，并返回其字符串形式，如果不存在，则会报错

		// 模拟数据表中的数据
		var secrets = gin.H{
			"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
			"austin": gin.H{"email": "austin@example.com", "phone": "666"},
			"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
		}

		// 在 authorized 路由组中已经通过 BasicAuth 中间件校验了登陆
		// 当下要校验的是账号在数据库中是否存在
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// 路由定义时可以使用 http 方法
	// r.GET("/someGet", getting)
	// r.POST("/somePost", posting)
	// r.PUT("/somePut", putting)
	// r.DELETE("/someDelete", deleting)
	// r.PATCH("/somePatch", patching)
	// r.HEAD("/someHead", head)
	// r.OPTIONS("/someOptions", options)
	// r.Any("/someOptions", options)

	// 使用全局中间件
	// Logger 中间件将日志写入 gin.DefaultWriter，即使 GIN_MODE 设置为 release
	// 中间件可以注册两遍，不会报错，但会执行两次，并且在 GIN-DEBUG 中的路由列表中，handlers 也会按两次来算
	// 就算后面使用 r.Use 注册了新的中间件，但只生效于后面注册的路由，前面已经注册的路由无效
	//r.Use(gin.Logger())
	// Recovery 中间件会  recover 任何 panic，并写入 500
	// r.Use(gin.Recovery())
	// 可以为一个路由单独定义中间件
	r.GET("/benchmark", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	}, func(c *gin.Context) {
		fmt.Println(666)
	})

	// 认证路由组
	// authorized := r.Group("/", AuthRequired())
	// 和使用以下两行代码的效果完全一样:
	auth := r.Group("/")
	// 注册中间件
	auth.Use(func(c *gin.Context) {
		fmt.Println("auth test")
	})

	// 这是个代码块，文档上因为和上面代码处于同一行，容易让人误会
	{
		auth.GET("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"ok": "ok",
			})
		}, func(c *gin.Context) {
			fmt.Println("路由组中间件下的定制中间件")
		})

		// 嵌套路由组
		testing := auth.Group("testing")
		testing.GET("/analytics", func(c *gin.Context) {
			c.JSON(200, gin.H{})
		})
	}

	r.Any("/onlyQuery", func(c *gin.Context) {
		type Person struct {
			Name    string `form:"name"`
			Address string `form:"address"`
		}

		var person Person

		// ShouldBindQuery 方法会将请求中 url query 部分的参数绑定到传入的结构体中，会忽略 post 的表单数据
		// 但 POST 请求的 url query 还是能够绑定的
		if c.ShouldBindQuery(&person) == nil {
			log.Println("====== Only Bind By Query String ======")
			log.Println(person.Name)
			log.Println(person.Address)
		}
		c.String(200, "Success")
	})

	// 中间件中使用 goroutine 必须为上下文创建一个副本
	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			// 用 time.Sleep() 模拟一个长任务。
			time.Sleep(5 * time.Second)

			// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 因为没有使用 goroutine，不需要拷贝上下文
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// QueryMap 和 PostFormMap
	r.POST("/QueryMapAndPostFormMap", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")
		c.JSON(200, gin.H{
			"ids":   ids,
			"names": names,
		})
	})

	// 查询字符串
	// /welcome?firstname=Jane&lastname=Doe
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest") // 带默认值
		lastname := c.Query("lastname")                   // c.Request.URL.Query().Get("lastname") 的一种快捷方式

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	// gin 默认只支持单模板，如果要使用多模板，需要用到另外的包

	// 一般想要将请求体绑定到结构体中时，会使用 ShouldBind 方法
	// 但使用 ShouldBind 时，只有 Query, Form, FormPost, FormMultipart 这些格式才能多次调用此方法进行多次绑定
	// 甚至是绑定到不同的结构体且不造成任何性能损失
	// 但如果要对 JSON, XML, MsgPack, ProtoBuf 的请求体格式进行绑定时，ShouldBind 方法只能绑定一次
	// 二次调用时，会找不到 Request.Body
	// 这种情况下，要使用 ShouldBindBodyWith 来绑定，这个方法在绑定之前，会将 body 存储到上下文中
	// 这个存储工作，会有额外的性能开销，虽然非常轻微，所以，能使用 ShouldBind 一次绑定时，就尽可能不使用 ShouldBindBodyWith 方法

	// gin 中将请求内容绑定到结构体，一般有两类方法，分别为 Must 和 Should
	// 区别是，Must 绑定时，如果绑定失败，则直接抛出 4xx 的响应，而 Should 绑定出错时，会将错误返回给开发者自己处理
	// 结构体中，可以通过为字段定义 `` 标签来具体指导绑定方式，如 `form:"user" json:"user" xml:"user" binding:"required"`
	// 其中 binding 如果指定为 required ，则请求中必须有值才能绑定成功，否则绑定失败
	// binding 还可以设定为 "-" 这样就算请求中不包含此字段，也能绑定成功

	type myForm struct {
		Colors []string `form:"colors[]"` // 这个结构体，可以将请求中的复选框的值绑定到 Colors 字段
	}

	type Person struct {
		ID   string `uri:"id" binding:"required,uuid"` // 其中的 uri 指将对应的路由参数绑定到结构体来，但是必须使用 ShouldBindUri 方法
		Name string `uri:"name" binding:"required"`
	}

	// 使用 c.Bind 来将请求绑定到结构体时，如果结构体有嵌套的情况，则嵌套字段不可有 form tag，否则无法绑定

	r.GET("/cookie", func(c *gin.Context) {

		// 先看 cookie 是否存在
		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			// 不存在时，设置 cookie
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "127.0.0.1", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})

	// 此 middlewares 将匹配 /user/john 但不会匹配 /user/ 或者 /user
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// 此 middlewares 将匹配 /user/john/ 和 /user/john/send
	// 如果没有其他路由匹配 /user/john，它将重定向到 /user/john/
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com/") // 重定向

		// 也可以这么写
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})

	// gin.DisableConsoleColor() 可以禁用控制台颜色
	// gin.ForceConsoleColor() 强制开启控制台颜色
	// goland 的 terminal 看不出颜色，要在系统自带的命令行上看

	// http.ListenAndServe(":8080", router) 可以使用此方法自定义 http 配置

	//go r.Run() // 监听并在 0.0.0.0:8080 上启动服务
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务

	// 如果要运行多个服务，请参考 gin 官方文档的写法

}
