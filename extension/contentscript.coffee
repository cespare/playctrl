s = document.createElement('script')
s.src = chrome.extension.getURL("playctrl.js")
s.onload = -> @parentNode.removeChild(@)
(document.head || document.documentElement).appendChild(s)
