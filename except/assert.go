package except

import "runtime/debug"

func ASSERT(exp bool, info ...string) { // 接受一个字符串参数
	if !exp {
		infostr := ""
		if len(info) > 0 {
			infostr = info[0]
		}
		Log.Errorf("ASSERT FAILED!\ninfo=[%v]\nstack = [%v]\n", infostr, string(debug.Stack()))
		panic("ASSERT FAILED")
	}
}
