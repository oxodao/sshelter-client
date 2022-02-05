# SSHelter Client


SSHelter is a simple tool letting you sync across your computers parts of your ssh config like the host machines and specific config for each (like Port Forwarding). It was inspired by the really neat [Termius](https://termius.com/) but since it's a bit pricy for simple hobbyists, I decided to make my own implementation of the features I most needed.

This does not sync your SSH keys yet, as I'm not comfortable to ensure the security required for this.

This is the client for the [SSHelter](https://github.com/oxodao/sshelter) server. Please don't mix&match the client and server version. Be sure to be always on the same version for both.

## Installation

### Linux
Requires [zenity](https://help.gnome.org/users/zenity/index.html.fr):
```sh
$ pacman -S zenity # Arch-based distros
$ sudo apt install zenity # Ubuntu-based distros
```

Simply download the release for your OS in the release tab and put it somewhere:

```sh
$ sudo curl -o /usr/bin/sshelter https://github.com/oxodao/sshelter-client/releases/download/v0.1/sshelter-0.1-linux-x86_64
$ sudo chmod +x /usr/bin/sshelter
```

Then create the config file:
`$HOME/.config/sshelter/sshelter.yml` (Linux)
```yaml
server:
    url: https://sshelter.localhost # Fill with your SSHelter base URL
```

### OS X 
The software has not been ported to OS X but it is planned

### Windows
The software has not been ported to Windows but it might be planned

## Usage

The first time you'll use a command that requires access to the server, your username/password will be prompted. It should not be required anymore after this and the password is not stored anywhere.

### Listing machines
This command will list every machines you have setup on your account
```sh
$ sshelter list
```

### Creating a machine
```sh
$ sshelter create --name "My machine" --shortname machine --host my.machine.tld
```
### Adding forwarded ports
```sh
$ sshelter port add --by-shortname machine --lh localhost --lp 1234 --rh localhost --rp 1234
```

The `--reverse` flag can be used to map a local port on the remote machine instead of mapping a remote port on the local machine.

### Removing a forwarded port
```sh
$ sshelter port delete --index 0
```

The index is listed between `[]` when you list machines

### Syncing
```sh
$ sshelter sync
```
This will rewrite your ssh config appending/updating the SSHelter machines to it.


### Exporting one or more machine(s)
```sh
$ sshelter export --by-shortname "shortnames,separated,by,commas"
```
This will export the machine creation file. This can be used as a backup for your account, or to share it with your team. This file should be directly ready to use with the create command.

### Deleting a machine
```sh
$ sshelter delete --by-shortname "machine shortname"
```
### Editing machines
```sh
$ sshelter edit --by-shortname machine --shortname newshortname
```

## Roadmap
- [x] Implementing basic sync
- [x] Implementing a way to create machines
- [x] Implementing a way to delete machines
- [x] Implementing a way to export machines
- [x] Implementing the ssh config writer
- [x] Implement remote forward ports
- [ ] Implementing a simple system tray icon to keep everything always in sync + start the port forwarding
- [ ] Add a --no-gui flag for people who wants to only use it in CLI (And because zenity MIGHT not work / be too much hassle under WSL)
- [ ] Porting to OS X
- [ ] Porting to Windows
- [ ] General system stability improvements to enhance the user's experience.

## License
> SSHelter Client - Simple ssh config sync software
> Copyright (C) 2021 - Oxodao
> 
> This program is free software: you can redistribute it and/or modify
> it under the terms of the GNU General Public License as published by
> the Free Software Foundation, either version 3 of the License, or
> (at your option) any later version.
> 
> This program is distributed in the hope that it will be useful,
> but WITHOUT ANY WARRANTY; without even the implied warranty of
> MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
> GNU General Public License for more details.
> 
> You should have received a copy of the GNU General Public License
> along with this program.  If not, see <https://www.gnu.org/licenses/>.

## Contributing

So yeah this is only a simple side project, if you want to contribute by all mean go ahead, no guidelines or so just don't be a jerk. I reserve the right to refuse any PR without explanation or stuff like that. Maybe one day this software will grow in feature-set enough for me to consider it as a **real** project and have a clearer way of handling this stuff.
