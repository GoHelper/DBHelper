package DBHelper

import(
	"database/sql"
	"reflect"
	"strings"
	"errors"
	//"fmt"
)

type DB struct{
	DbDriver string
	ConnectionString string
}

func (db DB) Exec(sqlText string) error {
	conn,err := sql.Open(db.DbDriver,db.ConnectionString)
	if(err!=nil){
		return err;
	}
	_,err = conn.Exec(sqlText)
	if(err!=nil){
		return err;
	}
	return nil
}

func (db DB) Query(T interface{},sqlText string,sqlArgs ...interface{}) (error){
	var structType reflect.Type = nil;
	var mapType reflect.Type = nil;
	var singleValueType reflect.Type = nil;
	var dataCount int = 0;
	outType := reflect.TypeOf(T).Elem()
	outTypeKind := outType.Kind()
	switch(outTypeKind){
	case reflect.Invalid,reflect.Uintptr,reflect.Complex64,reflect.Complex128,reflect.Chan,reflect.Func,reflect.Ptr,reflect.UnsafePointer://unsupported types
		return errors.New("type error")
	case reflect.Slice,reflect.Array://if slice will returns all rows.if array will returns rows that count <= len(array) 
		if(outTypeKind==reflect.Slice){
			dataCount = -1
		}else{
			dataCount = outType.Len()
		}
		
		itemType := outType.Elem()
		switch itemType.Kind() {
			case reflect.Invalid,reflect.Uintptr,reflect.Complex64,reflect.Complex128,reflect.Chan,reflect.Func,reflect.Ptr,reflect.UnsafePointer,reflect.Slice,reflect.Array://nonsupported types
				return errors.New("type error")
			case reflect.Struct:
				structType = itemType
			case reflect.Map:
				mapType = itemType
			default:
				singleValueType = itemType
		}
	case reflect.Struct: //will save first row into this struct
		dataCount = 1
		structType = outType
	case reflect.Map://will save first row into this map
		dataCount = 1
		mapType = outType
	default://will save first row's first column
		dataCount = 1
		singleValueType = outType
	}
	conn,err := sql.Open(db.DbDriver,db.ConnectionString)
	if(err!=nil){
		return err;
	}
	rows,err := conn.Query(sqlText,sqlArgs...)
	if(err!=nil){
		return err;
	}
	var args []interface{}
	if(structType!=nil){
		args = buildArgsByDbAndStruct(structType,rows)
	}else if(mapType!=nil){
		args = buildArgsByDbAndType(mapType.Elem(),rows)
	}else if(singleValueType != nil){
		args = buildArgsByDbAndType(singleValueType,rows)
	}
	i :=-1
	for(rows.Next()){
		i++;
		if(dataCount>=0&&i>=dataCount){
			break
		}
		rows.Scan(args...)
		if(structType!=nil){
			structValue,err := saveValueToStruct(args,rows,structType)
			if(err!=nil) {
				return err
			}
			switch outTypeKind {
			case reflect.Struct:
				reflect.ValueOf(T).Elem().Set(structValue)
			case reflect.Array:
				reflect.ValueOf(T).Elem().Index(i).Set(structValue)
			case reflect.Slice:
				s := reflect.ValueOf(T).Elem()
				s.Set(reflect.Append(s,structValue))
			}
		} else if mapType!=nil {
			// keyType := mapType.Key()
			// for(c:=0;c<rows.Columns())
			// keyValue := reflect.Zero(keyType)
			// keyValue.SetString()
			// switch outTypeKind {
			// case reflect.Map:
			// 	reflect.ValueOf(T).Elem().Set(structValue)
			// case reflect.Array:
			// 	reflect.ValueOf(T).Elem().Index(i).Set(structValue)
			// case reflect.Slice:
			// 	reflect.Append(reflect.ValueOf(T).Elem(),structValue)	
			// }
		}
	}
	return nil
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

func buildArgsByDbAndType(valueType reflect.Type,rows *sql.Rows) []interface{} {
	var args []interface{}
	cols,err := rows.Columns()
	if(err!=nil){
		return args
	}
	args = make([]interface{},len(cols))
	argsValue := reflect.ValueOf(&args)
	for i,_ := range(cols) {
		argsValue.Elem().Index(i).Set(reflect.New(valueType))
	}
	return args;
}

// function saveValutToMap(args[]interface{},rows *sql.Rows,mapType reflect.Type)(reflect.Value,error){
// 	val := reflect.Make()
// }

func saveValueToStruct(args []interface{},rows *sql.Rows,structType reflect.Type)(reflect.Value,error){
	var val reflect.Value
	if(val.IsValid()){
		val.Set(reflect.Zero(structType))
	}else{
		val = reflect.New(structType).Elem()
	}
	cols,err := rows.Columns()
	if(err!=nil){
		return val,err;
	}
	
	for i,field := range cols {
		setStructValue(val,field,args[i])
	}
	return val,nil
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
		structValue.Field(firstIndex).Set(reflect.ValueOf(fieldValue).Elem())
	} else if secondIndex>-1 {
		structValue.Field(secondIndex).Set(reflect.ValueOf(fieldValue).Elem())
	}
}