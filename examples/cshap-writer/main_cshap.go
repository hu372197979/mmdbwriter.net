package main
/*
#include<stdlib.h> //导入 stdlib.h 文件头以使用 C.free 方法
*/
import "C"
import (
	"log"
	"net"
	"os"
	"unsafe"
	"runtime"
	"encoding/json"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func main() {
}

var (
	glb_writer *mmdbwriter.Tree
	err error
)

//export Init
func Init(dbfile *C.char){

	defer func() {
		if err := recover(); err != nil {
			log.Println("Init: " + err.(string))
		}
	}()

	c_dbfile :=C.GoString(dbfile)
	//defer C.free(unsafe.Pointer(dbfile))
	
	//默认新建文件当文件存在则加载文件
    isLoad := false 
	if(c_dbfile != ""){
		_, _bool := IsExists(c_dbfile)
		if (_bool) {
			isLoad = true;
		}
	}

	//判断是加载文件还是新建文件
	if (isLoad) {
		glb_writer, err = mmdbwriter.Load(c_dbfile, mmdbwriter.Options{})
	}else {
		glb_writer, err = mmdbwriter.New(
			mmdbwriter.Options{
				DatabaseType: "Spring.IPLocation.City",
				RecordSize:  28,
				Languages: []string{"en", "zh-CN"},
				Description:map[string]string{"en": "Spring.IPLocation.City","zh-CN":"IP城市定位数据库"},
			},
		)
	}
	runtime.SetFinalizer(&glb_writer,Finalize)
	if err != nil {
		log.Println(err)
	}
	return
}

//export InsertRange
func InsertRange(start *C.char,end *C.char,value *C.char){

	defer func() {
		if err := recover(); err != nil {
			log.Println("InsertRange: " + err.(string))
		}
	}()

	var p interface{}
	c_json :=C.GoString(value)
	c_start :=C.GoString(start)
	c_end :=C.GoString(end)

	//defer func() {
		//C.free(unsafe.Pointer(start))
		//C.free(unsafe.Pointer(end))
		//C.free(unsafe.Pointer(value))
	//}()
 	json.Unmarshal([]byte(c_json), &p)
	_record := MapData(p.(map[string]interface{}))
	_start := net.ParseIP(c_start)
	_end := net.ParseIP(c_end)
	err = glb_writer.InsertRange(_start,_end,_record)
	if err != nil {
		log.Println(err)
	}
	return
}

//export Save
func Save(filename *C.char){

	defer func() {
		if err := recover(); err != nil {
			log.Println("Save: " + err.(string))
		}
	}()

	c_filename :=C.GoString(filename)
	//defer C.free(unsafe.Pointer(filename))

	fh, err := os.Create(c_filename)
	if err != nil {
		log.Println(err)
	}
	_, err = glb_writer.WriteTo(fh)
	if err != nil {
		log.Println(err)
	}
	return
}

//export Free
func Free(p unsafe.Pointer) {
	C.free(p)
}

//export Relase
func Relase(){
	runtime.GC()
}

func MapData(valueInput map[string]interface{}) mmdbtype.DataType{
	record := mmdbtype.Map{}
	for k, v := range valueInput {
		switch value := v.(type)  {
        case nil:
            log.Println(k, "is nil", "null")
        case string:
            record[mmdbtype.String(k)] = mmdbtype.String(v.(string))
        case int32:
            record[mmdbtype.String(k)] = mmdbtype.Int32(v.(int32))
        case float64:
            record[mmdbtype.String(k)] = mmdbtype.Float64(v.(float64))
        case []interface{}:
            log.Println(k, "is an array:")
            for i, u := range value {
                log.Println(i, u)
            }
        case map[string]interface{}:
            record[mmdbtype.String(k)]= MapData(v.(map[string]interface{}))
        default:
            log.Println(k, "is unknown type")
        }
	}
	return record
}

func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
    return f, err == nil || os.IsExist(err)
}

func Finalize(){
	log.Println("Finalize!!!!!!! ")
}