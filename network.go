package systemd

import (
	"fmt"
	"junheiot/util"
	"os/exec"

	"gopkg.in/ini.v1"
)

type SystemdNetwork struct {
	Match struct {
		Name string
	}
	Network struct {
		DHCP    string //此选项貌似喝下边不能同时
		Address []string
		Gateway string
		DNS     []string
	}
	DHCP struct {
		UseDNS    string
		UseRoutes string
	}
}

var confPathFmt = "/etc/systemd/network/20-wired-%s.network"

func confPath(name string) string {
	return fmt.Sprintf(confPathFmt, name)
}

func (n *SystemdNetwork) Load(name string) error {
	f, err := ini.ShadowLoad(confPath(name))
	if err != nil {
		return err
	}
	err = f.MapTo(n)
	if err != nil {
		return err
	}
	n.Network.Address = f.Section("Network").Key("Address").ValueWithShadows()
	n.Network.DNS = f.Section("Network").Key("DNS").ValueWithShadows()
	return nil
}

func (n *SystemdNetwork) SaveTo(name string) error {
	file := ini.Empty(ini.LoadOptions{AllowShadows: true})

	// Match.Name
	file.Section("Match").Key("Name").SetValue(n.Match.Name)
	// Network
	if n.Network.DHCP != "" {
		// DHCP
		file.Section("Network").Key("DHCP").SetValue(n.Network.DHCP)
		file.Section("DHCP").Key("UseDNS").SetValue(n.DHCP.UseDNS)
		file.Section("DHCP").Key("UseRoutes").SetValue(n.DHCP.UseRoutes)
	} else {
		// static
		// Address
		file.Section("Network").DeleteKey("Address")
		for _, addr := range n.Network.Address {
			file.Section("Network").Key("Address").AddShadow(addr)
		}
		// Gateway
		file.Section("Network").Key("Gateway").SetValue(n.Network.Gateway)
		// DNS
		for _, addr := range n.Network.DNS {
			file.Section("Network").Key("DNS").AddShadow(addr)
		}
	}
	return file.SaveTo(confPath(name))
}

// 删除这个文件
func (n *SystemdNetwork) RemoveFile(name string) error {
	return util.RemoveFile(confPath(name))
}

// func StopService() error {
// 	return exec.Command("systemctl", "stop", "systemd-networkd.service").Run()
// }

// func StartService() error {
// 	return exec.Command("systemctl", "start", "systemd-networkd.service").Run()
// }

func ReStartService() error {
	return exec.Command("systemctl", "restart", "systemd-networkd.service").Run()
}

// func EnableService() error {
// 	return exec.Command("systemctl", "enable", "systemd-networkd.service").Run()
// }
