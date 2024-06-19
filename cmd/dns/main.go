package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"auto-hosts/log"
)

var (
	nameserversFile     = "nameservers.txt"
	nameserversTempFile = "nameserversTemp.txt"
)

func init() {
	dir := ""
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Base(exe)
	if !strings.HasPrefix(exe, os.TempDir()) && !strings.HasPrefix(path, "___") {
		dir = filepath.Dir(exe)
	} else {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			// 需要根据当前文件所处目录，修改相对位置
			dir = filepath.Join(filepath.Dir(filename), "../../")
		}
	}
	wd := filepath.Join(dir, "temp")

	// 设置 app 工作目录
	err = os.Chdir(wd)
	if err != nil {
		panic(errors.WithStack(err))
	}

	log.Init()
}

func main() {
	// 删除 nameservers.txt 文件
	_ = os.Remove(nameserversFile)

	// 获取最新的 nameservers 并且按照响应时间排序
	nss := fetchNameservers()
	dst := filterNameservers(nss)

	// 保存 nameservers.txt 文件
	f1, err := os.OpenFile(nameserversFile, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("open nameservers file error")
		return
	}
	for i, s := range dst {
		_, _ = f1.WriteString(s + "\n")
		log.Info().Msgf("%d: %s", i, s)
	}
}

func fetchNameservers() (nss []string) {
	nss = make([]string, 0)
	url := "https://public-dns.info/nameservers.txt"
	_, err := resty.New().R().SetOutput(nameserversFile).Get(url)
	if err != nil {
		log.Error().
			Str("url", url).
			Err(errors.WithStack(err)).Msg("download nameservers error")
		return
	}

	var (
		f1 *os.File
		f2 *os.File
	)
	f1, err = os.OpenFile(nameserversFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversFile).
			Err(errors.WithStack(err)).Msg("read nameservers file error")
		return
	}
	f2, err = os.OpenFile(nameserversTempFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", nameserversTempFile).
			Err(errors.WithStack(err)).Msg("create nameservers file error")
	}
	fb := bufio.NewReader(f1)
	for {
		l, _, err1 := fb.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		ip := net.ParseIP(string(l))
		if ip.To4() != nil {
			_, _ = f2.Write(l)
			_, _ = f2.WriteString("\n")

			nss = append(nss, string(l))
		}
	}
	_ = f1.Close()
	_ = f2.Close()

	// 重命名文件
	err = os.Rename(nameserversTempFile, nameserversFile)
	if err != nil {
		log.Error().
			Str("file1", nameserversFile).
			Str("file2", nameserversTempFile).
			Err(errors.WithStack(err)).Msg("rename nameservers file error")
		return
	}
	return
}

func filterNameservers(src []string) (dst []string) {
	// 测试 nameservers 是否有效
	dev := `\Device\NPF_{81A86FFA-2C4F-4E6B-AD4E-29036647FB75}`
	loIp := net.IP{192, 168, 1, 27}
	loHw := net.HardwareAddr{0x54, 0x05, 0xdb, 0x83, 0x7f, 0xa5}
	//gwIp := net.IP{192, 168, 1, 1}
	gwHw := net.HardwareAddr{0xa0, 0x04, 0x60, 0x92, 0xa6, 0x02}

	// 创建设备
	send, err := pcap.OpenLive(dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", dev).
			Err(errors.WithStack(err)).Msg("open send device error")
		return
	}
	defer send.Close()
	recv, err := pcap.OpenLive(dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", dev).
			Err(errors.WithStack(err)).Msg("open recv device error")
		return
	}
	defer recv.Close()

	go func() {
		err = recv.SetBPFFilter("dst host " + loIp.To4().String() + " and icmp")
		if err != nil {
			log.Error().
				Str("1.dev", dev).
				Err(errors.WithStack(err)).Msg("handle set bpf filter error")
			return
		}

		// 读取 icmp 包
		source := gopacket.NewPacketSource(recv, recv.LinkType())
		for packet := range source.Packets() {
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer == nil {
				return
			}

			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			if icmpLayer == nil {
				return
			}
			icmp := icmpLayer.(*layers.ICMPv4)
			if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoReply {
				ip := ipLayer.(*layers.IPv4)
				ns := ip.SrcIP.String()

				log.Info().Msgf("recv packet from %s", ns)
				dst = append(dst, ns)
			}
		}
	}()

	id := uint16(os.Getpid())
	seq := uint16(0)
	buf := gopacket.NewSerializeBuffer()
	ops := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	for i, ns := range src {
		if i%100 == 0 {
			time.Sleep(100 * time.Millisecond)
		}

		en := &layers.Ethernet{
			SrcMAC:       loHw,
			DstMAC:       gwHw,
			EthernetType: layers.EthernetTypeIPv4,
		}
		ip := &layers.IPv4{
			Version:  4,
			TTL:      64,
			SrcIP:    loIp.To4(),
			DstIP:    net.ParseIP(ns).To4(),
			Protocol: layers.IPProtocolICMPv4,
		}
		icmp := &layers.ICMPv4{
			Id:       id,
			Seq:      seq,
			TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		}

		err1 := gopacket.SerializeLayers(buf, ops, en, ip, icmp)
		if err1 != nil {
			log.Error().
				Str("1.dev", dev).
				Str("2.loIp", loIp.String()).
				Str("3.loHw", loHw.String()).
				Err(errors.WithStack(err1)).Msg("serialize packet error")
			continue
		}
		err2 := send.WritePacketData(buf.Bytes())
		if err2 != nil {
			log.Error().
				Str("1.dev", dev).
				Str("2.loIp", loIp.String()).
				Str("3.loHw", loHw.String()).
				Err(errors.WithStack(err2)).Msg("write packet error")
			continue
		}

		seq++
		log.Info().Msgf("send packet to %s", ns)
	}

	time.Sleep(3 * time.Second) // 收集相应包
	return
}
