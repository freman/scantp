# ScanTP - Scan to FTP

Many modern scanners scan to FTP, SMB or USB. None of these options really work for me as depending on why I'm scanning I want it to be archived away, run through OCR or shoved in seafile to be shared.

## Just for scanners?

This tool isn't exclusivly for scanners but honestly ftp is an ancient protocol and if you have the option of using literally any other protocol you shuold.

I have deliberately tried to limit the permissions of this tool, it will not permit downloading files or deleting them, it will however currently permit overwriting them, I suspect I will have to change that.

## How it works

You can configure any number of virtual paths in the root of the virtual filesystem that can be powered by different drivers. You can point `documents` at `/home/$user/documents` and `www` at `/var/www`. If you really want to mix it up you can point `cloud` to `https://scanner:password@seafile.example.com` and have your scanned files show up there.

## Configuration

Take a look at the [example config](config.example.toml) you can configure any number of endpoints the only restriction is they have to have unique names.

- host (string)

Host address to listen on, for example "192.168.0.1" for a specific ip, or "0.0.0.0" for all the ips on the box

- port (int)

Port to listen on

- username (string)

Username the scanner will use to log in with

- password (string)

Password the scanner will use to log in with, you have the option of using bcrypt passwords or if you're lazy you can use a plaintext password and just set the plaintext flag to true

- plaintext

Lazy passwords

- path (map)

Define root paths

```toml
[path.$name]
type=$driver
```

Where:
* `$name` is the name of the path to display in `/`
* `$driver` is the name of the driver to use

Depending on the driver you will have vavrious other options you can specify

## Drivers?

### [Seafile](https://www.seafile.com/en/home/) `seafile`

The [seafile driver](drivers/seafile) allows you to upload your files directly to a seafile instance, I recommend using a `scanner` user and sharing a `Scanned Documents` library so you're not leaving important credentials laying around on disk, you can always move the scanned documents out of that shared library.

#### Configuration

- username (string)

Username to log into seafile with

- password (string)

Password to log into seafile with

- api (string)

URL to the seafile installation, without any /api/, /api2/, etc just the base url.

### [File System] `fs` `filesystem` `local`

The file system driver permits you to save your files anywhere on disk, provided the user you're running seafile as has write access to that location.

#### Configuration

- root

Root path

## TODO

- Documentation
- Tests
- Better, more consistant errors
- More drivers
- Either replace the ftp-server library or update go.mod when the patches get merged