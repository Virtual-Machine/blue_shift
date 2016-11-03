var $adminPanel = document.getElementById('adminPanel')
var $adminMessage = document.getElementById('adminMessage')
var $adminUserList = document.getElementById('adminUserList')
var $p1 = document.getElementById('p1')
var $p2 = document.getElementById('p2')
var $p3 = document.getElementById('p3')
var $p4 = document.getElementById('p4')

function allowDrop(ev) {
    ev.preventDefault()
}

function drag(ev) {
	var name = ev.target.textContent
	var names = [$p1.value, $p2.value, $p3.value, $p4.value]
	for(var i in names){
		if (names[i] === name) {
			$adminMessage.textContent = "This name has already been selected"
			ev.preventDefault()
			setTimeout(function(){
				$adminMessage.textContent = "Drag 2-4 Player Names Into The Relevant Slots And Click Submit To Begin A Game"
			},2000)
			return
		}
	}
    ev.dataTransfer.setData("text", name)
}

function drop(ev) {
    ev.preventDefault()
    var data = ev.dataTransfer.getData("text")
    ev.target.value = data
}

function adminClear(){
	$p1.value = ""
	$p2.value = ""
	$p3.value = ""
	$p4.value = ""
}

function adminSubmit(){
	var sendString = ""
	var names = [$p1.value, $p2.value, $p3.value, $p4.value]
	for(var i in names){
		if(names[i].trim().length > 0){
			sendString += ";" + names[i]
		}
	}
	sendString = sendString.substring(1)
	console.log(sendString)
	window.sConn.send(JSON.stringify({type: "StartGame",message: sendString}))
}