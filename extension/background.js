function makeHandler(port) {
  return function(e) {
    console.log("message received");
    console.log(e);
    port.postMessage(JSON.parse(e.data));
  };
}

chrome.runtime.onConnectExternal.addListener(function(port) {
  source = new EventSource("http://localhost:49133");
  source.addEventListener("message", makeHandler(port), false);
});
