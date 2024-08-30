# go-mail-discovery

This is a **quick** and **dirty** implementation of the following mail discovery protocols:

- Autoconfig
- Mobileconfig
- Autodiscover

If you want to contribute to this project, feel free to open a PR!
Considering it is still in very early stages, major changes can be accepted.

## Autoconfig

URL: `https://autoconfig.<domain>.<tld>/mail/config-v1.1.xml` and optionally `?emailaddress=<email>`
Method: `GET`

## Mobileconfig

URL: `https://<domain>.<tld>/email.mobileconfig?email=<email>`
Method: `GET`
Notes: This endpoint is not directly accessed from the mail application, but
is opened via the device's browser, a profile is downloaded and the user
can install it.

## Autodiscover

URL: `https://autodiscover.<domain>.<tld>/autodiscover/autodiscover.xml` or `https://autodiscover.<domain>.<tld>/Autodiscover/Autodiscover.xml`
Body:

```xml
<Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/responseschema/2006">
  <Request>
    <EMailAddress>user@domain.com</EMailAddress>
  </Request>
</Autodiscover>
```

Method: `POST`
