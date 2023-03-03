package main

// 分割文本字符串
func stringSplit(text string, ic byte) []string {
	vs := []string{}

	b := -1
	for i, c := range text {
		if c == rune(ic) {
			if b != -1 {
				vs = append(vs, string(text[b:i]))
				b = -1
			}
		} else if b == -1 {
			b = i
		}
	}
	if b != -1 {
		vs = append(vs, string(text[b:]))
	}
	return vs
}
func HandleMessage(msg string) bool {
	v := stringSplit(msg, ' ')
	if len(v) == 0 {
		return false
	}
	if v[0] == "a" {
		return add(v)
	} else if v[0] == "d" {
		return del(v)
	}
	return false

}
func checkGpNum(s string) bool {
	if len(s) != 6 {
		return false
	}
	for _, c := range s {
		if c > '9' || c < '0' {
			return false
		}
	}
	return true
}

func add(v []string) bool {
	if len(v) == 3 {
		if checkGpNum(v[1]) {
			data := Data_json{
				SubType:     "SUBON",
				Inst:        v[1],
				Market:      v[2],
				ServiceType: "STATISTICS",
				ReqID:       1,
			}
			Post(data)

			data = Data_json{
				SubType:     "SUBON",
				Inst:        v[1],
				Market:      v[2],
				ServiceType: "STATIC",
				ReqID:       1,
			}
			Post(data)
			data = Data_json{
				SubType:     "SUBON",
				Inst:        v[1],
				Market:      v[2],
				ServiceType: "TICK",
				ReqID:       1,
			}
			Post(data)

		}
	}

	return true

}

func del(v []string) bool {
	if len(v) >= 2 {
		if checkGpNum(v[1]) {

		}
	}
	return false
}
