msgToKeyCode = {
  "previous": 37,
  "playpause": 32,
  "next": 39,
  "volumeup": 187,
  "volumedown": 189
};

function pressKey(k) {
  e = document.createEvent("Events");
  e.initEvent("keydown", true, true);
  e.keyCode = k;
  e.which = k;
  document.dispatchEvent(e);
};

function handleMsg(msg) {
  if (msg.version !== 1) {
    console.warn("playctrl: unhandled protocol version", msg.version);
    return;
  }
  keyCode = msgToKeyCode[msg.value];
  if (keyCode == null) {
    console.warn("playctrl: unknown message", msg.value);
    return;
  }
  console.log("playctrl: dispatching message", msg.value);
  pressKey(keyCode);
}

function connect(extensionId) {
  chrome.runtime.connect(extensionId).onMessage.addListener(function(msg, _, _) {
    console.log("message received");
    console.log(msg);
    handleMsg(msg);
  });
}

addEventListener("message", function(e) {
  if (e.origin === "https://play.google.com" && e.data.magic === "playctrl") {
    connect(e.data.extensionId);
  }
}, false);
