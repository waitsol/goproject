package jfzt

func SendMsg(id, msg string, send_group bool) {

	if MGR[0].start {
		m := map[string]interface{}{}
		m["id"] = id
		if send_group {
			m["group_id"] = "853312133"
		}
		im.SendMsg(msg, m)
	}

}
