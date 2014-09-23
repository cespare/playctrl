// http://stackoverflow.com/questions/9515704/building-a-chrome-extension-inject-code-in-a-page-using-a-content-script/9517879#9517879
s = document.createElement("script");
s.src = chrome.extension.getURL("playctrl.js");
s.onload = function() {
  return this.parentNode.removeChild(this);
};
(document.head || document.documentElement).appendChild(s);
