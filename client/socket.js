function establishSocketConnection(token) {
	conn = new WebSocket("ws://192.168.5.10:8090/ws?id=" + token)

	conn.onopen = function (evt) {
		console.log("Socket connection established")
	}
	conn.onclose = function (evt) {
		console.log("Socket disconnected")
	}
	conn.onmessage = function (evt) {
		var socketPacket = evt.data
		console.log("Got socket packet data: ", socketPacket)
		var parsedPacket = JSON.parse(evt.data)
		console.log("Got parsed: ", parsedPacket)
	}
	conn.onerror = function (evt) {
	    console.log("Error:", evt)
	}

	window.sConn = conn
}