package DBHelper

import(
	"database/sql"
	"reflect"
	"strings"
	"errors"
)

type DB struct{
	DbDriver string
	ConnectionString string
}


func (db DB) Query(sqlText string,T *interface{}) (error){
	var structType reflect.Type = nil;
	s := nil
	
	
	outType := reflect.TypeOf(T).Elem()
	switch(outType.Kind()){
	case reflect.Invalid,reflect.Uintptr,reflect.Complex64,reflect.Complex128,reflect.Chan,reflect.Func,reflect.Ptr,reflect.UnsafePointer://nonsupported types
		return errors.New("type error")
	case reflect.Slice,reflect.Array://if slice will returns all rows.if array will returns rows that count <= len(array) 
		itemType := outType.Elem()
		switch itemType.Kind() {
			case reflect.Invalid,reflect.Uintptr,reflect.Complex64,reflect.Complex128,reflect.Chan,reflect.Func,reflect.Ptr,reflect.UnsafePointer,reflect.Slice,reflect.Array://nonsupported types
				return errors.New("type error")
			case reflect.Struct:
				structType = itemType
			case reflect.Map:
			default:
		}
	case reflect.Struct: //will save first row into this struct
		structType = outType
	case reflect.Map://will save first row into this map
	default://will save first row's first column
	}
	
	conn,err := sql.Open(db.DbDriver,db.ConnectionString)
	if(err==nil){
		return err;
	}
	rows,err := conn.Query(sqlText)
	if(err!=nil){
		return err;
	}
	args := buildArgsByDbAndStruct(structType,rows)
	for(rows.Next()){
		rows.Scan(args...)

	}


	reflect.ValueOf(T).Field(0)
}

func buildArgsByDbAndStruct(t reflect.Type,rows *sql.Rows) []interface{} {
	var args []interface{}
	cols,err := rows.Columns()
	if(err!=nil){
		return args
	}
	args = make([]interface{},len(cols))
	argsValue := reflect.ValueOf(&args)
	for i,field := range(cols) {
		fieldType := getTypeByName2(t,field)
		argsValue.Elem().Index(i).Set(reflect.New(fieldType))
	}
	return args;
}

func saveValueToStruct(args interface{},rows *sql.Rows,t reflect.Type)(error){
	cols,err := rows.Columns()
	if(err!=nil){
		return err;
	}
	val := reflect.New(t)
	for i,field := range cols {
		fieldType := getTypeByName2(t,field)
		val.FieldByName("")
	}
	return nil
}

func getTypeByName(s *interface{},name string) (reflect.Type){
	return nil;
}

func getTypeByName2(t reflect.Type,name string) (reflect.Type){
	var sencendType reflect.Type
	len := t.NumField();
	for i:=0;i<len;i++ {
		field := t.Field(i).Tag.Get("field")
		if(field == name){ 
			return t.Field(i).Type;
		}
		field = strings.ToLower(t.Field(i).Name)
		if(field == strings.ToLower(name)){
			sencendType = t.Field(i).Type
		}
	}
	return sencendType
}

func setStructValue(structValue reflect.Value,name string,fieldValue interface{}){
	var firstIndex int = -1
	var secondIndex int = -1
	lowName := strings.ToLower(name)
	len := structValue.NumField();
	t := structValue.Type()
	for i:=0;i<len;i++ {
		field := t.Field(i).Tag.Get("field")
		if(field==""||field=="_"||field=="-"){
			field = t.Field(i).Name
		}
		if(field==name) {
			firstIndex = i
			break;
		} else {
			field = strings.ToLower(field)
			if(field==lowName){
				secondIndex = i
			}
		}
	}
	if firstIndex>-1 {
		structValue.Field(firstIndex).Set(reflect.ValueOf(fieldValue))
	} else if secondIndex>-1 {
		structValue.Field(secondIndex).Set(reflect.ValueOf(fieldValue))
	}
}