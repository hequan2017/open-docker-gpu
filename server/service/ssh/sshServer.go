package ssh

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ssh"
	sshReq "github.com/flipped-aurora/gin-vue-admin/server/model/ssh/request"
	sshc "golang.org/x/crypto/ssh"
)

type SshServerService struct{}

// CreateSshServer 创建Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) CreateSshServer(ctx context.Context, record *ssh.SshServer) (err error) {
	if record.Port == nil {
		var p int64 = 22
		record.Port = &p
	}
	if record.Username == nil || *record.Username == "" {
		u := "root"
		record.Username = &u
	}
	if record.Ip != nil && record.Password != nil && record.Port != nil && record.Username != nil {
		addr := fmt.Sprintf("%s:%d", *record.Ip, int(*record.Port))
		cfg := &sshc.ClientConfig{User: *record.Username, Auth: []sshc.AuthMethod{sshc.Password(*record.Password)}, HostKeyCallback: sshc.InsecureIgnoreHostKey(), Timeout: 5 * time.Second}
		if c, e := sshc.Dial("tcp", addr, cfg); e == nil {
			s := "正常"
			record.Status = &s
			_ = c.Close()
		} else {
			s := "异常"
			record.Status = &s
		}
	}
	err = global.GVA_DB.Create(record).Error
	return err
}

// DeleteSshServer 删除Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) DeleteSshServer(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&ssh.SshServer{}, "id = ?", ID).Error
	return err
}

// DeleteSshServerByIds 批量删除Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) DeleteSshServerByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]ssh.SshServer{}, "id in ?", IDs).Error
	return err
}

// UpdateSshServer 更新Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) UpdateSshServer(ctx context.Context, record *ssh.SshServer) (err error) {
	if record.Port == nil {
		var p int64 = 22
		record.Port = &p
	}
	if record.Username == nil || *record.Username == "" {
		u := "root"
		record.Username = &u
	}
	if record.Ip != nil && record.Password != nil && record.Port != nil && record.Username != nil {
		addr := fmt.Sprintf("%s:%d", *record.Ip, int(*record.Port))
		cfg := &sshc.ClientConfig{User: *record.Username, Auth: []sshc.AuthMethod{sshc.Password(*record.Password)}, HostKeyCallback: sshc.InsecureIgnoreHostKey(), Timeout: 5 * time.Second}
		if c, e := sshc.Dial("tcp", addr, cfg); e == nil {
			s := "正常"
			record.Status = &s
			_ = c.Close()
		} else {
			s := "异常"
			record.Status = &s
		}
	}
	err = global.GVA_DB.Model(&ssh.SshServer{}).Where("id = ?", record.ID).Updates(record).Error
	return err
}

// GetSshServer 根据ID获取Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) GetSshServer(ctx context.Context, ID string) (record ssh.SshServer, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&record).Error
	return record, err
}

// GetSshServerInfoList 分页获取Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) GetSshServerInfoList(ctx context.Context, info sshReq.SshServerSearch) (list []ssh.SshServer, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&ssh.SshServer{})
	var sshs []ssh.SshServer
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Ip != nil && *info.Ip != "" {
		db = db.Where("ip LIKE ?", "%"+*info.Ip+"%")
	}
	if info.Port != nil {
		db = db.Where("port = ?", *info.Port)
	}
	if info.Username != nil && *info.Username != "" {
		db = db.Where("username LIKE ?", "%"+*info.Username+"%")
	}
	if info.Password != nil && *info.Password != "" {
		db = db.Where("password = ?", *info.Password)
	}
	if info.Status != nil && *info.Status != "" {
		db = db.Where("status = ?", *info.Status)
	}
	if info.Label != nil && *info.Label != "" {
		db = db.Where("label LIKE ?", "%"+*info.Label+"%")
	}
	if info.Region != nil && *info.Region != "" {
		db = db.Where("region LIKE ?", "%"+*info.Region+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["ip"] = true
	orderMap["status"] = true
	orderMap["label"] = true
	orderMap["region"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&sshs).Error
	return sshs, total, err
}
func (sshService *SshServerService) GetSshServerPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

type GPUInfo struct {
	Index       int    `json:"index"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	MemoryTotal int    `json:"memoryTotal"`
	MemoryUsed  int    `json:"memoryUsed"`
	Utilization int    `json:"utilization"`
	Driver      string `json:"driver"`
}

func (sshService *SshServerService) runSSH(ID string, cmd string) (string, error) {
	var rec ssh.SshServer
	if err := global.GVA_DB.Where("id = ?", ID).First(&rec).Error; err != nil {
		return "", err
	}
	if rec.Ip == nil || *rec.Ip == "" || rec.Username == nil || *rec.Username == "" || rec.Password == nil {
		return "", fmt.Errorf("ssh记录信息不完整")
	}
	port := 22
	if rec.Port != nil {
		port = int(*rec.Port)
	}
	addr := fmt.Sprintf("%s:%d", *rec.Ip, port)
	cfg := &sshc.ClientConfig{User: *rec.Username, Auth: []sshc.AuthMethod{sshc.Password(*rec.Password)}, HostKeyCallback: sshc.InsecureIgnoreHostKey(), Timeout: 10 * time.Second}
	// sshc.ClientConfig has Timeout from go1.20+, ensure available; fallback if older
	c, err := sshc.Dial("tcp", addr, cfg)
	if err != nil {
		// try with manual dial timeout if Dial hangs
		d := &net.Dialer{Timeout: 8 * time.Second}
		conn, e2 := d.Dial("tcp", addr)
		if e2 != nil {
			return "", err
		}
		cc, chans, reqs, e3 := sshc.NewClientConn(conn, addr, cfg)
		if e3 != nil {
			return "", err
		}
		c = sshc.NewClient(cc, chans, reqs)
	}
	defer c.Close()
	sess, e := c.NewSession()
	if e != nil {
		return "", e
	}
	defer sess.Close()
	var buf bytes.Buffer
	sess.Stdout = &buf
	sess.Stderr = &buf
	if e = sess.Run(cmd); e != nil {
		return buf.String(), e
	}
	return buf.String(), nil
}

// runSSHShell 使用远端shell执行复杂命令（管道/重定向）
func (sshService *SshServerService) runSSHShell(ID string, cmd string) (string, error) {
	// 优先 bash -lc，失败回退 sh -lc
	// 注意: 进行基本的单引号转义
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	// 尝试 bash
	out, err := sshService.runSSH(ID, "bash -lc "+quote(cmd))
	if err == nil && strings.TrimSpace(out) != "" {
		return out, nil
	}
	// 回退 sh
	out2, err2 := sshService.runSSH(ID, "sh -lc "+quote(cmd))
	if err2 != nil {
		return out2, err2
	}
	return out2, nil
}

// FetchNvidiaSmiText 获取 nvidia-smi 输出尾部摘要文本
func (sshService *SshServerService) FetchNvidiaSmiText(ctx context.Context, ID string, tail int) (string, error) {
	if tail <= 0 {
		tail = 50
	}
	if tail > 2000 {
		tail = 2000
	}
	cmd := fmt.Sprintf("command -v nvidia-smi >/dev/null 2>&1 && nvidia-smi | tail -n %d || echo 'nvidia-smi not found'", tail)
	return sshService.runSSHShell(ID, cmd)
}

// FetchGPUInfo 通过SSH执行命令获取GPU信息
func (sshService *SshServerService) FetchGPUInfo(ctx context.Context, ID string) (items []GPUInfo, err error) {
	// 首选 nvidia-smi
	out, e := sshService.runSSH(ID, "nvidia-smi --query-gpu=index,uuid,name,memory.total,memory.used,utilization.gpu,driver_version --format=csv,noheader,nounits")
	if e != nil || strings.TrimSpace(out) == "" {
		// ROCm 尝试
		rocmOut, e2 := sshService.runSSH(ID, "rocm-smi --showproductname --showuse --showmeminfo vram --showuniqueid --csv")
		if e2 == nil && strings.TrimSpace(rocmOut) != "" {
			lines := strings.Split(strings.TrimSpace(rocmOut), "\n")
			for i, ln := range lines {
				parts := strings.Split(ln, ",")
				if len(parts) < 4 {
					continue
				}
				name := strings.TrimSpace(parts[0])
				util := parseInt(strings.TrimSpace(parts[1]))
				memTotal := parseInt(strings.TrimSpace(parts[2]))
				memUsed := parseInt(strings.TrimSpace(parts[3]))
				uuid := ""
				if len(parts) > 4 {
					uuid = strings.TrimSpace(parts[4])
				}
				items = append(items, GPUInfo{Index: i, UUID: uuid, Name: name, MemoryTotal: memTotal, MemoryUsed: memUsed, Utilization: util, Driver: "ROCm"})
			}
			return items, nil
		}
		// 最后回退 lspci 简要信息
		lspci, e3 := sshService.runSSH(ID, "lspci | grep -i 'vga\\|3d' | grep -i 'nvidia\\|amd'")
		if e3 == nil && strings.TrimSpace(lspci) != "" {
			lines := strings.Split(strings.TrimSpace(lspci), "\n")
			for i, ln := range lines {
				name := strings.TrimSpace(ln)
				items = append(items, GPUInfo{Index: i, Name: name})
			}
			return items, nil
		}
		return nil, fmt.Errorf("无法获取GPU信息")
	}
	// 解析 nvidia-smi 输出
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.Split(ln, ",")
		if len(parts) < 6 {
			continue
		}
		idx := parseInt(strings.TrimSpace(parts[0]))
		uuid := strings.TrimSpace(parts[1])
		name := strings.TrimSpace(parts[2])
		memTotal := parseInt(strings.TrimSpace(parts[3]))
		memUsed := parseInt(strings.TrimSpace(parts[4]))
		util := parseInt(strings.TrimSpace(parts[5]))
		driver := ""
		if len(parts) > 6 {
			driver = strings.TrimSpace(parts[6])
		} else {
			driver = ""
		}
		items = append(items, GPUInfo{Index: idx, UUID: uuid, Name: name, MemoryTotal: memTotal, MemoryUsed: memUsed, Utilization: util, Driver: driver})
	}
	return items, nil
}

func parseInt(s string) int {
	s = strings.TrimSpace(strings.TrimSuffix(s, "%"))
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}
