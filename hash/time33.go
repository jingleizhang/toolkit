package hash

func Time33(buf string, shards uint32) uint32 {
	var h, p uint64

	for i := 0; i < len(buf); i++ {
		p = uint64((buf)[i])
		h = (h + (h << 5)) + p
	}

	return uint32(h % uint64(shards))
}
