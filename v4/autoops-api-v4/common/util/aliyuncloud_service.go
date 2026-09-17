package util

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// AliyunCloudService 提供阿里云API调用服务
type AliyunCloudService struct {
	AccessKey    string
	AccessSecret string
	Region      string
}

// NewAliyunCloudService 创建新的阿里云服务实例
func NewAliyunCloudService(accessKey, accessSecret, region string) *AliyunCloudService {
	return &AliyunCloudService{
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
		Region:      region,
	}
}

// AliyunInstanceInfo 阿里云实例信息结构
type AliyunInstanceInfo struct {
	InstanceID   string            `json:"instance_id"`
	InstanceName string            `json:"instance_name"` // 实例名称(ECS控制台显示的名称)
	PrivateIPs   []string          `json:"private_ips"`
	PublicIPs    []string          `json:"public_ips"`
	CPU          int               `json:"cpu"`
	MemoryGB     int               `json:"memory_gb"`
	DiskGB       string            `json:"disk_gb"` // 注意: 这里应该返回纯数字，不带"GB"单位
	BillingType  string            `json:"billing_type"`
	OSName       string            `json:"os_name"`
	CreateTime   string            `json:"create_time"`
	ExpireTime   string            `json:"expire_time"`
	Region       string            `json:"region"` // 实例所在区域
	Status       string            `json:"status"` // 实例状态: Running(运行中), Stopped(已停止), Expired(已过期)等
	Tags         map[string]string `json:"tags"`   // 实例标签: key-value格式
}

// GetInstances 获取阿里云实例列表
// 如果 s.Region 为空字符串，则扫描所有区域；否则只查询指定区域
func (s *AliyunCloudService) GetInstances() ([]AliyunInstanceInfo, error) {
	// 验证AccessKey和AccessSecret
	if s.AccessKey == "" || s.AccessSecret == "" {
		return nil, fmt.Errorf("阿里云AccessKey或AccessSecret为空")
	}

	// 创建初始客户端用于查询区域列表
	client, err := ecs.NewClientWithAccessKey("cn-hangzhou", s.AccessKey, s.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云客户端失败: %v", err)
	}

	var allInstances []AliyunInstanceInfo

	// 获取所有可用区域
	regionRequest := ecs.CreateDescribeRegionsRequest()
	regionRequest.Scheme = "https"

	regionResponse, err := client.DescribeRegions(regionRequest)
	if err != nil {
		return nil, fmt.Errorf("查询可用区域失败: %v", err)
	}

	// 遍历每个区域
	for _, regionInfo := range regionResponse.Regions.Region {
		regionID := regionInfo.RegionId

		// 为每个区域创建客户端
		regionClient, err := ecs.NewClientWithAccessKey(regionID, s.AccessKey, s.AccessSecret)
		if err != nil {
			fmt.Printf("[WARN] 创建区域 %s 客户端失败: %v\n", regionID, err)
			continue
		}

		// 分页查询该区域的所有实例
		pageNumber := 1
		pageSize := 100 // 每页100条(阿里云最大值)

		for {
			request := ecs.CreateDescribeInstancesRequest()
			request.Scheme = "https"
			request.PageNumber = requests.NewInteger(pageNumber)
			request.PageSize = requests.NewInteger(pageSize)

			response, err := regionClient.DescribeInstances(request)
			if err != nil {
				fmt.Printf("[WARN] 查询区域 %s 实例失败: %v\n", regionID, err)
				break
			}

			// 如果没有实例了，退出循环
			if len(response.Instances.Instance) == 0 {
				break
			}

			// 处理当前页的实例
			for _, instance := range response.Instances.Instance {
				// 解析私网IP - 经典网络
				privateIPs := make([]string, 0)
				for _, ip := range instance.InnerIpAddress.IpAddress {
					privateIPs = append(privateIPs, ip)
				}
				// VPC网络的私网IP
				for _, ip := range instance.VpcAttributes.PrivateIpAddress.IpAddress {
					privateIPs = append(privateIPs, ip)
				}

				// 解析公网IP
				publicIPs := make([]string, 0)
				for _, ip := range instance.PublicIpAddress.IpAddress {
					publicIPs = append(publicIPs, ip)
				}
				// 弹性公网IP
				if instance.EipAddress.IpAddress != "" {
					publicIPs = append(publicIPs, instance.EipAddress.IpAddress)
				}

				// 解析标签
				tags := make(map[string]string)
				for _, tag := range instance.Tags.Tag {
					tags[tag.TagKey] = tag.TagValue
				}

				// 解析计费类型
				billingType := "按量付费"
				if instance.InstanceChargeType == "PrePaid" {
					billingType = "包年包月"
				}

				instanceInfo := AliyunInstanceInfo{
					InstanceID:   instance.InstanceId,
					InstanceName: instance.InstanceName,
					PrivateIPs:   privateIPs,
					PublicIPs:    publicIPs,
					CPU:          instance.Cpu,
					MemoryGB:     instance.Memory / 1024, // 阿里云返回的是MB，转换为GB
					DiskGB:       "0",                    // 阿里云实例信息中没有直接的磁盘容量字段
					BillingType:  billingType,
					OSName:       instance.OSName,
					CreateTime:   instance.CreationTime,
					ExpireTime:   instance.ExpiredTime,
					Region:       regionID,
					Status:       instance.Status, // Running, Stopped等
					Tags:         tags,
				}

				allInstances = append(allInstances, instanceInfo)
			}

			// 如果返回的实例数少于pageSize，说明已经是最后一页了
			if len(response.Instances.Instance) < pageSize {
				break
			}

			// 移动到下一页
			pageNumber++
		}
	}

	return allInstances, nil
}

// GetInstancesJSON 获取阿里云实例列表的JSON格式
func (s *AliyunCloudService) GetInstancesJSON() ([]byte, error) {
	instances, err := s.GetInstances()
	if err != nil {
		return nil, err
	}
	return json.Marshal(instances)
}
