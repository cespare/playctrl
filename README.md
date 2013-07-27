# playctrl

Playctrl is a combination of a small webserver and a chrome extension to enable you control Google Play
(Music) from your desktop.

Playctrl is heavily inspired by [playplay](https://github.com/jsharkey/playplay) but was designed to be easier
for users to install and run.

## How it works

Playctrl consists of three pieces: the server (`playctrld`), the client (`playctrl`), and a Chrome extension.
The server listens for UDP commands from the client. The Chrome extension listens for server send events from
the server. The server just relays commands through.

## Usage

Convenient downloads and Ubuntu packages coming soon.

Once you've installed the client, server, and chrome extension, and the server is running, you'll need to
reload Google Play so that the extension takes effect.

You can run, for example, `playctrl playpause` from the commandline to verify that it works.

You'll want to hook up the shortcuts to global shortcuts in your operating system. For example, I use XFCE
with Ubuntu, so in the Settings Manager I go to Keyboard -> Application Shortcuts and set `playctrl playpause`
to the shortcut `<Primary><Super>space` (that is, ctrl-super-space).

The commands are:

* previous
* playpause
* next
* volumeup
* volumedown

## Developing

Build everything via `make` (you'll need Go and Coffeescript installed).

Run the server with `./bin/playctrld`. Run the client with `./bin/playctrl`. You can use `-h` to see the
options. Typically you won't give the server any options, and you'll run the client like so: `./bin/playctrl
playpause`.
