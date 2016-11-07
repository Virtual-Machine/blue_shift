function establishSocketConnection(token) {
	conn = new WebSocket("ws://192.168.5.10:8090/ws?id=" + token)

	conn.onopen = function (evt) {
		console.log("Socket connection established")
		$chatInput.addEventListener('keydown', function(event){
			if (event.which === 13 && $chatInput.value.trim() != ""){
				var message = $chatInput.value
				$chatInput.value = ""
				conn.send(JSON.stringify({type: "ChatMessage",message: message}))
			}
			event.stopPropagation()
		})
	}
	conn.onclose = function (evt) {
		console.log("Socket disconnected")
	}
	conn.onmessage = function (evt) {
		// MARKER Client -> Client is receiving data from socket hub
		var parsedPacket = JSON.parse(evt.data)
		if (parsedPacket instanceof Array){
			console.log("Got map data: ", parsedPacket)
			// TODO Process map data
		} else {
			console.log("Got data packet: ", parsedPacket)
			processPacket(parsedPacket)
		}
	}
	conn.onerror = function (evt) {
	    console.log("Error:", evt)
	}

	window.sConn = conn
}

function processPacket(parsedPacket){
	if(parsedPacket.display_admin_panel){
		$adminPanel.style.display = "block"
		updateUserList(parsedPacket.user_list)
		return
	}
	if(parsedPacket.user_list){
		updateUserList(parsedPacket.user_list)
		return
	}
	if(parsedPacket.author){
		appendChatMessage(parsedPacket.author, parsedPacket.message)
		return
	}
	if(parsedPacket.error){
		appendMessage(parsedPacket.error)
		return
	}
	if(parsedPacket.success){
		appendMessage(parsedPacket.success, parsedPacket.players)
		$adminPanel.style.display = "none"
		window.canvas.state.setActive(parsedPacket.players[0])
		window.canvas.state.start()
		return
	}
	if(parsedPacket.admin_error){
		displayAdminMessage(parsedPacket.admin_error)
		return
	}
}

function appendMessage(message, players){
	var history = $history
	if (history.childNodes.length >= 100){
		history.removeChild(history.firstChild)
	}
	var element = document.createElement('div')
	var span = document.createElement('span')
	element.classList.add("history-text")
	span.classList.add("history-time")
	span.textContent = new Date().toTimeString().split(" ")[0] + " - "
	element.appendChild(span)
	if (players){
		for(var i in players){
			message += " p" + (1 + parseInt(i)) + " : " + players[i]
		}
	}
	var textMessage = document.createTextNode(message)
	element.appendChild(textMessage)
	history.appendChild(element)
	scrollIfNotScrolled(history)
}

function updateUserList(userList){
	var user_list = $userList
	$adminUserList.innerHTML = ""
	user_list.innerHTML = ""
	for(var i in userList){
		var active = userList[i].name === window.activeClient
		var status = userList[i].status
		var element = document.createElement('div')
		var span = document.createElement('span')
		element.classList.add("chat-user")
		if (active) { element.classList.add("active-user") }
		element.textContent = userList[i].name
		span.classList.add("chat-status")
		if(status === "Online") { span.classList.add("status-online") }
		span.textContent = status
		element.appendChild(span)
		user_list.appendChild(element)

		if(userList[i].status === "Online"){
			var userItem = document.createElement('div')
			var text = document.createTextNode(userList[i].name)
			userItem.appendChild(text)
			userItem.setAttribute("draggable", true)
			userItem.classList.add('choice')
			userItem.addEventListener('dragstart', drag)
			$adminUserList.appendChild(userItem)
		}
	}
}

function displayAdminMessage(message){
	$adminMessage.textContent = message
}

function appendChatMessage(author, message){
	var chatDisplay = $chatDisplay
	if (chatDisplay.childNodes.length >= 100){
		chatDisplay.removeChild(chatDisplay.firstChild)
	}
	var element = document.createElement('div')
	var span = document.createElement('span')
	element.classList.add('chat-message')
	element.appendChild(span)
	span.classList.add('chat-id')
	if(author == window.activeClient){ span.classList.add('active-user') }
	span.textContent = author
	var textMessage = document.createTextNode(message)
	element.appendChild(textMessage)
	chatDisplay.appendChild(element)
	scrollIfNotScrolled(chatDisplay)
}

function scrollIfNotScrolled(element) {
	var shouldScroll = ((element.scrollHeight - element.scrollTop - element.clientHeight) < element.clientHeight)
	if(shouldScroll) {
		element.scrollTop = element.scrollHeight
	}
}