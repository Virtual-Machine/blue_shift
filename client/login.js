$loginButton.addEventListener('click', function(){
	var xhr = new XMLHttpRequest()
	var submitName = $userName.value
	xhr.open('PUT', '/login')
	xhr.setRequestHeader('Content-Type', 'application/json')
	xhr.onload = function() {
		if (xhr.status === 200) {
			var serverResponse = JSON.parse(xhr.responseText)
			if(serverResponse.Type == "Success"){
				// MARKER Client -> Login successful
				$loginButton.parentNode.style.display = 'none'
				window.activeClient = submitName
				window.canvas.setBindings()
				establishSocketConnection(serverResponse.Message)
			} else {
				$warnText.textContent = serverResponse.Message
			}
		}
	}
	xhr.send(JSON.stringify({
		name: submitName,
		password: $userPassword.value
	}))
})