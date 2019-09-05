package service

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"gopub3.0/model"
)

func Tocken() (*sdk.Client, error) {
	client, err := sdk.NewClientWithAccessKey("cn-hongkong", "0ctDrS1xoFucwkoX", "Pu8wQ1lfPfhxrZTK382zF24clOgVQE")
	if err != nil {
		panic(err)
	}
	return client, err
}

func InfoHostList(domain model.DomainAccess, page int) (response *alidns.DescribeDomainRecordsResponse, err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", domain.AccessKeyId, domain.AccessSecret)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = domain.Domain
	request.PageNumber = requests.NewInteger(page)
	request.PageSize = requests.NewInteger(500)

	response, err = client.DescribeDomainRecords(request)
	return
}

func InfoHost(Domain string) (string, error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", "<accessKeyId>", "<accessSecret>")

	request := alidns.CreateDescribeSubDomainRecordsRequest()
	request.Scheme = "https"

	request.SubDomain = Domain

	response, err := client.DescribeSubDomainRecords(request)
	if err != nil {
		return "", err
	}
	return response.RequestId, nil
}

func DelHost(domain model.DomainAccess, DomailId string) (bool, error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", domain.AccessKeyId, domain.AccessSecret)

	request := alidns.CreateDeleteDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = DomailId

	_, err = client.DeleteDomainRecord(request)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewHost(domain string, subdomain string, ip string) (bool, error) {
	client, err := alidns.NewClientWithAccessKey("cn-hongkong", "0ctDrS1xoFucwkoX", "Pu8wQ1lfPfhxrZTK382zF24clOgVQE")

	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"

	request.DomainName = domain
	request.RR = subdomain
	request.Type = "A"
	request.Value = ip

	_, err = client.AddDomainRecord(request)
	if err != nil {
		return false, err
	}
	return true, nil
}
