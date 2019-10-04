# r53-dyn-dns
Dynamic DNS daemon for use with Amazon Route 53

This requires that the AWS CLI be installed on the machine running this application as we make use of the shared credentials file.
In order to use this, you must have credentials configured on your local machine. You can perform this following the link below:

https://docs.aws.amazon.com/cli/latest/userguide/install-linux.html
https://docs.aws.amazon.com/cli/latest/userguide/install-macos.html
https://docs.aws.amazon.com/cli/latest/userguide/install-windows.html

Configure credentials:
https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html

After credentials are created, login to an AWS account where you have access to an Amazon Route 53 Hosted Zone. Create a hosted zone if necessary and ensure that public name resolution is operating as expected.

## Installation

```
go get -u github.com/routebyintuition/r53-dyn-dns.git
```

## Configuration
Create a configuration file (a template is provided in this repository). Update the configuration to reflect your values. An example is shown below:

``` 
hostname = "dns.host.name" 
dns_url = "https://api.ipify.org?format=text"

[aws]
aws_profile = ""
hosted_zone_id = "XXXXXXXXXX"

[server]
refresh_interval = 600
log_directory = "./"
```

**hostname:** Set a hostname matching your Amazon Route 53 hosted zone. If I owned a hosted zone named, "home.name", this would work.

**dns_url:** Set to a public DNS URL that you trust. This must return the public IP resolution you will set to the hostname value. A working example is provided.

**aws_profile:** After creating your shared credentials, leave this blank to use the default profile or enter the name of the created a dedicated profile.

**hosted_zone_id:** This should reflect the Amazon Route 53 hosted zone ID you control.
refresh_interval: This is how often the check is run against the public IP resolution in seconds. 

**log_directory:** Log location directory where service log file is written.

## Run
If the configuration file is named config.toml and placed in the same directory as the running binary, you can just run the application as:
```
./r53-dyn-dns
```

Otherwise, you must pass the configuration file name to the binary:
```
/path/to/r53-dyn-dns -config /another/path/file.toml
```

