package lvxml

import (
	"encoding/xml"
)

type Memory struct {
	Text string `xml:",chardata"`
	Unit string `xml:"unit,attr,omitempty"`
}

type Vcpu struct {
	Text      string `xml:",chardata"`
	Placement string `xml:"placement,attr"`
}

type OsType struct {
	Text    string `xml:",chardata"`
	Arch    string `xml:"arch,attr"`
	Machine string `xml:"machine,attr"`
}

type Boot struct {
	Dev string `xml:"dev,attr"`
}

type Os struct {
	Type OsType `xml:"type"`
	Boot []Boot `xml:"boot"`
}

type LibosinfoOs struct {
	ID string `xml:"id,attr"`
}
type Libosinfo struct {
	Libosinfo string      `xml:"xmlns:libosinfo,attr"`
	Os        LibosinfoOs `xml:"libosinfo:os"`
}

type Metadata struct {
	Libosinfo Libosinfo `xml:"libosinfo:libosinfo"`
}

type Features struct {
	Acpi string `xml:"acpi"`
	Apic string `xml:"apic"`
}

type Cpu struct {
	Text  string `xml:",chardata"`
	Mode  string `xml:"mode,attr"`
	Check string `xml:"check,attr"`
}

type Timer struct {
	Name       string `xml:"name,attr"`
	Tickpolicy string `xml:"tickpolicy,attr,omitempty"`
	Present    string `xml:"present,attr,omitempty"`
}

type Clock struct {
	Offset string  `xml:"offset,attr"`
	Timer  []Timer `xml:"timer"`
}

type SuspendToMem struct {
	Enabled string `xml:"enabled,attr"`
}

type SuspendToDisk struct {
	Enabled string `xml:"enabled,attr"`
}

type Pm struct {
	SuspendToMem  SuspendToMem  `xml:"suspend-to-mem"`
	SuspendToDisk SuspendToDisk `xml:"suspend-to-disk"`
}

type DiskDriver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type DiskSource struct {
	File string `xml:"file,attr"`
}

type DiskTarget struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type Disk struct {
	Type     string     `xml:"type,attr"`
	Device   string     `xml:"device,attr"`
	Driver   DiskDriver `xml:"driver"`
	Source   DiskSource `xml:"source"`
	Target   DiskTarget `xml:"target"`
	Readonly string     `xml:"readonly,omitempty"`
}

type Controller struct {
	Text  string `xml:",chardata"`
	Type  string `xml:"type,attr"`
	Index string `xml:"index,attr,omitempty"`
	Model string `xml:"model,attr,omitempty"`
	Ports string `xml:"ports,attr,omitempty"`
}

type Mac struct {
	Address string `xml:"address,attr"`
}

type InterfaceSource struct {
	Network string `xml:"network,attr"`
}

type InterfaceModel struct {
	Type string `xml:"type,attr"`
}

type Interface struct {
	Type   string          `xml:"type,attr"`
	Mac    Mac             `xml:"mac"`
	Source InterfaceSource `xml:"source"`
	Model  InterfaceModel  `xml:"model"`
}

type SerialModel struct {
	Name string `xml:"name,attr"`
}

type SerialTarget struct {
	Type  string      `xml:"type,attr"`
	Port  string      `xml:"port,attr,omitempty"`
	Model SerialModel `xml:"model"`
}
type Serial struct {
	Type   string       `xml:"type,attr"`
	Target SerialTarget `xml:"target"`
}

type ConsoleTarget struct {
	Type string `xml:"type,attr"`
	Port string `xml:"port,attr,omitempty"`
}

type Console struct {
	Type   string        `xml:"type,attr"`
	Target ConsoleTarget `xml:"target"`
}

type ChannelTarget struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr,omitempty"`
}

type Channel struct {
	Type   string        `xml:"type,attr"`
	Target ChannelTarget `xml:"target"`
}

type Input struct {
	Type string `xml:"type,attr"`
	Bus  string `xml:"bus,attr"`
}

type GraphicsListen struct {
	Type    string `xml:"type,attr"`
	Address string `xml:"address,attr"`
}

type Graphics struct {
	Type       string         `xml:"type,attr"`
	Port       string         `xml:"port,attr"`
	Autoport   string         `xml:"autoport,attr"`
	ListenAttr string         `xml:"listen,attr"`
	Listen     GraphicsListen `xml:"listen"`
}

type Audio struct {
	ID   string `xml:"id,attr"`
	Type string `xml:"type,attr"`
}

type VideoModel struct {
	Type    string `xml:"type,attr"`
	Vram    string `xml:"vram,attr,omitempty"`
	Heads   string `xml:"heads,attr"`
	Primary string `xml:"primary,attr"`
}

type Video struct {
	Model VideoModel `xml:"model"`
}

type Memballoon struct {
	Model string `xml:"model,attr"`
}

type RngBackend struct {
	Text  string `xml:",chardata"`
	Model string `xml:"model,attr"`
}

type Rng struct {
	Model   string     `xml:"model,attr"`
	Backend RngBackend `xml:"backend"`
}

type Devices struct {
	Text       string       `xml:",chardata"`
	Emulator   string       `xml:"emulator"`
	Disk       []Disk       `xml:"disk"`
	Controller []Controller `xml:"controller"`
	Interface  []Interface  `xml:"interface"`
	Serial     Serial       `xml:"serial"`
	Console    Console      `xml:"console"`
	Channel    Channel      `xml:"channel,omitempty"`
	Input      []Input      `xml:"input"`
	Graphics   Graphics     `xml:"graphics"`
	Audio      Audio        `xml:"audio"`
	Video      Video        `xml:"video"`
	Memballoon Memballoon   `xml:"memballoon"`
	Rng        Rng          `xml:"rng"`
}

type Domain struct {
	XMLName       xml.Name `xml:"domain"`
	Text          string   `xml:",chardata"`
	Type          string   `xml:"type,attr"`
	Name          string   `xml:"name"`
	Description   string   `xml:"description"`
	Metadata      Metadata `xml:"metadata"`
	Memory        Memory   `xml:"memory"`
	CurrentMemory Memory   `xml:"currentMemory,omitempty"`
	//	MaxMemory     Memory   `xml:"maxMemory,omitempty"`
	Vcpu       Vcpu     `xml:"vcpu"`
	Os         Os       `xml:"os"`
	Features   Features `xml:"features"`
	Cpu        Cpu      `xml:"cpu"`
	Clock      Clock    `xml:"clock"`
	OnPoweroff string   `xml:"on_poweroff"`
	OnReboot   string   `xml:"on_reboot"`
	OnCrash    string   `xml:"on_crash"`
	Pm         Pm       `xml:"pm"`
	Devices    Devices  `xml:"devices"`
}

//getDefaultDomainXML returns a Domain struct with all necessary info to define a domain
func GetDefaultDomainXML(dom *Domain) {
	dom.Type = "kvm"
	dom.Name = "Default-VM"
	dom.Description = "Created with lvapi"
	dom.Metadata = Metadata{
		Libosinfo{
			Libosinfo: "http://libosinfo.org/xmlns/libvirt/domain/1.0",
			Os:        LibosinfoOs{ID: "http://archlinux.org/archlinux/rolling/"},
		},
	}
	dom.Memory = Memory{Unit: "MiB", Text: "2048"}
	dom.CurrentMemory = Memory{Unit: "MiB", Text: "2048"}
	dom.Vcpu = Vcpu{Placement: "static", Text: "2"}
	dom.Os = Os{
		Type: OsType{Arch: "x86_64", Machine: "pc-q35-6.2", Text: "hvm"},
		Boot: []Boot{
			{Dev: "cdrom"},
			{Dev: "hd"},
		},
	}
	dom.Cpu = Cpu{
		Mode:  "host-model",
		Check: "partial",
	}
	dom.Clock = Clock{
		Offset: "utc",
		Timer: []Timer{
			{Name: "rtc", Tickpolicy: "catchup"},
			{Name: "pit", Tickpolicy: "delay"},
			{Name: "hpet", Present: "no"},
		},
	}
	dom.OnPoweroff = "destroy"
	dom.OnReboot = "restart"
	dom.OnCrash = "destroy"
	dom.Pm = Pm{
		SuspendToMem:  SuspendToMem{Enabled: "no"},
		SuspendToDisk: SuspendToDisk{Enabled: "no"},
	}

	dom.Devices = Devices{
		Emulator: "/usr/bin/qemu-system-x86_64",
		Disk: []Disk{
			{Type: "file", Device: "disk", Driver: DiskDriver{Name: "qemu", Type: "qcow2"}, Source: DiskSource{File: "/virt/images/_frr-base.qcow2"}, Target: DiskTarget{Dev: "vda", Bus: "virtio"}},
			{Type: "file", Device: "cdrom", Driver: DiskDriver{Name: "qemu", Type: "raw"}, Source: DiskSource{File: "/virt/isos/arch.iso"}, Target: DiskTarget{Dev: "sda", Bus: "sata"}, Readonly: " "},
		},
		Controller: []Controller{
			{Type: "usb", Model: "qemu-xhci", Ports: "15"},
			{Type: "sata"},
			{Type: "pci", Model: "pcie-root"},
			{Type: "pci", Model: "pcie-root-port"},
			{Type: "virtio-serial"},
		},
		Interface: []Interface{
			{Type: "network", Mac: Mac{Address: "52:54:00:00:00:00"}, Source: InterfaceSource{Network: "VID100"}, Model: InterfaceModel{Type: "virtio"}},
		},
		Serial:  Serial{Type: "pty", Target: SerialTarget{Type: "isa-serial", Model: SerialModel{Name: "isa-serial"}}},
		Console: Console{Type: "pty", Target: ConsoleTarget{Type: "serial"}},
		Channel: Channel{Type: "unix", Target: ChannelTarget{Type: "virtio"}},
		Input: []Input{
			{Type: "tablet", Bus: "usb"},
			{Type: "mouse", Bus: "ps2"},
			{Type: "keyboard", Bus: "ps2"},
		},
		Graphics:   Graphics{Type: "vnc", Port: "-1", Autoport: "yes", ListenAttr: "0.0.0.0", Listen: GraphicsListen{Type: "address", Address: "0.0.0.0"}},
		Audio:      Audio{ID: "1", Type: "none"},
		Video:      Video{Model: VideoModel{Type: "vga", Heads: "1", Primary: "yes"}},
		Memballoon: Memballoon{Model: "virtio"},
		Rng:        Rng{Model: "virtio", Backend: RngBackend{Model: "random", Text: "/dev/urandom"}},
	}
}
