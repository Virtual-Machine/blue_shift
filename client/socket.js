conn = new WebSocket("ws://localhost:8090/ws")

conn.onopen = function (evt) {
    console.log("Opening connection")
    console.log(evt)
}
conn.onclose = function (evt) {
    console.log("Closing connection")
    console.log(evt)
}
conn.onmessage = function (evt) {
    console.log("A message was received")
    console.log(evt)
}
conn.onerror = function (evt) {
    console.log("An error occurred")
    console.log(evt)
}