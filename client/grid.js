class Grid {
	constructor(width, height, scale){
		let cellsWide = width / scale
		let cellsHigh = height / scale
		this.scale = scale
		this.cells = {}
		for(var i = 0; i < cellsWide; i++){
			for(var k = 0; k < cellsHigh; k++){
				this.cells[i + "-" + k] = new Coordinate(i, k)
			}
		}
	}
}