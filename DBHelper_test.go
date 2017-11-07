package DBHelper

import(
	"fmt"
	"testing"
	_ "github.com/mattn/go-adodb"
	_ "encoding/json"
)

type Goods struct{
	Id string
	Name string
}


func TestQuery(t * testing.T){
	
	db := DB{"adodb","Provider=SQLOLEDB;Data Source=.\\sqlexpress;Integrated Security=SSPI;Initial Catalog=mybusiness;user id=sa;password=123456"}
	
	var g1 Goods
	db.Query(&g1,"select * from goods where id=?","00030004")
	fmt.Println(g1)

	fmt.Println("--------------------------------------------------")

	var g2 [3]Goods
	db.Query(&g2,"select name,id from goods")
	fmt.Println(g2)

	fmt.Println("--------------------------------------------------")

	var g3 []Goods
	db.Query(&g3,"select * from goods")
	fmt.Println(g3)
}