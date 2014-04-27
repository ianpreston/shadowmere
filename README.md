# shadowmere

Shadowmere is a really barebones, postgresql-backed IRC services package for UnrealIRCd. Very much in a work-in-progress, pre-alpha state.

## Install

Currently `shadowmere` only supports UnrealIRCd, though support for more IRCds may be added in the future.

1) In `unrealircd.conf`, add a link block:

    link your.services.hostname
    {
        username        *;
        hostname        your.services.hostname;
        bind-ip         *;
        port            6668;
        hub             *;
        password-connect "foo";
        password-receive "foo";
        class           servers;
    };

2) In `unrealircd.conf`, add a U-Line:

    ulines {
        your.services.hostname;
    };

3) Open `main.go` and edit the configuration directives

4) Configure the PostgreSQL database:

    $ psql -c "CREATE DATABASE shadowmere;"
    $ psql -d shadowmere -f scripts/schema.sql

Finally, start the services daemon (in production, you'll want to use `supervisord` or something similar):

    $ go run main.go

## Usage

    /msg NickServ REGISTER email password
    /msg NickServ IDENTIFY password

## License

`shadowmere` is provided under the [MIT License](http://en.wikipedia.org/wiki/MIT_License "MIT License")