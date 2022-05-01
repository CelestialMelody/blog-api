package util

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(value string) string {
	m := md5.New()                        //创建一个md5算法的hash对象
	m.Write([]byte(value))                //写入要处理的字符串
	return hex.EncodeToString(m.Sum(nil)) //返回一个[]byte类型的hash值
}
