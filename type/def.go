package def

type Ip struct {
	inner string
}

func IpAddr(ip string) Ip {
	return Ip{ip}
}

type IpHeader struct {
	Src     Ip
	Dst     Ip
	Payload string
}

type Mac struct {
	mac string
}

func Addr(mac string) Mac {
	return Mac{mac}
}

type MacHeader struct {
	Src      Mac
	Dst      Mac
	IpHeader IpHeader
}
