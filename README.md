# playctrl

Playctrl is a combination of a commandline tool and a chrome extension to enable you control Google Play
(Music) from your desktop.

Playctrl is heavily inspired by [playplay](https://github.com/jsharkey/playplay) but was designed to be easier
for users to install and run.

## How it works

Playctrl consists of two pieces: the `playctrl` tool and a Chrome extension. The `playctrl` tool runs itself
as a daemon and then sends the daemon commands via a unix domain socket. The Chrome extension listens for
server send events from the server. The server just relays commands through.

## Usage

Convenient downloads and Ubuntu packages coming soon.

Once you've installed the `playctrl` tool and chrome extension, you'll need to reload Google Play so that the
extension takes effect.

You can run, for example, `playctrl playpause` from the commandline to verify that it works.

You can explicitly run `playctrl start-daemon` to start the daemon. You can run `playctrl stop-daemon`
if you want to make the daemon quit.

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

Run the tool with `./bin/playctrl`. You can use `-h` to see the options. You can launch the server in the
foreground with `./bin/playctrl daemon`.
