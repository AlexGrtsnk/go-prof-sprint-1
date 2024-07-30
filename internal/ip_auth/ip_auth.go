package subnetchecker

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func ResolveIP(r *http.Request) (net.IP, error) {

	// смотрим заголовок запроса X-Real-IP
	ipStr := r.Header.Get("X-Real-IP")
	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
		// этот заголовок содержит адреса отправителя и промежуточных прокси
		// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
		ips := r.Header.Get("X-Forwarded-For")
		// разделяем цепочку адресов
		ipStrs := strings.Split(ips, ",")
		// интересует только первый
		ipStr = ipStrs[0]
		// парсим
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}
	return ip, nil
	//}
}

func CheckUserIPinSubnet(userIP net.IP, subnetIP string) (flag bool, err error) {
	_, ipnetA, err := net.ParseCIDR(subnetIP)
	if err != nil{
		return flag, err
	}
	if ipnetA.Contains(userIP) {
		return true, nil
	} else {
		return false, nil
	}

}
