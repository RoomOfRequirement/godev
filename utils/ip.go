package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strings"
)

// IsPrivateIPv4 ...
//	IPv4 has private ip range but v6 not
func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 ||
			ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) ||
			ip[0] == 192 && ip[1] == 168)
}

// GetPublicIP ...
func GetPublicIP() (net.IP, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(strings.TrimSpace(string(content))), nil
}

// GetLocalIP ...
//	non-loopback
func GetLocalIP(onlyV4 bool) ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	// for certain iface: iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, 0, len(addrs))
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() {
			var ip net.IP
			if onlyV4 {
				ip = ipnet.IP.To4()
			} else {
				ip = ipnet.IP.To16()
			}
			if ip != nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}

// IsIPBetween ...
//	v4 / v6
func IsIPBetween(from, to, in net.IP) (bool, error) {
	if from == nil || to == nil || in == nil {
		return false, fmt.Errorf("invalid input: %v, %v, %v", from, to, in)
	}
	// 16-byte representation
	//	To16 can handle v4 and v6, while To4 can only handle v4
	from16, to16, in16 := from.To16(), to.To16(), in.To16()
	if from16 == nil || to16 == nil || in16 == nil {
		return false, fmt.Errorf("invalid input: %v, %v, %v", from, to, in)
	}
	if bytes.Compare(in16, from16) >= 0 && bytes.Compare(in16, to16) <= 0 {
		return true, nil
	}
	return false, nil
}

// IPv4AtoN ...
//	notice: inet_aton() returns non-zero if the address is a valid one, and it returns zero if the address is invalid
func IPv4AtoN(ip net.IP) int64 {
	/*
		if ip == nil || ip.To4() == nil {
				return 0
		}
		// IPv4 use ".", IPv6 use ":" or with only one "::"
		bits := strings.Split(ip.String(), ".")
		if len(bits) != 4 {
			return 0
		}
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		return int64(b0) << 24 | int64(b1) << 16 | int64(b2) << 8 | int64(b3)
	*/
	ipv4Int := big.NewInt(0)
	ipv4Int.SetBytes(ip.To4())
	return ipv4Int.Int64()
}

// IPv4NtoA ...
func IPv4NtoA(ipv4 int64) net.IP {
	bs := make([]byte, 4)
	bs[0] = byte(ipv4 & 0xFF)
	bs[1] = byte((ipv4 >> 8) & 0xFF)
	bs[2] = byte((ipv4 >> 16) & 0xFF)
	bs[3] = byte((ipv4 >> 24) & 0xFF)
	return net.IPv4(bs[3], bs[2], bs[1], bs[0]).To4()
}

// IPv6AtoN returns hex string
//	need handle v6 representation of v4
//	notice: maybe better to use ip.String() for string representation
func IPv6AtoN(ip net.IP) string {
	ipInt := big.NewInt(0)
	hexStr := ""
	// v4
	if ip.To4() != nil {
		ipInt.SetBytes(ip.To4())
		hexStr = hex.EncodeToString(ipInt.Bytes())
		return hexStr
	}
	ipInt.SetBytes(ip.To16())
	hexStr = hex.EncodeToString(ipInt.Bytes())
	return hexStr
}

// IPv6NtoA ...
func IPv6NtoA(hexStr string) net.IP {
	bs, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil
	}
	l := len(bs)
	// v4 or v6
	if l == 4 || l == 16 {
		return bs
	}
	return nil
}
