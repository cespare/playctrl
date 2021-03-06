# playctrl

**NOTE: As of Chrome 39, playctrl may not be needed for you because Chome has "global scope" shortcuts!**
Install the [Google Play Music
extension](https://chrome.google.com/webstore/detail/google-play-music/fahmaaghhglfmonjliepjlchgpgfmobi) and
change the shortcuts to your liking by clicking 'Keyboard shortcuts' at the bottom of the extensions page
(chrome://extensions). See [this](https://plus.google.com/+FrancoisBeaufort/posts/SDDSiLV5WVF) and
[this](http://www.omgchrome.com/keyboard-control-scope-chrome-extensions/) for more info.

---

Playctrl is a combination of a commandline tool and a chrome extension to enable you control Google Play
(Music) from your desktop.

Playctrl is heavily inspired by [playplay](https://github.com/jsharkey/playplay) but was designed to be easier
for users to install and run.

## How it works

Playctrl consists of two pieces: the `playctrl` tool and a Chrome extension. The `playctrl` tool runs itself
as a daemon and then sends the daemon commands via a unix domain socket. The Chrome extension listens for
server send events from the server. The server just relays commands through.

## Installation

**Tool**

You can download the tool prebuilt for your system (Mac OS or Linux, 64-bit only), at the following urls:

* [Mac OS x64 latest build](http://dl.ctrl-c.us/playctrl-darwin-x64-latest)
* [Linux x64 latest build](http://dl.ctrl-c.us/playctrl-linux-x64-latest)

You will have to `chmod +x` the file after you've downloaded it, and put it somewhere in your `$PATH`.

Alternatively, you can build it yourself from a clone of this repo. You'll need to have Go (1.1+) installed.
Just run `make bin/playctrl` and the binary will be at `bin/playctrl`.

If your Go environment is set up properly, you can use `go get` even more easily:

    $ go get github.com/cespare/playctrl

**Chrome Extension**

You can download the Chrome extension [from the Chrome
store](https://chrome.google.com/webstore/detail/playctrl/loakeafbjkkagnmmlpadfmknpeedckjg).

## Usage

Once you've installed the `playctrl` tool and chrome extension, you'll need to reload Google Play so that the
extension takes effect.

You can run, for example, `playctrl playpause` from the commandline to verify that it works.

You can explicitly run `playctrl start-daemon` to start the daemon. (You don't need to do this; playctrl
launches the daemon the first time you run a command.)  You can run `playctrl stop-daemon` if you want to make
the daemon quit.

You'll want to hook up the shortcuts to global shortcuts in your operating system. For example, I use XFCE
with Ubuntu, so in the Settings Manager I go to Keyboard -> Application Shortcuts and set `playctrl playpause`
to the shortcut `<Primary><Super>space` (that is, ctrl-super-space).

On Mac OS X, you can [create an Automator service and assign a keyboard
shortcut](http://computers.tutsplus.com/tutorials/how-to-launch-any-app-with-a-keyboard-shortcut--mac-31463).

The commands are:

* previous
* playpause
* next
* volumeup
* volumedown
* thumbsup
* thumbsdown

## Developing

Build everything via `make` (you'll need Go installed).

Run the tool with `./bin/playctrl`. You can use `-h` to see the options. You can launch the server in the
foreground with `./bin/playctrl daemon`.
