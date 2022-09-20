package util

var Hosts []Host

const FileName = "test.properties"
const ConfigFilePath = "/root/test.properties"

type Host struct {
	Ip    string
	State int
}
