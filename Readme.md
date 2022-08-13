[![Coverage Status](https://coveralls.io/repos/github/zebox/go-oc-currency-updater/badge.svg?branch=master)](https://coveralls.io/github/zebox/go-oc-currency-updater?branch=master)
### OPENCART CURRENCY UPDATER
Offered ways for fixed it sometimes brake auto update function. Also auto update function, in default module, triggers when admin logged in admin panel.
But I need get currency update daily and automaticaly without login to admin panel. I'm don't have big skills in PHP and OpenCart framework 
and that i decision make simple external CLI util for do it (OpencCartCurrencyUpdate - OCCU).

The util call internal OpenCart API (`route=localisation/currency/refresh`) which refresh currency when user click a refresh buttun 
in 'Settings'-'Localisation'-'Currenct' section. Before call this API user should get `user_token` using `route=common/login` call with admin credentails.

This util required base OpenCart URL and credentials for loggin to admin panel (required parameter set as CLI flags). 
Then  user can schedule execute this utils by either coron or some different schedulers.

### All application Options
```text
-b, --base-url:  Base URL for OpenCart site (https://example.org) [%OCCU_BASE_URL%]
-u, --username:  Username for access to OpenCart admin panel [%OCCU_LOGIN%]
-p, --password:  Password for access to OpenCart admin panel [%OCCU_PASSWORD%]

Help Options:
-?              Show this help message
-h, --help      Show this help message

```

### Credits
[Umputun](https://github.com/umputun) - ideas and patterns  

jessevdk/go-flags for cli parameters parser.