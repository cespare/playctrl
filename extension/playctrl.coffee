msgToKeyCode =
  "previous": 37
  "playpause": 32
  "next": 39
  "volumeup": 187
  "volumedown": 189

pressKey = (k) ->
  e = document.createEvent("Events")
  e.initEvent("keydown", true, true)
  e.keyCode = k
  e.which = k
  document.dispatchEvent(e)

eventHandler = (e) ->
  msg = JSON.parse(e.data)
  if msg.version != 1
    console.warn "playctrl: unhandled protocol version", msg.version
    return
  keyCode = msgToKeyCode[msg.value]
  if !keyCode?
    console.warn "playctrl: unknown message", msg.value
    return
  console.log "playctrl: dispatching message", msg.value
  pressKey(keyCode)

source = new EventSource("http://localhost:49133")
source.addEventListener "message", eventHandler, false
