package DBHelper

import(
	
)

type DB struct{
	DbDriver string
	ConnectionString string
}

func (db DB) Query(sql string,T *interface{}) {

}