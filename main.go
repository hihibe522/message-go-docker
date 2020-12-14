package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	// "net/http"
)

type Msg struct {
	Msg_id int    `json:"msg_id"`
	Name   string `json:"name"`
	Msg    string `json:"msg"`
	Time   string `json:"time"`
}

var DB *sql.DB

func main() {

	initDB()
	Router()
}
func initDB() {

	DB, _ = sql.Open("mysql", "root:root@tcp(mysql)/message")
	// 設定 database 最大連接數
	DB.SetConnMaxLifetime(100)
	//設定上 database 最大閒置連接數
	DB.SetMaxIdleConns(10)
	// 驗證是否連上 db
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail:", err)
		return
	}
	fmt.Println("connnect success")
}

func Router() {
	router := gin.Default()
	router.LoadHTMLGlob("dist/*.html")        // 添加入口index.html
	router.LoadHTMLFiles("static/*/*")        // 添加资源路径
	router.Static("/static", "./dist/static") // 添加资源路径
	router.StaticFile("/", "dist/index.html") //前端接口
	router.GET("/api/board", InitPage)
	router.POST("/api/board", AddMsg)
	router.DELETE("/api/board", DeleteMsg)
	router.PATCH("/api/board", UpdateMsg)
	// router.GET("/board", InitPage)
	// router.POST("/board", AddMsg)
	// router.DELETE("/board", DeleteMsg)
	// router.PATCH("/board", UpdateMsg)

	router.Run(":3000")
}

func InitPage(c *gin.Context) {
	msg := new(Msg)
	Msgs := []Msg{}
	rows, err := DB.Query("SELECT* FROM `msg`")
	defer rows.Close()
	if err != nil {
		fmt.Print(err.Error())
	}
	// 開始輪查
	for rows.Next() {
		err = rows.Scan(&msg.Msg_id, &msg.Name, &msg.Msg, &msg.Time)
		Msgs = append(Msgs, *msg)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"result": Msgs,
		"count":  len(Msgs),
	})

}

func AddMsg(c *gin.Context) {
	msg := &Msg{}
	if err := c.BindJSON(&msg); err != nil {
		fmt.Println(err)
		return
	}

	rs, err := DB.Exec("INSERT INTO `msg`(`name`,`msg`) VALUES (?,?)", msg.Name, msg.Msg)
	if err != nil {
		fmt.Println(err)
	}

	rowCount, err := rs.RowsAffected()
	rowId, err := rs.LastInsertId() // 資料表中有Auto_Increment欄位才起作用，回傳剛剛新增的那筆資料ID

	fmt.Printf("新增 %d 筆資料，id = %d \n", rowCount, rowId)
	c.JSON(http.StatusOK, gin.H{
		"status": " insert ok",
	})

}

func UpdateMsg(c *gin.Context) {
	msg := &Msg{}
	if err := c.BindJSON(&msg); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg.Msg_id, msg.Name, msg.Msg)
	_, err := DB.Exec("UPDATE `msg` SET `name`= ?, `msg`=? WHERE `msg_id` = ?;", msg.Name, msg.Msg, msg.Msg_id)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "update ok",
	})
}

func DeleteMsg(c *gin.Context) {
	msg_id := c.Query("msg_id")
	_, err := DB.Exec("DELETE FROM `msg` WHERE `msg_id` = ?;", msg_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("delete ok")
	c.JSON(http.StatusOK, gin.H{
		"status": "delete ok",
	})
}
