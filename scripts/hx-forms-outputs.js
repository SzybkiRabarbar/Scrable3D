function updateMakePlayHxVals(form) {
    const currentVals = {
        actionType: "makePlay",
        side: getNormalizedRotation(),
        chars: getSquaresPositions(),
    }
    form.setAttribute('hx-vals', JSON.stringify(currentVals));
}