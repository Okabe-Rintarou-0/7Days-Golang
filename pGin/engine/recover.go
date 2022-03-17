package engine

import (
	"fmt"
	"runtime"
	"strings"
)

func trace(errorMsg string) string {
	var pcs [32]uintptr

	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(errorMsg)
	str.WriteString("Trace back\n")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\tat %s:%d\n\t\t%s\n", file, line, fn.Name()))
	}

	return str.String()
}

func recoverPanic() {
	if err := recover(); err != nil {
		fmt.Println(trace(fmt.Sprintf("%v\n", err)))
	}
}
