var $warnText = document.getElementById('warnText')
var $loginButton = document.getElementById('loginSubmit')
var $userName = document.getElementById('uName')
var $userPassword = document.getElementById('uPassword')
var $chatInput = document.getElementById('chatInput')
var $history = document.getElementById('history')
var $userList = document.getElementById('users')
var $chatDisplay = document.getElementById('chatDisplay')
var $minimapCursor = document.getElementById('minimapCursor')


class State {
	constructor(){
		this.active = null
	}

	getActive(){
		return this.active
	}

	setActive(player){
		this.active = player 
	}
}