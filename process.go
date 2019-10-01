package main

import "time"

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
		Info.Println("Checking record matching IP entry: ", dnsIP)
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
