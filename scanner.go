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

	// pcap handle
	handle  *pcap.Handle
	options gopacket.SerializeOptions

	// 超时时间
	timeout time.Duration
}

func NewScanner() *Scanner {
	return &Scanner{
		options: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		timeout: 3 * time.Second,
	}
}

func (s *Scanner) init(d Device) (err error) {
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

	s.handle, err = pcap.OpenLive(s.dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Err(errors.WithStack(err)).Msg("open device error")
		return
	}
	err = s.write(en, arp)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Str("4.gwIp", s.gwIp.String()).
			Err(errors.WithStack(err)).Msg("send packet error")
	}
	//defer s.handle.Close()
	source := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
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

func (s *Scanner) write(l ...gopacket.SerializableLayer) error {
	buf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buf, s.options, l...)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Err(errors.WithStack(err)).Msg("serialize packet error")
		return err
	}
	return s.handle.WritePacketData(buf.Bytes())
}

// DetectIps 探测 ip 按照延时排序
func (s *Scanner) DetectIps(src []string) (dst []string) {
	return src
}

// ResolveIps 解析域名对应的 ip 地址
func (s *Scanner) ResolveIps(nss []string, domain string) (ips []string) {
	ipm := make(map[string]bool)
	go func() {
		handle, err := pcap.OpenLive(s.dev, 65535, true, pcap.BlockForever)
		if err != nil {
			log.Error().
				Str("1.dev", s.dev).
				Err(errors.WithStack(err)).Msg("open device error")
			return
		}
		defer handle.Close()

		// 设置过滤条件
		err = handle.SetBPFFilter("udp and port 53")
		if err != nil {
			log.Error().
				Str("1.dev", s.dev).
				Err(errors.WithStack(err)).Msg("handle set bpf filter error")
			return
		}

		// 读取 dns 响应包
		source := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range source.Packets() {
			layer := packet.Layer(layers.LayerTypeDNS)
			if layer == nil {
				continue
			}
			dns, _ := layer.(*layers.DNS)
			for _, q := range dns.Questions {
				if string(q.Name) != domain {
					continue
				}
				for _, a := range dns.Answers {
					if a.IP.To4() == nil || a.IP.IsLoopback() {
						continue
					}
					ip := a.IP.To4().String()
					if ipm[ip] {
						continue
					} else {
						ipm[ip] = true
					}
				}
			}
		}
	}()

	for i, ns := range nss {
		if i%100 == 0 {
			time.Sleep(60 * time.Millisecond)
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
			Questions: []layers.DNSQuestion{
				{
					Name:  []byte(domain),
					Type:  layers.DNSTypeA,
					Class: layers.DNSClassIN,
				},
			},
		}
		err := s.write(en, ip, udp, dns)
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
	ips = make([]string, 0, len(ipm))
	for ip := range ipm {
		ips = append(ips, ip)
	}
	return ips
}
