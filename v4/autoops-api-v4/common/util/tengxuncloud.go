package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentCloudService struct {
	AccessKey    string
	AccessSecret string
}

func NewTencentCloudService(accessKey, accessSecret string) *TencentCloudService {
	return &TencentCloudService{
		AccessKey:    strings.TrimSpace(accessKey),
		AccessSecret: strings.TrimSpace(accessSecret),
	}
}

type TencentInstanceInfo struct {
	InstanceID   string            `json:"instance_id"`
	InstanceName string            `json:"instance_name"` // 实例名称(腾讯云控制台显示的名称)
	PrivateIPs   []string          `json:"private_ips"`
	PublicIPs    []string          `json:"public_ips"`
	CPU          int               `json:"cpu"`
	Memory       int               `json:"memory"`
	SystemDisk   int               `json:"system_disk"`
	DataDisks    []int             `json:"data_disks"`
	OSName       string            `json:"os_name"`
	CreateTime   string            `json:"create_time"`
	ExpireTime   string            `json:"expire_time"`
	Region       string            `json:"region"`
	Status       string            `json:"status"`          // 实例状态: RUNNING(运行中), STOPPED(已停止), EXPIRED(已过期)等
	Tags         map[string]string `json:"tags"`            // 实例标签: key-value格式
}

func (s *TencentCloudService) GetInstances() ([]TencentInstanceInfo, error) {
	// 验证AccessKey和AccessSecret
	if s.AccessKey == "" || s.AccessSecret == "" {
		return nil, fmt.Errorf("腾讯云SecretId或SecretKey为空")
	}

	credential := common.NewCredential(s.AccessKey, s.AccessSecret)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

	// 获取所有可用区域
	regionClient, err := cvm.NewClient(credential, "ap-guangzhou", clientProfile)
	if err != nil {
		return nil, fmt.Errorf("创建区域查询客户端失败: %v", err)
	}

	regionRequest := cvm.NewDescribeRegionsRequest()
	regionResponse, err := regionClient.DescribeRegions(regionRequest)
	if err != nil {
		return nil, fmt.Errorf("查询可用区域失败: %v", err)
	}

	var allInstances []TencentInstanceInfo
	
	for _, regionInfo := range regionResponse.Response.RegionSet {
		if regionInfo.Region == nil {
			continue
		}
		regionName := *regionInfo.Region
		client, err := cvm.NewClient(credential, regionName, clientProfile)
		if err != nil {
			continue
		}

		// 分页查询该区域的所有实例
		offset := int64(0)
		limit := int64(100) // 每页100条

		for {
			request := cvm.NewDescribeInstancesRequest()
			request.Offset = &offset
			request.Limit = &limit

			response, err := client.DescribeInstances(request)
			if err != nil {
				fmt.Printf("[WARN] 查询区域 %s 实例失败: %v\n", regionName, err)
				break
			}

			// 如果没有实例了，退出循环
			if response.Response.InstanceSet == nil || len(response.Response.InstanceSet) == 0 {
				break
			}

			fmt.Printf("[DEBUG] 区域 %s: 第 %d-%d 条实例\n", regionName, offset+1, offset+int64(len(response.Response.InstanceSet)))

			for _, instance := range response.Response.InstanceSet {
			privateIPs := make([]string, 0)
			publicIPs := make([]string, 0)

			if instance.PrivateIpAddresses != nil {
				for _, ip := range instance.PrivateIpAddresses {
					if ip != nil {
						privateIPs = append(privateIPs, *ip)
					}
				}
			}
			if instance.PublicIpAddresses != nil {
				for _, ip := range instance.PublicIpAddresses {
					if ip != nil {
						publicIPs = append(publicIPs, *ip)
					}
				}
			}

			// 获取系统盘大小
			systemDiskSize := 0
			if instance.SystemDisk != nil && instance.SystemDisk.DiskSize != nil {
				systemDiskSize = int(*instance.SystemDisk.DiskSize)
			}

			// 获取数据盘大小列表
			dataDiskSizes := make([]int, 0)
			if instance.DataDisks != nil {
				for _, disk := range instance.DataDisks {
					if disk.DiskSize != nil {
						dataDiskSizes = append(dataDiskSizes, int(*disk.DiskSize))
					}
				}
			}

			// 获取实例标签
			tags := make(map[string]string)
			if instance.Tags != nil {
				for _, tag := range instance.Tags {
					if tag.Key != nil && tag.Value != nil {
						tags[*tag.Key] = *tag.Value
					}
				}
			}

			instanceInfo := TencentInstanceInfo{
				PrivateIPs:   privateIPs,
				PublicIPs:    publicIPs,
				Region:       regionName,
				SystemDisk:   systemDiskSize,
				DataDisks:    dataDiskSizes,
				Tags:         tags,
			}

			// 安全处理可能为nil的字段
			if instance.InstanceId != nil {
				instanceInfo.InstanceID = *instance.InstanceId
			}
			if instance.InstanceName != nil {
				instanceInfo.InstanceName = *instance.InstanceName
			}
			if instance.InstanceState != nil {
				instanceInfo.Status = *instance.InstanceState
			}
			if instance.CPU != nil {
				instanceInfo.CPU = int(*instance.CPU)
			}
			if instance.Memory != nil {
				instanceInfo.Memory = int(*instance.Memory)
			}
			if instance.CreatedTime != nil {
				instanceInfo.CreateTime = *instance.CreatedTime
			}
			if instance.ExpiredTime != nil {
				instanceInfo.ExpireTime = *instance.ExpiredTime
			}
			if instance.OsName != nil {
				instanceInfo.OSName = *instance.OsName
			}

			allInstances = append(allInstances, instanceInfo)
			}

			// 如果返回的实例数少于limit，说明已经是最后一页了
			if len(response.Response.InstanceSet) < int(limit) {
				break
			}

			// 移动到下一页
			offset += limit
		}
	}

	return allInstances, nil
}

func (s *TencentCloudService) GetInstancesJSON() ([]byte, error) {
	instances, err := s.GetInstances()
	if err != nil {
		return nil, err
	}
	return json.Marshal(instances)
}
