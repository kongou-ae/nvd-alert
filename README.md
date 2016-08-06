# nvd-alert
Scripts for alerting the information of NVD

## Usage

1. install `kotakanbe/go-cve-dictionary`
2. set nvd-alert-cron to cron
3. make config.json 
4. set nvd-alert to cron

if new CVE appear in NVD, nvd-alert sends mail to you.

## Requirement

- `kotakanbe/go-cve-dictionary`
- SendGrid