# Remote

![](https://i.imgur.com/lf6Bzjd.png)

Native UI for installing PS4 packages remotely using [`ui`](https://github.com/andlabs/ui) for the frontend and [`net/http`](https://pkg.go.dev/net/http) for the backend.

## Requirements

* Go (1.16)
* [Remote Package Installer](https://github.com/flatz/ps4_remote_pkg_installer)

## Getting Started

```bash
git clone https://github.com/ramadan8/Remote && cd Remote
go build ./src/remote
```

From here, you will have generated a `remote` binary. For Windows users, this will show as a singular `.exe` file.

## Usage

1. Start the Remote Package Installer on your PlayStation 4.
2. Start Remote and select the IPv4 address you would like to use for the host under `Server > Address`.
3. Enter the local IPv4 address of your PlayStation 4 into the `Console > Address` field.
4. Click `Add file` and double click the package file you'd like to install. It **must** end with `.pkg`.