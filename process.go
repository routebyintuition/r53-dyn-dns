package main

import (
	"fmt"
	"strings"
	"time"
)

// Process is the main function that checks for the existence of the records in Route 53 and updates as necessary
func (conf *Config) Process(changeDate time.Time) error {
	recordList, err := conf.GetRecords()
	if err != nil {
		Error.Println("Error in Route 53 query: ", err)
		return err
	}

	dnsIP, err := conf.GetPublicIPService()
	if err != nil {
		Error.Println("Error in DNS query: ", err)
		return err
	}

	if len(recordList) == 0 {
		Info.Println("No Route 53 records matching hostname...adding one now.")
		err := conf.ChangeRecord("CREATE", dnsIP.String(), changeDate.String())
		if err != nil {
			return err
		}
	}

	if len(recordList) == 1 {
		if len(recordList[0].ResourceRecords) != 1 {
			err := conf.DeleteRecords()
			if err != nil {
				Error.Printf("Error in deleting Route53 records...(count %d): %s", len(recordList), err)
				return err
			}

			err = conf.ChangeRecord("CREATE", dnsIP.String(), changeDate.String())
			if err != nil {
				Error.Println("Error in creating Route 53 entry: ", err)
				return err
			}
			return nil
		}
		r53IPString := *recordList[0].ResourceRecords[0].Value
		dnsIPString := fmt.Sprintf("%s", dnsIP)
		if !strings.EqualFold(r53IPString, dnsIPString) {
			Info.Printf("Record IPs not matching, performing update of %s to %s \n", r53IPString, dnsIPString)
			err := conf.ChangeRecord("UPSERT", dnsIP.String(), changeDate.String())
			if err != nil {
				Error.Println("Error in updating record set to new address: ", err)
				return err
			}
		}
	}

	if len(recordList) > 1 {
		Info.Println("Found more than one record in the list...will delete the associated records and insert a new one.")
		err := conf.DeleteRecords()
		if err != nil {
			Error.Printf("Error in deleting Route53 records...(count %d): %s", len(recordList), err)
			return err
		}

		err = conf.ChangeRecord("CREATE", dnsIP.String(), changeDate.String())
		if err != nil {
			Error.Println("Error in creating Route 53 entry: ", err)
			return err
		}
	}

	return nil
}
