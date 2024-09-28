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

---

It uses MX Lookup to determine the provider ID.
It goes over the MX records, extracts the **host** and this becomes the ID of the provider.

For example: `mx.zoho.eu.` becomes `zoho.eu`.
So a provider would have to be defined in the config file like so:

```yaml
# In order to use the config file you need to set the env variable MAIL_PROVIDER_CONFIG_PATH
# otherwise a default provider (zoho.eu) will be used
providers:
  - id: zoho.eu # Notice the id here
    imap_server:
      hostname: imappro.zoho.eu
      port: 993
      socket_type: SSL
      authentication: password-cleartext
    pop3_server:
      hostname: poppro.zoho.eu
      port: 995
      socket_type: SSL
      authentication: password-cleartext
    smtp_server:
      hostname: smtppro.zoho.eu
      port: 587
      socket_type: STARTTLS
      authentication: password-cleartext
      use_global_preferred_server: false
```

---

## Credits

Big inspirations and help from the following projects:

- [email-autoconf](https://gitlab.com/onlime/email-autoconf)
- [automx2](https://github.com/rseichter/automx2)

## Resources

- [Apple Configuration Profile Reference](https://developer.apple.com/business/documentation/Configuration-Profile-Reference.pdf)
