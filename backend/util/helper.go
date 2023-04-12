package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/dop251/goja"
)

func WriteError(level string, error_string string) {
	log.Println(error_string)
}

func ExecJsXbogus(url_string string, user_agent string) (string, error) {

	content, err1 := ioutil.ReadFile("./resource/Signer.js")
	if err1 != nil {
		return "", err1
	}
	vm := goja.New()
	_, err := vm.RunString(string(content))
	if err != nil {
		fmt.Println("JS代码有问题！")
		return "", err
	}
	var fn func(string, string) string
	err = vm.ExportTo(vm.Get("sign"), &fn)
	if err != nil {
		fmt.Println("Js函数映射到 Go 函数失败！")
		return "", err
	}
	sign := fn(url_string, user_agent)

	return sign, nil
}

func DouyinMsToken(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
