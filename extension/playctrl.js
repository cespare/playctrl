msgToKeyCode = {
  "previous": {code: 37},
  "playpause": {code: 32},
  "next": {code: 39},
  "volumeup": {code: 187},
  "volumedown": {code: 189},
  "thumbsup": {code: 187, alt: true},
  "thumbsdown": {code: 189, alt: true}
};

function pressKey(k) {
  e = document.createEvent("Events");
  e.initEvent("keydown", true, true);
  if (k.alt) {
    e.altKey = true;
  }
  e.keyCode = k.code;
  e.which = k.code;
  document.dispatchEvent(e);
};

function handleMsg(msg) {
  if (msg.version > 2) {
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
