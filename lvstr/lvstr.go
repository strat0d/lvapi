package lvstr

// converts libvirt objects to more usable
import "libvirt.org/go/libvirt"

// GetBlkioParameters(flags DomainModificationImpact) (*DomainBlkioParameters, error)
// GetBlockInfo(disk string, flags uint32) (*DomainBlockInfo, error)
// GetBlockIoTune(disk string, flags DomainModificationImpact) (*DomainBlockIoTuneParameters, error)
// GetBlockJobInfo(disk string, flags DomainBlockJobInfoFlags) (*DomainBlockJobInfo, error)
// GetCPUStats(startCpu int, nCpus uint, flags uint32) ([]DomainCPUStats, error)
// GetControlInfo(flags uint32) (*DomainControlInfo, error)
// GetDiskErrors(flags uint32) ([]DomainDiskError, error)
// GetEmulatorPinInfo(flags DomainModificationImpact) ([]bool, error)
// GetFSInfo(flags uint32) ([]DomainFSInfo, error)
// GetGuestInfo(types DomainGuestInfoTypes, flags uint32) (*DomainGuestInfo, error)
// GetGuestVcpus(flags uint32) (*DomainGuestVcpus, error)
// GetIOThreadInfo(flags DomainModificationImpact) ([]DomainIOThreadInfo, error)
// GetInfo() (*DomainInfo, error)
// GetInterfaceParameters(device string, flags DomainModificationImpact) (*DomainInterfaceParameters, error)
// GetJobInfo() (*DomainJobInfo, error)
// GetJobStats(flags DomainGetJobStatsFlags) (*DomainJobInfo, error)
// GetLaunchSecurityInfo(flags uint32) (*DomainLaunchSecurityParameters, error)
// GetMemoryParameters(flags DomainModificationImpact) (*DomainMemoryParameters, error)
// GetNumaParameters(flags DomainModificationImpact) (*DomainNumaParameters, error)
// GetPerfEvents(flags DomainModificationImpact) (*DomainPerfEvents, error)
// GetSchedulerParameters() (*DomainSchedulerParameters, error)
// GetSchedulerParametersFlags(flags DomainModificationImpact) (*DomainSchedulerParameters, error)
// GetSecurityLabel() (*SecurityLabel, error)
// GetSecurityLabelList() ([]SecurityLabel, error)
// GetTime(flags uint32) (int64, uint, error)
// GetUUID() ([]byte, error)
// GetVcpuPinInfo(flags DomainModificationImpact) ([][]bool, error)
// GetVcpus() ([]DomainVcpuInfo, error)
// GetVcpusFlags(flags DomainVcpuFlags) (int32, error)

// GetMetadata(metadataType DomainMetadataType, uri string, flags DomainModificationImpact) (string, error)

type Domain struct {
	Autostart   bool
	Hostname    string
	ID          uint
	MaxMemory   uint64
	MaxVcpus    uint
	Messages    []string
	Name        string
	OSType      string
	StateReason int
	State       domainState
	UUID        string
	XMLDesc     string
	VcpuInfo    domainVcpuInfo
}

type domainState struct {
	State int
}

func (s *domainState) String() string {
	switch s.State {
	case int(libvirt.DOMAIN_NOSTATE):
		return "NO STATE"
	case int(libvirt.DOMAIN_RUNNING):
		return "RUNNING"
	case int(libvirt.DOMAIN_BLOCKED):
		return "BLOCKED"
	case int(libvirt.DOMAIN_PAUSED):
		return "PAUSED"
	case int(libvirt.DOMAIN_SHUTDOWN):
		return "SHUTDOWN"
	case int(libvirt.DOMAIN_CRASHED):
		return "CRASHED"
	case int(libvirt.DOMAIN_PMSUSPENDED):
		return "SUSPENDED"
	case int(libvirt.DOMAIN_SHUTOFF):
		return "SHUTOFF"
	default:
		return "UNKNOWN"
	}
}

type domainVcpuInfo struct {
	Number  uint32
	State   vcpuState
	CpuTime uint64
	Cpu     int32
	CpuMap  []bool
}

type vcpuState struct {
	State int
}

func (s *vcpuState) String() string {
	switch s.State {
	case int(libvirt.VCPU_OFFLINE):
		return "OFFLINE"
	case int(libvirt.VCPU_RUNNING):
		return "RUNNING"
	case int(libvirt.VCPU_BLOCKED):
		return "BLOCKED"
	default:
		return "UNKNOWN"
	}
}

func GetDomain(l *libvirt.Domain) *Domain {
	var d = Domain{}

	d.Autostart, _ = l.GetAutostart()
	d.Hostname, _ = l.GetHostname(libvirt.DOMAIN_GET_HOSTNAME_AGENT | libvirt.DOMAIN_GET_HOSTNAME_LEASE)
	d.ID, _ = l.GetID()
	d.MaxMemory, _ = l.GetMaxMemory()
	d.MaxVcpus, _ = l.GetMaxVcpus()
	d.Messages, _ = l.GetMessages(libvirt.DOMAIN_MESSAGE_DEPRECATION | libvirt.DOMAIN_MESSAGE_TAINTING)
	d.Name, _ = l.GetName()
	d.OSType, _ = l.GetOSType()
	d.UUID, _ = l.GetUUIDString()
	s, r, _ := l.GetState()
	d.State = domainState{State: int(s)}
	d.StateReason = r

	return &d
}

/// GetAutostart() (bool, error)
/// GetHostname(flags DomainGetHostnameFlags) (string, error)
// GetID() (uint, error)
// GetMaxMemory() (uint64, error)
// GetMaxVcpus() (uint, error)
// GetMessages(flags DomainMessageType) ([]string, error)

// GetName() (string, error)
// GetOSType() (string, error)
// GetUUIDString() (string, error)
// GetState() (DomainState, int, error)
// GetXMLDesc(flags DomainXMLFlags) (string, error)
