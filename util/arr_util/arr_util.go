package arr_util

func Remove[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice // 索引越界，不进行删除
	}
	return append(slice[:index], slice[index+1:]...) // 拼接前后部分
}
