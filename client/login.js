var warnText = document.getElementById('warnText')
var loginButton = document.getElementById('loginSubmit')

loginButton.addEventListener('click', function(){
	var xhr = new XMLHttpRequest()
	var submitName = document.getElementById('uName').value
	xhr.open('PUT', '/login')
	xhr.setRequestHeader('Content-Type', 'application/json')
	xhr.onload = function() {
		if (xhr.status === 200) {
			var serverResponse = JSON.parse(xhr.responseText)
			if(serverResponse.Type == "Success"){
				loginButton.parentNode.style.display = 'none'
				window.activeClient = submitName
				establishSocketConnection(serverResponse.Message)
			} else {
				warnText.textContent = serverResponse.Message
			}
		}
	}
	xhr.send(JSON.stringify({
		name: submitName,
		password: document.getElementById('uPassword').value
	}))
})