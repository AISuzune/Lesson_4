package dao

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// 假数据库，用 map 实现
var database = map[string]string{
	"yxh": "123456",
	"wx":  "654321",
}

var lock = sync.RWMutex{} // lock 是一个互斥锁，用于并发安全操作数据库

// init 函数在程序启动时调用，从文件中加载数据到内存中
func init() {
	loadDataFromFile()
}

// loadDataFromFile 从文件中读取数据到内存中
func loadDataFromFile() {
	file, err := os.Open("database.json") // 打开名为 database.json 的文件
	if err != nil {
		return
	}
	defer file.Close() // 关闭文件

	data, err := ioutil.ReadAll(file) // 读取文件中的数据
	if err != nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	err = json.Unmarshal(data, &database) // 反序列化数据到 map 变量中
	if err != nil {
		return
	}
}

// saveDataToFile 将内存中的数据写入到文件中
func saveDataToFile() {
	lock.RLock()
	defer lock.RUnlock()

	data, err := json.Marshal(database) // 将 map 变量序列化为 JSON 字符串
	if err != nil {
		return
	}

	err = ioutil.WriteFile("database.json", data, 0644) // 将 JSON 字符串写入到文件中
	if err != nil {
		return
	}
}

// AddUser 添加一个新用户到数据库中
func AddUser(username, password string) {
	lock.Lock()
	defer lock.Unlock()

	database[username] = password // 将新用户信息添加到 map 变量中

	saveDataToFile() // 将新数据写入到文件中
}

// SelectUser 根据用户名查找用户是否存在
func SelectUser(username string) bool {
	lock.RLock()
	defer lock.RUnlock()

	_, ok := database[username] // 判断用户名是否存在

	return ok
}

// SelectPasswordFromUsername 根据用户名获取用户密码
func SelectPasswordFromUsername(username string) string {
	lock.RLock()
	defer lock.RUnlock()

	return database[username] // 返回用户名对应的密码
}

/*func AddUser(username, password string) {
	database[username] = password
}

// 若没有这个用户返回 false，反之返回 true

func SelectUser(username string) bool {
	if database[username] == "" {
		return false
	}
	return true
}

func SelectPasswordFromUsername(username string) string {
	return database[username]
}*/
