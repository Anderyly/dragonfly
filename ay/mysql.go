package ay

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Sql() {
	dsn := Yaml.GetString("mysql.user") + ":" + Yaml.GetString("mysql.password") + "@tcp(" + Yaml.GetString("mysql.localhost") + ":" + Yaml.GetString("mysql.port") + ")/" + Yaml.GetString("mysql.database") + "?parseTime=true&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
	}
	Db = database

	//database.SingularTable(true)
	// 把模型与数据库中的表对应起来
	//Db.AutoMigrate(&models.Area{})
}
