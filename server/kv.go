package server

type kvcache struct{
	kv map[string]string
}

var Cache *kvcache

func (cache *kvcache) SET(key string,val string){
	cache.kv[key] = val
}

func (cache * kvcache) GET(key string) string {
	if _, ok := cache.kv[key]; ok {
		return cache.kv[key]
	}
	return ""
}

func (cache *kvcache) EXIS(key string)bool{
	if _, ok := cache.kv[key]; ok {
		return true
	}
	return false
}

func (cache *kvcache) GET_KEYS()[]string{
	var res []string
	for k , _  := range cache.kv {
		res = append(res,k)
	}
	return res
}

func (cache *kvcache) DELETE(key string){
	delete(cache.kv, key)
}

func InitCache()*kvcache{
	Cache = &kvcache{
		make(map[string]string),
	}
	return Cache
}