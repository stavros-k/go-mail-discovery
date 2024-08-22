# go-mail-discovery

## Autoconfig

URL: `https://autoconfig.<domain>.<tld>/mail/config-v1.1.xml` and optionally `?emailaddress=<email>`
Method: `GET`

## Mobileconfig

URL: `https://<domain>.<tld>/email.mobileconfig?email=<email>`
Method: `GET`
Notes: This endpoint is not directly accessed from the mail application, but
is opened via the device's browser, a profile is downloaded and the user
can install it.
