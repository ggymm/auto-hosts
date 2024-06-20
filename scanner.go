package main

import (
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

type Scanner struct {
	dev string // 设备 Id

	// 本地与网关的设备信息
	loIp, gwIp net.IP
	loHw, gwHw net.HardwareAddr

	// 超时时间
	timeout time.Duration
}

func NewScanner() *Scanner {
	return &Scanner{
		timeout: 5 * time.Second,
	}
}

func (s *Scanner) Init(d *Device) (err error) {
	s.dev = d.Name
	s.loIp = d.IpAddr
	s.loHw = d.HwAddr

	// 获取网关 mac 地址
	s.gwIp, err = ParseGateway(s.loIp)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Err(errors.WithStack(err)).Msg("parse gateway ip error")
		return
	}

	en := &layers.Ethernet{
		SrcMAC:       s.loHw,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		Operation:         layers.ARPRequest,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(s.gwIp),
		SourceHwAddress:   []byte(s.loHw),
		SourceProtAddress: []byte(s.loIp),
	}

	handle, err := pcap.OpenLive(s.dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Err(errors.WithStack(err)).Msg("open device error")
		return
	}
	buf := gopacket.NewSerializeBuffer()
	ops := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err = gopacket.SerializeLayers(buf, ops, en, arp)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Err(errors.WithStack(err)).Msg("build packet error")
		return err
	}
	err = handle.WritePacketData(buf.Bytes())
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Err(errors.WithStack(err)).Msg("write packet error")
	}
	defer handle.Close()
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			arp = layer.(*layers.ARP)
			if arp.Operation == layers.ARPReply && net.IP(arp.SourceProtAddress).Equal(s.gwIp) {
				s.gwHw = arp.SourceHwAddress
				return
			}
		}
	}
	return
}

func (s *Scanner) Start(domains []string, nameservers []string) (ips map[string][]string) {
	ips = make(map[string][]string)
	for _, domain := range domains {
		ips[domain] = make([]string, 0)
	}

	// 创建设备
	read, err := pcap.OpenLive(s.dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Err(errors.WithStack(err)).Msg("open read device error")
		return
	}
	defer read.Close()
	write, err := pcap.OpenLive(s.dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Err(errors.WithStack(err)).Msg("open write device error")
		return
	}
	defer write.Close()

	go func() {
		// 设置过滤条件
		err = read.SetBPFFilter("udp and port 53")
		if err != nil {
			log.Error().
				Str("1.dev", s.dev).
				Err(errors.WithStack(err)).Msg("handle set bpf filter error")
			return
		}

		// 读取 dns 响应包
		source := gopacket.NewPacketSource(read, read.LinkType())
		for packet := range source.Packets() {
			layer := packet.Layer(layers.LayerTypeDNS)
			if layer == nil {
				continue
			}
			dns, _ := layer.(*layers.DNS)
			for _, q := range dns.Questions {
				if ip, ok := ips[string(q.Name)]; ok {
					for _, a := range dns.Answers {
						if a.IP.To4() == nil || a.IP.IsLoopback() {
							continue
						}

						ip = append(ip, a.IP.To4().String())
					}
				}
			}
		}
	}()

	for i, ns := range nameservers {
		if i%100 == 0 {
			time.Sleep(100 * time.Millisecond)
		}

		en := &layers.Ethernet{
			SrcMAC:       s.loHw,
			DstMAC:       s.gwHw,
			EthernetType: layers.EthernetTypeIPv4,
		}
		ip := &layers.IPv4{
			Version:  4,
			TTL:      64,
			SrcIP:    s.loIp.To4(),
			DstIP:    net.ParseIP(ns).To4(),
			Protocol: layers.IPProtocolUDP,
		}
		udp := &layers.UDP{
			SrcPort: layers.UDPPort(54321),
			DstPort: layers.UDPPort(53),
		}
		_ = udp.SetNetworkLayerForChecksum(ip)
		dns := &layers.DNS{
			QR:           false,
			OpCode:       layers.DNSOpCodeQuery,
			AA:           false,
			TC:           false,
			RD:           true,
			RA:           false,
			Z:            2,
			ResponseCode: layers.DNSResponseCodeNoErr,
			QDCount:      1,
			ANCount:      0,
			NSCount:      0,
			ARCount:      0,
			Questions:    make([]layers.DNSQuestion, 0),
		}

		// 一次性查询多个域名
		for _, domain := range domains {
			dns.Questions = append(dns.Questions, layers.DNSQuestion{
				Name:  []byte(domain),
				Type:  layers.DNSTypeA,
				Class: layers.DNSClassIN,
			})
		}

		buf := gopacket.NewSerializeBuffer()
		ops := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		err = gopacket.SerializeLayers(buf, ops, en, ip, udp, dns)
		if err != nil {
			log.Error().
				Str("1.dev", s.dev).
				Str("2.loIp", s.loIp.String()).
				Str("3.loHw", s.loHw.String()).
				Err(errors.WithStack(err)).Msg("serialize packet error")
			continue
		}
		err = write.WritePacketData(buf.Bytes())
		if err != nil {
			log.Error().
				Str("1.dev", s.dev).
				Str("2.loIp", s.loIp.String()).
				Str("3.loHw", s.loHw.String()).
				Str("4.gwIp", s.gwIp.String()).
				Err(errors.WithStack(err)).Msg("send packet error")
			continue
		}
	}

	time.Sleep(s.timeout) // 收集相应包
	return ips
}
