function getSquaresPositions() {
    let activeFace = getActiveFace(gridContainers);
    let squares = activeFace.querySelectorAll('.draggable-square');
    let result = []
    squares.forEach(square => {
        result.push({
            id: square.id,
            val: square.textContent,
            pos: [
                parseInt(square.getAttribute('data-X-pos'), 10),
                parseInt(square.getAttribute('data-Y-pos'), 10)
            ]
        })
    });
    return result;
}