class Coordinate {
	constructor(x, y){
		this.x = x
		this.y = y
	}

	distanceTo(x, y){
		return (Math.abs(x - this.x) + Math.abs(y - this.y))
	}

	getXPixels(scale){
		return this.x * scale
	}

	getYPixels(scale){
		return this.y * scale
	}
}