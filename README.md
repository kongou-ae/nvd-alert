# nvd-alert
Scripts for alerting the information of NVD

- nvd-alert
    - check the database of go-cve-dictionary's server by the cpes written in config file.
    - if there is new CVE which corresponded to cpes, inform you.

- nvd-alert-update
    - update the database of go-cve-dictionary's server.

## Usage

1. install `kotakanbe/go-cve-dictionary`
2. set nvd-alert-update to cron
3. make config.json 
4. set nvd-alert to cron

if new CVE appear in NVD, nvd-alert sends mail to you.

## Requirement

- `kotakanbe/go-cve-dictionary`
- SendGrid