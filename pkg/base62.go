package pkg

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode_base62(id uint64) string {
	var result []byte
	for id > 0 {
		remainder := id % 62
		result = append(result, alphabet[remainder])
		id /= 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
