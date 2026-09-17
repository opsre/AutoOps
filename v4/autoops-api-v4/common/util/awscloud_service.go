package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// AWSCloudService 提供AWS云API调用服务
type AWSCloudService struct {
	AccessKey    string
	AccessSecret string
}

// NewAWSCloudService 创建新的AWS云服务实例
func NewAWSCloudService(accessKey, accessSecret string) *AWSCloudService {
	return &AWSCloudService{
		AccessKey:    strings.TrimSpace(accessKey),
		AccessSecret: strings.TrimSpace(accessSecret),
	}
}

// AWSInstanceInfo AWS实例信息结构
type AWSInstanceInfo struct {
	InstanceID   string            `json:"instance_id"`
	InstanceName string            `json:"instance_name"` // 实例名称(从Name标签获取)
	PrivateIPs   []string          `json:"private_ips"`
	PublicIPs    []string          `json:"public_ips"`
	CPU          int               `json:"cpu"`       // vCPU数量
	Memory       int               `json:"memory"`    // 内存(GB)
	DiskGB       string            `json:"disk_gb"`   // 磁盘大小(GB)
	OSName       string            `json:"os_name"`   // 操作系统
	CreateTime   string            `json:"create_time"`
	Region       string            `json:"region"`     // 实例所在区域
	Status       string            `json:"status"`     // 实例状态: running, stopped, terminated等
	InstanceType string            `json:"instance_type"` // 实例类型: t2.micro, m5.large等
	Tags         map[string]string `json:"tags"`       // 实例标签: key-value格式
}

// GetInstances 获取AWS实例列表(扫描所有区域)
func (s *AWSCloudService) GetInstances() ([]AWSInstanceInfo, error) {
	// 验证AccessKey和AccessSecret
	if s.AccessKey == "" || s.AccessSecret == "" {
		return nil, fmt.Errorf("AWS AccessKey或SecretKey为空")
	}

	ctx := context.Background()

	// 创建AWS配置
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.AccessKey,
			s.AccessSecret,
			"",
		)),
		config.WithRegion("us-east-1"), // 默认区域用于获取区域列表
	)
	if err != nil {
		return nil, fmt.Errorf("创建AWS配置失败: %v", err)
	}

	// 创建EC2客户端获取所有区域
	ec2Client := ec2.NewFromConfig(cfg)
	regionsOutput, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false), // 只获取启用的区域
	})
	if err != nil {
		return nil, fmt.Errorf("获取AWS区域列表失败: %v", err)
	}

	fmt.Printf("[DEBUG] 查询到 %d 个AWS区域\n", len(regionsOutput.Regions))

	var allInstances []AWSInstanceInfo

	// 遍历每个区域
	for _, region := range regionsOutput.Regions {
		if region.RegionName == nil {
			continue
		}
		regionName := *region.RegionName

		// 为每个区域创建配置和客户端
		regionCfg, err := config.LoadDefaultConfig(ctx,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				s.AccessKey,
				s.AccessSecret,
				"",
			)),
			config.WithRegion(regionName),
		)
		if err != nil {
			fmt.Printf("[WARN] 创建区域 %s 配置失败: %v\n", regionName, err)
			continue
		}

		regionClient := ec2.NewFromConfig(regionCfg)

		// 查询该区域的所有实例
		instances, err := s.getRegionInstances(ctx, regionClient, regionName)
		if err != nil {
			fmt.Printf("[WARN] 查询区域 %s 实例失败: %v\n", regionName, err)
			continue
		}

		if len(instances) > 0 {
			fmt.Printf("[INFO] 区域 %s: 找到 %d 台实例\n", regionName, len(instances))
			allInstances = append(allInstances, instances...)
		}
	}

	fmt.Printf("[INFO] 从所有区域共获取到 %d 台AWS实例\n", len(allInstances))
	return allInstances, nil
}

// getRegionInstances 获取指定区域的所有实例
func (s *AWSCloudService) getRegionInstances(ctx context.Context, client *ec2.Client, regionName string) ([]AWSInstanceInfo, error) {
	var instances []AWSInstanceInfo
	var nextToken *string

	for {
		// 查询实例
		input := &ec2.DescribeInstancesInput{
			MaxResults: aws.Int32(1000), // 每次最多1000个实例
			NextToken:  nextToken,
		}

		result, err := client.DescribeInstances(ctx, input)
		if err != nil {
			return nil, err
		}

		// 解析实例信息
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				instanceInfo := s.parseInstance(instance, regionName)
				instances = append(instances, instanceInfo)
			}
		}

		// 检查是否还有下一页
		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return instances, nil
}

// parseInstance 解析AWS实例信息
func (s *AWSCloudService) parseInstance(instance types.Instance, regionName string) AWSInstanceInfo {
	info := AWSInstanceInfo{
		Region:       regionName,
		PrivateIPs:   []string{},
		PublicIPs:    []string{},
		Tags:         make(map[string]string),
		InstanceType: string(instance.InstanceType),
	}

	// 实例ID
	if instance.InstanceId != nil {
		info.InstanceID = *instance.InstanceId
	}

	// 实例名称(从Name标签获取)
	for _, tag := range instance.Tags {
		if tag.Key != nil && tag.Value != nil {
			info.Tags[*tag.Key] = *tag.Value
			if *tag.Key == "Name" {
				info.InstanceName = *tag.Value
			}
		}
	}

	// 如果没有Name标签，使用InstanceID作为名称
	if info.InstanceName == "" {
		info.InstanceName = info.InstanceID
	}

	// IP地址
	if instance.PrivateIpAddress != nil {
		info.PrivateIPs = append(info.PrivateIPs, *instance.PrivateIpAddress)
	}
	if instance.PublicIpAddress != nil {
		info.PublicIPs = append(info.PublicIPs, *instance.PublicIpAddress)
	}

	// 状态
	if instance.State != nil && instance.State.Name != "" {
		info.Status = string(instance.State.Name)
	}

	// CPU和内存(从实例类型推断)
	info.CPU, info.Memory = s.getInstanceTypeSpecs(string(instance.InstanceType))

	// 磁盘信息(从块设备映射获取)
	totalDiskGB := 0
	if instance.BlockDeviceMappings != nil {
		for _, bdm := range instance.BlockDeviceMappings {
			if bdm.Ebs != nil && bdm.Ebs.VolumeId != nil {
				// 默认每个卷为8GB(需要通过DescribeVolumes API获取精确大小，这里简化处理)
				totalDiskGB += 8
			}
		}
	}
	if totalDiskGB == 0 {
		totalDiskGB = 8 // 默认8GB
	}
	info.DiskGB = strconv.Itoa(totalDiskGB)

	// 操作系统(从平台和镜像信息推断)
	if instance.Platform != "" {
		info.OSName = string(instance.Platform)
	} else {
		info.OSName = "Linux" // 默认Linux
	}

	// 创建时间
	if instance.LaunchTime != nil {
		info.CreateTime = instance.LaunchTime.Format("2006-01-02T15:04Z")
	}

	return info
}

// getInstanceTypeSpecs 根据实例类型获取CPU和内存规格(简化版)
func (s *AWSCloudService) getInstanceTypeSpecs(instanceType string) (cpu int, memoryGB int) {
	// 简化的实例类型规格映射(实际应该查询AWS API或使用完整映射表)
	specs := map[string]struct{ cpu, memory int }{
		// t系列 - 通用型突发性能实例
		"t2.nano":    {1, 1},
		"t2.micro":   {1, 1},
		"t2.small":   {1, 2},
		"t2.medium":  {2, 4},
		"t2.large":   {2, 8},
		"t2.xlarge":  {4, 16},
		"t2.2xlarge": {8, 32},
		"t3.nano":    {2, 1},
		"t3.micro":   {2, 1},
		"t3.small":   {2, 2},
		"t3.medium":  {2, 4},
		"t3.large":   {2, 8},
		"t3.xlarge":  {4, 16},
		"t3.2xlarge": {8, 32},
		// m系列 - 通用型
		"m5.large":    {2, 8},
		"m5.xlarge":   {4, 16},
		"m5.2xlarge":  {8, 32},
		"m5.4xlarge":  {16, 64},
		"m5.8xlarge":  {32, 128},
		"m5.12xlarge": {48, 192},
		"m5.16xlarge": {64, 256},
		"m5.24xlarge": {96, 384},
		// c系列 - 计算优化型
		"c5.large":    {2, 4},
		"c5.xlarge":   {4, 8},
		"c5.2xlarge":  {8, 16},
		"c5.4xlarge":  {16, 32},
		"c5.9xlarge":  {36, 72},
		"c5.12xlarge": {48, 96},
		"c5.18xlarge": {72, 144},
		"c5.24xlarge": {96, 192},
		// r系列 - 内存优化型
		"r5.large":    {2, 16},
		"r5.xlarge":   {4, 32},
		"r5.2xlarge":  {8, 64},
		"r5.4xlarge":  {16, 128},
		"r5.8xlarge":  {32, 256},
		"r5.12xlarge": {48, 384},
		"r5.16xlarge": {64, 512},
		"r5.24xlarge": {96, 768},
	}

	if spec, ok := specs[instanceType]; ok {
		return spec.cpu, spec.memory
	}

	// 默认值(如果不在映射表中)
	return 2, 4
}

// GetInstancesJSON 获取JSON格式的实例列表
func (s *AWSCloudService) GetInstancesJSON() ([]byte, error) {
	instances, err := s.GetInstances()
	if err != nil {
		return nil, err
	}
	return json.Marshal(instances)
}
