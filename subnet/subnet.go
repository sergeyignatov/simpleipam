package subnet

import (
	"fmt"
	"github.com/sergeyignatov/simpleipam/common"
	icfg "github.com/sergeyignatov/simpleipam/config"
	"github.com/tatsushid/go-fastping"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var config *icfg.Config

type Subnets struct {
	subnets map[string]*Subnet
	sync.RWMutex
}

func checkHostAlive(ip string) bool {
	p := fastping.NewPinger()
	p.Network("udp")
	p.AddIPAddr(&net.IPAddr{net.ParseIP(ip), ""})
	resp := make(chan bool)
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		resp <- true
	}
	p.OnIdle = func() {
		resp <- false
	}

	p.MaxRTT = 100 * time.Millisecond
	p.RunLoop()
	t := <-resp
	//fmt.Printf("ping %s\n", ip)
	return t
}
func (ss *Subnets) Load(cfg *icfg.Config) error {
	config = cfg
	for k, v := range cfg.Subnets {
		s := NewSubnet(k, v)
		ss.Add(s)
	}
	return nil
}
func NewSubnets() *Subnets {
	s := Subnets{}
	rand.Seed(time.Now().Unix())
	s.subnets = make(map[string]*Subnet)
	return &s
}

func (s *Subnet) Add(cl *common.Client) {
	s.Lock()
	defer s.Unlock()
	s.inuse_ip[common.IPAddr(cl.Ip)] = cl
	s.inuse_mac[common.HardwareAddr(cl.Mac)] = cl
	s.inuse_fqdn[cl.Hostname] = common.HardwareAddr(cl.Mac)
}

func (s *Subnet) AddSave(cl *common.Client) error {
	f := path.Join(s.datadir, cl.Mac)
	err := icfg.SaveClient(cl, f)
	if err != nil {
		return err
	}
	s.Add(cl)
	return nil
}
func NewSubnet(cidr string, c icfg.Subnet) *Subnet {
	s := Subnet{}
	s.cidr = cidr
	s.start = net.ParseIP(c.Start)
	s.end = net.ParseIP(c.End)
	s.gateway = net.ParseIP(c.Gateway)
	s.inuse_ip = make(map[common.IPAddr]*common.Client)
	s.inuse_mac = make(map[common.HardwareAddr]*common.Client)
	s.inuse_fqdn = make(map[string]common.HardwareAddr)
	s.datadir = path.Join(config.DataDir, strings.Replace(cidr, "/", "_", 1))
	if _, err := os.Stat(s.datadir); err != nil {
		os.Mkdir(s.datadir, 0755)
	}
	files, err := ioutil.ReadDir(s.datadir)
	if err != nil {
		return &s
	}
	for _, f := range files {
		fpath := path.Join(s.datadir, f.Name())
		if cl, err := icfg.LoadClient(fpath); err == nil {
			s.Add(cl)
		}
	}
	return &s
}
func (ss *Subnets) Add(s *Subnet) {
	ss.Lock()
	defer ss.Unlock()
	ss.subnets[s.cidr] = s
}

type Subnet struct {
	cidr       string
	start      net.IP
	end        net.IP
	gateway    net.IP
	datadir    string
	inuse_mac  map[common.HardwareAddr]*common.Client
	inuse_ip   map[common.IPAddr]*common.Client
	inuse_fqdn map[string]common.HardwareAddr

	sync.RWMutex
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
func (s *Subnet) delete(cl *common.Client) error {
	s.Lock()
	defer s.Unlock()
	f := path.Join(s.datadir, cl.Mac)
	err := icfg.DeleteClient(f)
	if err != nil {
		return err
	}
	delete(s.inuse_ip, common.IPAddr(cl.Ip))
	delete(s.inuse_fqdn, cl.Hostname)
	delete(s.inuse_mac, common.HardwareAddr(cl.Mac))
	return nil
}
func (s *Subnets) getSubnet(subnet string) (*Subnet, error) {
	if s, ok := s.subnets[subnet]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("unable find subnet")
}
func (s *Subnet) releaseIP(mac common.HardwareAddr, ip, fqdn string) error {
	if c, ok := s.inuse_mac[mac]; ok {
		if c.Ip == ip {
			return s.delete(c)
		}
	}
	if m, ok := s.inuse_fqdn[fqdn]; ok {
		if c, ok := s.inuse_mac[m]; ok {
			if c.Ip == ip {
				return s.delete(c)
			}
		}
	}
	return fmt.Errorf("unable find leases")
}
func (s *Subnet) setIP(mac common.HardwareAddr, ip, fqdn string) (*common.Response, error) {
	if c, ok := s.inuse_mac[mac]; ok {
		c.Ip = ip
	}
	cl := common.Client{Ip: ip, Mac: string(mac), Hostname: fqdn, CreateTime: time.Now().Unix()}
	if err := s.AddSave(&cl); err != nil {
		return nil, err
	}
	return &common.Response{ip, s.gateway.String(), fqdn, string(mac), s.cidr}, nil

}
func (s *Subnet) getIP(mac common.HardwareAddr, fqdn string) (*common.Response, error) {
	if c, ok := s.inuse_mac[mac]; ok {
		return &common.Response{c.Ip, s.gateway.String(), fqdn, string(mac), s.cidr}, nil
	}
	if m, ok := s.inuse_fqdn[fqdn]; ok {
		if c, ok := s.inuse_mac[m]; ok {
			if !checkHostAlive(c.Ip) {
				return &common.Response{c.Ip, s.gateway.String(), fqdn, string(m), s.cidr}, nil
			}
		}
	}
	_, ipnet, err := net.ParseCIDR(s.cidr)
	if err != nil {
		return nil, err
	}
	for ip := s.start; ipnet.Contains(ip); inc(ip) {
		if _, ok := s.inuse_ip[common.IPAddr(ip.String())]; !ok {
			if ip.String() == s.end.String() {
				return nil, fmt.Errorf("No more ip to use")
			}
			cl := common.Client{Ip: ip.String(), Mac: string(mac), Hostname: fqdn, CreateTime: time.Now().Unix()}
			if err := s.AddSave(&cl); err != nil {
				return nil, err
			}

			return &common.Response{ip.String(), s.gateway.String(), fqdn, string(mac), s.cidr}, nil
		}
	}
	return nil, fmt.Errorf("unable get ip")
}
func (ss *Subnets) List() map[common.IPAddr]*common.Client {
	for _, s := range ss.subnets {
		return s.inuse_ip
	}
	return nil
}
func (ss *Subnets) macinuse(mac common.HardwareAddr) bool {
	for _, s := range ss.subnets {
		if _, ok := s.inuse_mac[mac]; ok {
			return true
		}
	}
	return false
}
func (s *Subnets) generatemac() (string, error) {

	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] |= 2
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	if s.macinuse(common.HardwareAddr(mac)) {
		return s.generatemac()
	}
	return mac, nil
}
func (s *Subnets) ReleaseIP(subnet, mac, ip, fqdn string) error {
	if _, err := net.ParseMAC(mac); err != nil {
		return err
	}
	if _, nt, err := net.ParseCIDR(subnet); err != nil {
		return err
	} else {
		if !nt.Contains(net.ParseIP(ip)) {
			return fmt.Errorf("ip is not belong to subnet")
		}
	}

	sn, err := s.getSubnet(subnet)
	if err != nil {
		return err
	}
	return sn.releaseIP(common.HardwareAddr(strings.ToLower(mac)), ip, fqdn)
	//return nil
}
func (s *Subnets) GetNewIp(subnet, mac, fqdn, oldip string) (*common.Response, error) {
	var macaddr string
	if mac != "" {
		if _, err := net.ParseMAC(mac); err != nil {
			return nil, err
		}
		macaddr = mac
	} else {
		m, err := s.generatemac()
		if err != nil {
			return nil, err
		}
		macaddr = m
	}
	_, nt, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, err
	}
	sn, err := s.getSubnet(subnet)
	if err != nil {
		return nil, err
	}
	if fqdn == "" {
		return nil, fmt.Errorf("empty fqdn")
	}
	if nt.Contains(net.ParseIP(oldip)) {
		return sn.setIP(common.HardwareAddr(strings.ToLower(macaddr)), oldip, fqdn)
	}
	return sn.getIP(common.HardwareAddr(strings.ToLower(macaddr)), fqdn)
}