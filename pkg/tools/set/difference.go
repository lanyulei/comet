package set

/*
  @Author : lanyulei
  @Desc :
*/

// Difference 切片差集
func Difference(val1, val2 []interface{}) (d []interface{}) {
	tmp := map[interface{}]struct{}{}

	for _, v := range val2 {
		if _, ok := tmp[v]; !ok {
			tmp[v] = struct{}{}
		}
	}

	for _, v := range val1 {
		if _, ok := tmp[v]; !ok {
			d = append(d, v)
		}
	}

	return
}
