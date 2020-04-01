package server

type kchan struct{
	kc map[string]chan string
}

var DeviceChan *kchan
var ClientChan *kchan

func (kc *kchan)NewChan(key string)chan string{
	ch := make(chan string)
	kc.kc[key] = ch
	return ch
}

func (kc *kchan)GetChan(key string)chan string{
	return kc.kc[key]
}

func (kc *kchan) Exist(key string)bool{
	if _, ok := kc.kc[key]; ok {
		return true
	}
	return false
}

func (kc *kchan) GetKeys()[]string{
	var res []string
	for k , _  := range kc.kc {
		res = append(res,k)
	}
	return res
}

func (kc *kchan) Delete(key string){
	delete(kc.kc, key)
}

func InitKChan(){
	 DeviceChan = &kchan{
		make(map[string]chan string),
	}
	ClientChan = &kchan{
		make(map[string]chan string),
	}
}