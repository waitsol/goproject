package main

import "strconv"

var mNameFollor map[string]*Follow

type Follow struct {
	Id   map[string]int //关注的股票id
	Name string
}

func (this *Follow) follow(id string) {
	if _, ok := this.Id[id]; ok {
		Post(this.Name, id)
		this.Id[id] = 400
	} else {
		//退出监听
		delete(mId2Listener[id], this.Name)
		delete(mNameFollor, this.Name)
	}
}

func (this *Follow) HandleMessage(msg string) (bool, string) {
	v := stringSplit(msg, ' ')
	if len(v) == 0 {
		return false, ""
	}
	if checkGpNum(v[0]) {
		this.follow(v[0])
		return true, "苟富贵"
	} else if v[0] == "set" {
		if len(v) == 3 && checkGpNum(v[1]) {
			x, err := strconv.Atoi(v[2])
			if err != nil || x < 200 {
				return true, "err " + v[2]
			}
			_, ok := this.Id[v[1]]
			if !ok {
				return true, "err " + v[1]
			}
			this.Id[v[1]] = x
			return true, "ok"
		} else {
			return true, "err args"
		}
	}
	return false, ""
}
func getFoller(name string) *Follow {
	if mNameFollor == nil {
		mNameFollor = map[string]*Follow{}
	}
	if x, ok := mNameFollor[name]; ok {
		return x
	}
	x := &Follow{Name: name}
	mNameFollor[name] = x
	return x
}

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
