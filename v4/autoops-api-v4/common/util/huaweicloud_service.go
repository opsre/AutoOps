package util

import (
    "encoding/json"
    "fmt"

    coreauth "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
    coreRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
    ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
    ecsmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
    ecsregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
    iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
    iammodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
    iamregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
    "strconv"
)

// HuaweiCloudService 提供华为云API调用服务
type HuaweiCloudService struct {
    AccessKey    string
    AccessSecret string
}

// NewHuaweiCloudService 创建新的华为云服务实例
func NewHuaweiCloudService(accessKey, accessSecret string) *HuaweiCloudService {
    return &HuaweiCloudService{
        AccessKey:    accessKey,
        AccessSecret: accessSecret,
    }
}

// HuaweiInstanceInfo 华为云实例信息结构
type HuaweiInstanceInfo struct {
    InstanceID   string   `json:"instance_id"`
    InstanceName string   `json:"instance_name"` // 实例名称
    PrivateIPs   []string `json:"private_ips"`
    PublicIPs    []string `json:"public_ips"`
    CPU          int      `json:"cpu"`         // vCPU 数量
    MemoryGB     int      `json:"memory_gb"`   // 内存(GB)
    SystemDisk   int      `json:"system_disk"` // 系统盘(GB)（若无法获取则为0）
    DataDisks    []int    `json:"data_disks"`  // 数据盘(GB)（若无法获取则为空）
    BillingType  string   `json:"billing_type"`
    OSName       string   `json:"os_name"`
    CreateTime   string   `json:"create_time"`
    ExpireTime   string   `json:"expire_time"`
    Region       string   `json:"region"`
}

// GetInstances 获取华为云实例列表（自动遍历可用区域；失败则用常见区域兜底）
func (s *HuaweiCloudService) GetInstances() ([]HuaweiInstanceInfo, error) {
    regions := s.discoverRegions()

    // 建立 IAM 客户端用于查询项目信息
    auth := coreauth.NewCredentialsBuilder().
        WithAk(s.AccessKey).
        WithSk(s.AccessSecret).
        Build()
    ir := iamregion.ValueOf("cn-north-4")
    if ir == nil {
        ir = coreRegion.NewRegion("cn-north-4", "https://iam.cn-north-4.myhuaweicloud.com")
    }
    ic := iam.NewIamClient(iam.IamClientBuilder().
        WithRegion(ir).
        WithCredential(auth).
        Build())

    // 获取 region -> projectId 映射
    projMap := s.fetchProjectMap(ic)

    var allInstances []HuaweiInstanceInfo

    for _, regionID := range regions {
        projectId, ok := projMap[regionID]
        if !ok || projectId == "" {
            // 无项目ID时跳过该区域
            fmt.Printf("[WARN] 区域 %s 未找到项目ID，跳过\n", regionID)
            continue
        }
        instances, err := s.getInstancesByRegion(regionID, projectId)
        if err != nil {
            // 个别区域失败不影响其他区域
            fmt.Printf("[WARN] 查询华为云区域 %s 失败: %v\n", regionID, err)
            continue
        }
        allInstances = append(allInstances, instances...)
    }

    return allInstances, nil
}

func (s *HuaweiCloudService) getInstancesByRegion(regionID string, projectId string) ([]HuaweiInstanceInfo, error) {
    auth := coreauth.NewCredentialsBuilder().
        WithAk(s.AccessKey).
        WithSk(s.AccessSecret).
        WithProjectId(projectId).
        Build()

    r := ecsregion.ValueOf(regionID)
    // 如果 SDK 未内置该区域，使用通用 Region 构造，避免中断
    if r == nil {
        r = coreRegion.NewRegion(regionID, fmt.Sprintf("https://ecs.%s.myhuaweicloud.com", regionID))
    }

    client := ecs.NewEcsClient(ecs.EcsClientBuilder().
        WithRegion(r).
        WithCredential(auth).
        Build())

    req := &ecsmodel.ListServersDetailsRequest{}
    resp, err := client.ListServersDetails(req)
    if err != nil {
        return nil, fmt.Errorf("区域 %s 列表查询失败: %v", regionID, err)
    }

    if resp == nil || resp.Servers == nil {
        return []HuaweiInstanceInfo{}, nil
    }

    var result []HuaweiInstanceInfo
    for _, sdetail := range *resp.Servers {
        var id, name, osName, created, expired string
        var vcpus, memGB int
        var privIPs, pubIPs []string

        // 基本字段（SDK为非指针类型）
        id = sdetail.Id
        name = sdetail.Name // 获取实例名称
        created = sdetail.Created
        expired = sdetail.Updated // 华为云无显式过期时间，这里留空或复用更新时间

        // OS 信息从元数据中获取（image_name 或 os_type）
        if sdetail.Metadata != nil {
            if imageName, ok := sdetail.Metadata["image_name"]; ok && imageName != "" {
                osName = imageName
            } else if t, ok := sdetail.Metadata["os_type"]; ok && t != "" {
                osName = t
            }
        }

        // 口味信息（vCPU/内存）
        if sdetail.Flavor != nil {
            if n, err := strconv.Atoi(sdetail.Flavor.Vcpus); err == nil {
                vcpus = n
            }
            if mb, err := strconv.Atoi(sdetail.Flavor.Ram); err == nil {
                memGB = mb / 1024
                if memGB == 0 && mb > 0 {
                    memGB = 1
                }
            }
        }

        // IP 地址
        if sdetail.Addresses != nil {
            // Addresses 是 map[string][]ServerAddress
            for _, addrs := range sdetail.Addresses {
                for _, addr := range addrs {
                    ip := addr.Addr
                    if addr.OSEXTIPStype != nil && addr.OSEXTIPStype.Value() == "floating" {
                        pubIPs = append(pubIPs, ip)
                    } else {
                        privIPs = append(privIPs, ip)
                    }
                }
            }
        }

        // 计费方式仅占位（ECS列表通常不直接返回），统一标记为"按需/包年包月未知"
        billing := "未知"

        result = append(result, HuaweiInstanceInfo{
            InstanceID:   id,
            InstanceName: name,
            PrivateIPs:   privIPs,
            PublicIPs:    pubIPs,
            CPU:          vcpus,
            MemoryGB:     memGB,
            SystemDisk:   0,
            DataDisks:    []int{},
            BillingType:  billing,
            OSName:       osName,
            CreateTime:   created,
            ExpireTime:   expired,
            Region:       regionID,
        })
    }

    return result, nil
}

// discoverRegions 自动发现可用区域（优先 SDK 已知区域，失败时提供兜底列表）
func (s *HuaweiCloudService) discoverRegions() []string {
    // 兜底列表（常见区域）
    fallback := []string{
        "cn-north-4",     // 北京4
        "cn-north-1",     // 北京1
        "cn-east-3",      // 上海3
        "cn-east-2",      // 上海2
        "cn-south-1",     // 广州
        "cn-southwest-2", // 重庆
        "ap-southeast-1", // 香港
        "ap-southeast-3", // 新加坡
    }

    // 优先通过 IAM API 自动发现区域
    defer func() {
        if r := recover(); r != nil {
            // 任何异常，使用兜底列表
        }
    }()

    auth := coreauth.NewCredentialsBuilder().
        WithAk(s.AccessKey).
        WithSk(s.AccessSecret).
        Build()

    ir := iamregion.ValueOf("cn-north-4")
    if ir == nil {
        ir = coreRegion.NewRegion("cn-north-4", "https://iam.cn-north-4.myhuaweicloud.com")
    }
    ic := iam.NewIamClient(iam.IamClientBuilder().
        WithRegion(ir).
        WithCredential(auth).
        Build())

    resp, err := ic.KeystoneListRegions(&iammodel.KeystoneListRegionsRequest{})
    if err != nil || resp == nil || resp.Regions == nil || len(*resp.Regions) == 0 {
        return fallback
    }

    out := make([]string, 0, len(*resp.Regions))
    for _, rg := range *resp.Regions {
        if rg.Id != "" {
            out = append(out, rg.Id)
        }
    }
    if len(out) == 0 {
        return fallback
    }
    return out
}

// fetchProjectMap 拉取账户下可用的项目，映射为 region -> projectId
func (s *HuaweiCloudService) fetchProjectMap(ic *iam.IamClient) map[string]string {
    out := map[string]string{}
    // 拉取项目列表
    resp, err := ic.KeystoneListProjects(&iammodel.KeystoneListProjectsRequest{})
    if err != nil || resp == nil || resp.Projects == nil {
        return out
    }
    for _, p := range *resp.Projects {
        if p.Enabled && p.Name != "" && p.Id != "" {
            // 项目名称通常形如 cn-north-4 / ap-southeast-3 等，作为 regionId 使用
            out[p.Name] = p.Id
        }
    }
    return out
}

// GetInstancesJSON 获取华为云实例列表的JSON格式
func (s *HuaweiCloudService) GetInstancesJSON() ([]byte, error) {
    instances, err := s.GetInstances()
    if err != nil {
        return nil, err
    }
    return json.Marshal(instances)
}
