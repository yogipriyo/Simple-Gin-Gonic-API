package main
import (
 	"github.com/gin-gonic/gin"
 	"database/sql"
 	"gopkg.in/gorp.v1"
 	_ "github.com/lib/pq"
 	"strconv"
 	"log"
)
////

var dbmap = initDb()

func initDb() *gorp.DbMap {
 	db, err := sql.Open("postgres", "postgres://postgres:09030015@localhost/gopgtest")
 	checkErr(err, "sql.Open failed")
 	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
 	dbmap.AddTableWithName(User{}, "User").SetKeys(true, "Id")
 	err = dbmap.CreateTablesIfNotExists()
 	checkErr(err, "Create table failed")
	return dbmap
}

func checkErr(err error, msg string) {
 	if err != nil {
 		log.Fatalln(msg, err)
 	}
}
////

func GetUsers(c *gin.Context) {
 	var users []User
 	_, err := dbmap.Select(&users, "SELECT * FROM user")
	if err == nil {
 		c.JSON(200, users)
 	} else {
 		c.JSON(404, gin.H{"error": "no user(s) into the table"})
 	}
	// curl -i http://localhost:8080/api/v1/users
}

func GetUser(c *gin.Context) {
 	id := c.Params.ByName("id")
 	var user User
 	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
	if err == nil {
 		user_id, _ := strconv.ParseInt(id, 0, 64)
		content := &User{
 			Id: user_id,
 			Firstname: user.Firstname,
 			Lastname: user.Lastname,
 		}
 		c.JSON(200, content)
 	} else {
 	c.JSON(404, gin.H{"error": "user not found"})
 	}
	// curl -i http://localhost:8080/api/v1/users/1
}

func PostUser(c *gin.Context) {
 	var user User
 	c.Bind(&user)
	if user.Firstname != "" && user.Lastname != "" {
		if insert, _ := dbmap.Exec(`INSERT INTO user (firstname, lastname) VALUES (?, ?)`, user.Firstname, user.Lastname); insert != nil {
		 	user_id, err := insert.LastInsertId()
		 	if err == nil {
		 		content := &User{
		 			Id: user_id,
		 			Firstname: user.Firstname,
		 			Lastname: user.Lastname,
				}
		 		c.JSON(201, content)
			} else {
		 		checkErr(err, "Insert failed")
		 	}
	 	}
	} else {
 		c.JSON(422, gin.H{"error": "fields are empty"})
 	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func UpdateUser(c *gin.Context) {
 	id := c.Params.ByName("id")
 	var user User
 	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
	if err == nil {
 		var json User
 		c.Bind(&json)
		user_id, _ := strconv.ParseInt(id, 0, 64)
		user := User{
	 		Id: user_id,
	 		Firstname: json.Firstname,
	 		Lastname: json.Lastname,
	 	}
		if user.Firstname != "" && user.Lastname != ""{
	 		_, err = dbmap.Update(&user)
			if err == nil {
	 			c.JSON(200, user)
	 		} else {
	 			checkErr(err, "Updated failed")
	 		}
		} else {
	 		c.JSON(422, gin.H{"error": "fields are empty"})
	 	}
	} else {
 		c.JSON(404, gin.H{"error": "user not found"})
 	}
	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeleteUser(c *gin.Context) {
 	id := c.Params.ByName("id")
	var user User
 	err := dbmap.SelectOne(&user, "SELECT id FROM user WHERE id=?", id)
	if err == nil {
	 	_, err = dbmap.Delete(&user)
		if err == nil {
			c.JSON(200, gin.H{"id #" + id: " deleted"})
	 	} else {
	 		checkErr(err, "Delete failed")
		}
	} else {
 		c.JSON(404, gin.H{"error": "user not found"})
 	}
	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}

////
type User struct {
 	Id int64 `db:"id" json:"id"`
 	Firstname string `db:"firstname" json:"firstname"`
 	Lastname string `db:"lastname" json:"lastname"`
}

func main() {
 	r := gin.Default()
	v1 := r.Group("api/v1")
 	{
 		v1.GET("/users", GetUsers)
 		v1.GET("/users/:id", GetUser)
 		v1.POST("/users", PostUser)
 		v1.PUT("/users/:id", UpdateUser)
 		v1.DELETE("/users/:id", DeleteUser)
 	}
	r.Run(":8081")
}


//static example//
// func GetUsers(c *gin.Context) {
//  	type Users []User
// 	var users = Users{
//  		User{Id: 1, Firstname: "Oliver", Lastname: "Queen"},
//  		User{Id: 2, Firstname: "Malcom", Lastname: "Merlyn"},
//  	}
// 	c.JSON(200, users)
// 	// curl -i http://localhost:8080/api/v1/users
// }

// func GetUser(c *gin.Context) {
//  	id := c.Params.ByName("id")
//  	user_id, _ := strconv.ParseInt(id, 0, 64)
// 	if user_id == 1 {
//  		content := gin.H{"id": user_id, "firstname": "Oliver", "lastname": "Queen"}
//  		c.JSON(200, content)
//  	} else if user_id == 2 {
//  		content := gin.H{"id": user_id, "firstname": "Malcom", "lastname": "Merlyn"}
//  		c.JSON(200, content)
//  	} else {
//  		content := gin.H{"error": "user with id#" + id + " not found"}
//  		c.JSON(404, content)
//  	}
// 	// curl -i http://localhost:8080/api/v1/users/1
// }

// func PostUser(c *gin.Context) {
//  // The futur code…
// }
// func UpdateUser(c *gin.Context) {
//  // The futur code…
// }
// func DeleteUser(c *gin.Context) {
//  // The futur code…
// }