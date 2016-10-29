function establishSocketConnection(token) {
	conn = new WebSocket("ws://192.168.5.10:8090/ws?id=" + token)

	conn.onopen = function (evt) {
		console.log("Socket connection established")
	}
	conn.onclose = function (evt) {
		console.log("Socket disconnected")
	}
	conn.onmessage = function (evt) {
		var parsedPacket = JSON.parse(evt.data)
		if (parsedPacket instanceof Array){
			console.log("Got map data: ", parsedPacket)
		} else {
			console.log("Got data packet: ", parsedPacket)
			if (parsedPacket.count) {
				var history = document.getElementById('history')
				var message = document.createElement('div')
				message.textContent = "There are now " + parsedPacket.count + " active user(s)."
				history.appendChild(message)
				history.scrollTop = history.scrollHeight
			}
		}
	}
	conn.onerror = function (evt) {
	    console.log("Error:", evt)
	}

	window.sConn = conn
}