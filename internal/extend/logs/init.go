// @Author xuanshuiyuan 2023/4/20 15:10:00
package logs

type LogsEr interface {
	Add()
}

//记录日志方法
type MeansEr interface {
	Add(*LogsParams) error
	SetOperation(string, ...interface{}) string
}

func Add(l LogsEr) {
	l.Add()
}
