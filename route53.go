package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func (conf *Config) initAWS() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", conf.Profile),
	})
	if err != nil {
		return err
	}

	conf.Route53 = route53.New(sess)

	return nil
}

// getRecord returns all of the associated records matching the request but also includes
// records not related so must be filtered. For results that DO NOT need to be filtered,
// use the GetRecords method.
func (conf *Config) getRecord() (*route53.ListResourceRecordSetsOutput, error) {
	listParams := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(conf.ZoneID),
		MaxItems:        aws.String("100"),
		StartRecordName: aws.String(conf.Hostname),
		StartRecordType: aws.String("A"),
	}
	respList, err := conf.Route53.ListResourceRecordSets(listParams)

	if err != nil {
		return nil, err
	}

	return respList, nil
}

// GetRecords returns a filtered list of records that match the requested query
func (conf *Config) GetRecords() ([]*route53.ResourceRecordSet, error) {
	recordList, err := conf.getRecord()
	if err != nil {
		return nil, err
	}

	recordListUpdate := []*route53.ResourceRecordSet{}

	for _, recordSet := range recordList.ResourceRecordSets {
		setName := *recordSet.Name
		getName := fmt.Sprintf("%s%s", conf.Hostname, ".")

		if strings.EqualFold(getName, setName) {
			recordListUpdate = append(recordListUpdate, recordSet)
		}
	}

	return recordListUpdate, nil
}

// ChangeRecord will make adjustments to the Amazon Route 53 record set
func (conf *Config) ChangeRecord(changeType string, address string, changeDate string) error {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(changeType),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(conf.Hostname),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(address),
							},
						},
						TTL:           aws.Int64(60),
						SetIdentifier: aws.String(fmt.Sprintf("Dynamic DNS record updated on: %s", changeDate)),
					},
				},
			},
			Comment: aws.String("Sample update."),
		},
		HostedZoneId: aws.String(conf.ZoneID),
	}
	_, err := conf.Route53.ChangeResourceRecordSets(params)

	if err != nil {
		return err
	}

	return nil
}

// DeleteRecords will delete the associated records in the event that there are issues and need to be re-populated.
// One such circumstance is in the event that there are more than one record.
func (conf *Config) DeleteRecords() error {
	recordList, err := conf.getRecord()
	if err != nil {
		return err
	}

	if *recordList.IsTruncated {
		return errors.New("Over 100 records for this entry...something is wrong")
	}

	changeBatchList := []*route53.Change{}

	for _, recordSet := range recordList.ResourceRecordSets {
		setName := *recordSet.Name
		getName := fmt.Sprintf("%s%s", conf.Hostname, ".")
		if strings.EqualFold(getName, setName) {
			recordSetToUpdate := route53.Change{}
			recordSetToUpdate.SetAction("DELETE")
			recordSetToUpdate.SetResourceRecordSet(recordSet)
			changeBatchList = append(changeBatchList, &recordSetToUpdate)
		}
	}

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changeBatchList,
			Comment: aws.String("Deleting records to reset config."),
		},
		HostedZoneId: aws.String(conf.ZoneID),
	}
	_, err = conf.Route53.ChangeResourceRecordSets(params)

	if err != nil {
		return err
	}

	return nil
}
