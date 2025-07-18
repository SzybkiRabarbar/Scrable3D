package ctrl

import (
	"bytes"
	"container/heap"
	"errors"
	"fmt"
	"log"
	"scrable3/internal/cfg"
	"scrable3/internal/common"
	"scrable3/internal/dto"
	"scrable3/internal/model"
	"scrable3/internal/svc"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

type GameController interface {
	GetCurrentFields(ctx *dto.WsContext) ([]byte, error)
	GetAvaibleChars(ctx *dto.WsContext) (senderResponse []byte, err error)
	ReceiveChars(
		ctx *dto.WsContext,
		p *dto.PlayData,
	) (
		broadcastResponse []byte,
		senderResponse []byte,
		err error,
	)
}

type gameController struct {
	wordsController WordsController
	playerService   svc.PlayerService
	fieldService    svc.FieldService
	avCharService   svc.AvCharService
}

func NewGameController(
	wordsController WordsController,
	playerService svc.PlayerService,
	fieldService svc.FieldService,
	avCharService svc.AvCharService,
) GameController {
	return &gameController{
		wordsController: wordsController,
		playerService:   playerService,
		fieldService:    fieldService,
		avCharService:   avCharService,
	}
}

// |PRIVATE| //

// 1 is horizontal, 0 is vertical
func (gc *gameController) whichAxisIsStraight(chars *[]dto.Char) int {
	firstElementNumber := (*chars)[0].Position[0]
	lastElementNumber := (*chars)[len(*chars)-1].Position[0]
	if firstElementNumber == lastElementNumber {
		return 0
	} else {
		return 1
	}
}

func (gc *gameController) findDepthFromLeftMostPosition(depthMap *map[int]int) int {
	var value int
	var ok bool
	for i := range cfg.BOARD_SIZE {
		value, ok = (*depthMap)[int(i)]
		if ok {
			return value
		}
	}
	return -1
}

// Maps characters to their positions along a specified axis, considering their
// depth positions. The 'sideInt' and 'straightAxisId' parameters determines
// which axis (X, Y, Z) is used as 'changingPos', 'staticPos', and 'depthPos'.
//
// The 'depthPos' varies based on the value of 'sideInt':
//
//	0   (smallest Z as depth) // North
//	90  (highest Y as depth) // East
//	180 (highest Z as depth)  // South
//	270 (smallest Y as depth)  // West
//
// The 'staticPos' varies based on the value of 'straightAxisId':
//
//	0 (Y as static)
//	1 (X or Z as static)
func (gc *gameController) makeCharsInStraightAxisFromFieldsMap(
	fields *[]model.Field,
	straightAxisNumber int,
	straightAxisId int,
	sideInt int,
) (*map[int]string, *map[int]int) {
	existingCharsInStraightAxis := make(map[int]string)
	currentDepthMap := make(map[int]int)

	for _, field := range *fields {
		// Pointer to the axis position that changes based on 'sideInt' and
		// 'straightAxisId'
		var changingPos *int
		// Pointer to the axis position that remains static based on
		// 'straightAxisNumber'
		var staticPos *int
		// Pointer to the axis position used for depth comparison based on
		// 'sideInt'
		var depthPos *int

		if sideInt == 0 {
			depthPos = &field.PosZ
			if straightAxisId == 0 {
				// @@ Check
				changingPos = &field.PosX
				staticPos = &field.PosY
			} else {
				// @@ Check
				changingPos = &field.PosY
				staticPos = &field.PosX
			}
		} else if sideInt == 90 {
			depthPos = &field.PosY
			if straightAxisId == 0 {
				// @@ Check
				var temp = cfg.BOARD_SIZE - 1 - field.PosZ
				changingPos = &field.PosX
				staticPos = &temp
				// log.Println(*changingPos, *staticPos, *depthPos)
			} else {
				changingPos = &field.PosX
				staticPos = &field.PosZ
			}
		} else if sideInt == 180 {
			depthPos = &field.PosZ
			if straightAxisId == 0 {
				changingPos = &field.PosX
				staticPos = &field.PosY
			} else {
				changingPos = &field.PosY
				staticPos = &field.PosX
			}
		} else if sideInt == 270 {
			depthPos = &field.PosY
			if straightAxisId == 0 {
				// @@ Check
				changingPos = &field.PosX
				staticPos = &field.PosZ
			} else {
				// @@ Check
				changingPos = &field.PosZ
				staticPos = &field.PosX
			}
		}

		// If the field's static position does not match the specified straight
		// axis number, skip this field.
		if *staticPos != straightAxisNumber {
			continue
		}

		currDepthPos, ok := currentDepthMap[*changingPos]
		if !ok { // initialize the depth position
			currentDepthMap[*changingPos], currDepthPos = *depthPos, *depthPos
		}
		// Update the depth position and character mapping if the current depth
		// position meets the criteria based on 'sideInt'
		if ((sideInt == 90 || sideInt == 180) && currDepthPos <= *depthPos) ||
			((sideInt == 0 || sideInt == 270) && currDepthPos >= *depthPos) {
			currentDepthMap[*changingPos] = *depthPos
			existingCharsInStraightAxis[*changingPos] = field.Value
		}
	}
	return &existingCharsInStraightAxis, &currentDepthMap
}

// Generates a dto.FieldData structure based on the provided character data,
// side orientation, depth level and word orientation.
func (gc *gameController) createFieldData(
	char *dto.Char,
	sideInt int,
	depthLevel int,
	isHorizonatal int,
) *dto.FieldData {
	fieldData := dto.FieldData{
		Value: char.Value,
		Pos:   [3]int{},
	}
	log.Println("\n*** createFieldData ***")
	log.Println("sideInt")
	log.Println(sideInt)
	log.Println("char.Value")
	log.Println(char.Value)
	log.Println("char.Position[1]")
	log.Println(char.Position[1])
	log.Println("char.Position[0]")
	log.Println(char.Position[0])
	log.Println("depthLevel")
	log.Println(depthLevel)
	log.Println("isHorizonatal")
	log.Println(isHorizonatal)
	log.Println("\n")

	if isHorizonatal == 1 {
		// Word is horizontal (along X axis)
		if sideInt == 0 {
			fieldData.Pos[0] = char.Position[1]
			fieldData.Pos[1] = char.Position[0]
			fieldData.Pos[2] = depthLevel
		} else if sideInt == 90 {
			fieldData.Pos[0] = depthLevel
			fieldData.Pos[1] = char.Position[1]
			fieldData.Pos[2] = cfg.BOARD_SIZE - 1 - char.Position[0]
		} else if sideInt == 180 {
			fieldData.Pos[0] = cfg.BOARD_SIZE - 1 - char.Position[1]
			fieldData.Pos[1] = cfg.BOARD_SIZE - 1 - char.Position[0]
			fieldData.Pos[2] = depthLevel
		} else { // sideInt == 270
			// @@ Check
			fieldData.Pos[0] = char.Position[1]
			fieldData.Pos[1] = depthLevel
			fieldData.Pos[2] = char.Position[0]
		}
	} else {
		// Word is vertical (along Y axis)
		if sideInt == 0 {
			fieldData.Pos[0] = char.Position[1]
			fieldData.Pos[1] = char.Position[0]
			fieldData.Pos[2] = depthLevel
		} else if sideInt == 90 {
			fieldData.Pos[0] = char.Position[1]
			fieldData.Pos[1] = char.Position[0]
			fieldData.Pos[2] = depthLevel
		} else if sideInt == 180 {
			fieldData.Pos[0] = char.Position[1]
			fieldData.Pos[1] = char.Position[0]
			fieldData.Pos[2] = depthLevel
		} else { // sideInt == 270
			// @@ Check
			fieldData.Pos[0] = char.Position[1]
			fieldData.Pos[1] = depthLevel
			fieldData.Pos[2] = char.Position[0]
		}
	}
	log.Println(fieldData)

	return &fieldData
}

func (gc *gameController) obtainWordAndFieldsData(
	gameUUID uuid.UUID, pl *dto.PlayData,
) (string, *[]dto.FieldData, error) {
	var obtainedWord string
	var fieldsData []dto.FieldData

	fields, err := gc.fieldService.GetWithGameUUID(gameUUID)
	if err != nil {
		return obtainedWord, &fieldsData, err
	}

	straightAxisId := gc.whichAxisIsStraight(&pl.Chars)
	nonStraightAxisId := straightAxisId ^ 1

	charsPositionsMinHeap := &common.PriorityQueue[*dto.Char, int]{}
	heap.Init(charsPositionsMinHeap)
	for i := range pl.Chars {
		heap.Push(
			charsPositionsMinHeap,
			&common.PQItem[*dto.Char, int]{
				Content:  &pl.Chars[i],
				Priority: pl.Chars[i].Position[nonStraightAxisId],
			},
		)
	}

	firstCharInHeap := charsPositionsMinHeap.GetTop().Content
	lastPosition := firstCharInHeap.Position[nonStraightAxisId]
	straightAxisNumber := firstCharInHeap.Position[straightAxisId]
	sideInt := common.Abs(pl.SideInt) % 360

	for _, field := range *fields {
		log.Println(field)
	}
	log.Println("straightAxisNumber")
	log.Println(straightAxisNumber)
	log.Println("straightAxisId")
	log.Println(straightAxisId)
	log.Println("sideInt")
	log.Println(sideInt)

	existingChars, existingCharsDepth := gc.makeCharsInStraightAxisFromFieldsMap(
		fields,
		straightAxisNumber,
		straightAxisId,
		sideInt,
	)
	log.Println("* existingChars")
	log.Println(*existingChars)

	log.Println("* existingCharsDepth")
	log.Println(*existingCharsDepth)

	if len(*existingChars) == 0 || len(*existingCharsDepth) == 0 {
		err := errors.New("no existing characters on straight axis")
		return obtainedWord, &fieldsData, err
	}

	depthLevel := gc.findDepthFromLeftMostPosition(existingCharsDepth)
	if depthLevel == -1 {
		err := errors.New("no depth level found")
		return obtainedWord, &fieldsData, err
	}

	// Boolean indicating if any char is near existing fields
	isTouching := false

	if val, ok := (*existingChars)[lastPosition-1]; ok {
		isTouching = true
		obtainedWord += val
	}

	for charsPositionsMinHeap.Len() > 0 {
		// get char
		pqi := heap.Pop(charsPositionsMinHeap).(*common.PQItem[*dto.Char, int])
		var char *dto.Char = pqi.Content
		for { // search fields in gaps from input WHILE char isn't in next pos
			if char.Position[nonStraightAxisId] == lastPosition {
				break
			}
			val, ok := (*existingChars)[lastPosition]
			if !ok {
				err := fmt.Errorf(
					"gap on position changingAxis: %v straightAxis: %v",
					lastPosition,
					straightAxisNumber,
				)
				return obtainedWord, &fieldsData, err
			}
			obtainedWord += val
			// Change depth used to create fieldData
			currDepth, ok := (*existingCharsDepth)[lastPosition]
			if !ok {
				err := fmt.Errorf("no depth in position %v", lastPosition)
				return obtainedWord, &fieldsData, err
			}
			depthLevel = currDepth
			isTouching = true
			lastPosition += 1
		}
		// Check if char does not colide with any field
		if _, ok := (*existingChars)[char.Position[nonStraightAxisId]]; ok {
			err := fmt.Errorf("char %v colide with field", char)
			return obtainedWord, &fieldsData, err
		}
		obtainedWord += char.Value
		lastPosition = char.Position[nonStraightAxisId] + 1
		fieldData := gc.createFieldData(
			char,
			sideInt,
			depthLevel,
			straightAxisId,
		)
		fieldsData = append(fieldsData, *fieldData)
	}

	if val, ok := (*existingChars)[lastPosition]; ok {
		isTouching = true
		obtainedWord += val
	}

	if !isTouching {
		err := errors.New("chars does not touch any field")
		return obtainedWord, &fieldsData, err
	}

	return obtainedWord, &fieldsData, nil
}

// Check if characters are in available characters
func (gc *gameController) areCharsInAvChars(
	playerUUID uuid.UUID,
	chars *[]dto.Char,
) error {
	avChars, err := gc.avCharService.GetWithPlayerUUID(playerUUID)
	if err != nil {
		return err
	}

	for _, char := range *chars {
		charID, err := char.ParseID()
		if err != nil {
			return err
		}
		isInAvChars := false
		for _, avChar := range *avChars {
			if avChar.ID == charID && avChar.Value == char.Value {
				isInAvChars = true
				break
			}
		}
		if !isInAvChars {
			return fmt.Errorf(
				"char %v is not in available characters",
				char.HtmlIdentifier,
			)
		}
	}

	return nil
}

func (gc *gameController) checkChars(
	playerUUID uuid.UUID,
	chars *[]dto.Char,
) error {
	charsLength := len(*chars)
	if charsLength == 0 {
		return errors.New("characters cannot be empty")
	}

	uniquePositions := [2]map[int]bool{
		make(map[int]bool),
		make(map[int]bool),
	}
	for _, char := range *chars {
		switch { // Check field value
		case len(char.Value) != 1:
			return fmt.Errorf(
				"field.value length should be 1, not %v",
				len(char.Value),
			)
		case !strings.Contains(cfg.ALLOWED_CHARACTERS, char.Value):
			return fmt.Errorf(
				"field.value ('%v') should be allowed character; allowed characters: %v",
				char.Value,
				cfg.ALLOWED_CHARACTERS,
			)
		}
		for i := range 2 {
			uniquePositions[i][char.Position[i]] = true
		}
	}

	uniqeCountY := len(uniquePositions[0])
	uniqeCountX := len(uniquePositions[1])
	// Ensure the fields are aligned and do not overlap. One dimension should
	// have a length of 1, and the other should match the length of the word.
	if !((uniqeCountY == 1 && uniqeCountX == charsLength) ||
		(uniqeCountX == 1 && uniqeCountY == charsLength)) {
		return fmt.Errorf(
			"fields overlap or are not aligned; Y: %v X: %v",
			uniqeCountY, uniqeCountX,
		)
	}

	if err := gc.areCharsInAvChars(playerUUID, chars); err != nil {
		return err
	}

	return nil
}

func (gc *gameController) buildHtmlFields(fields *[]model.Field) ([]byte, error) {
	var htmlContent bytes.Buffer
	tmpl, err := template.ParseFiles("views/game/field.html")
	if err != nil {
		return nil, err
	}

	for _, field := range *fields {
		data := dto.NewHtmlFieldData(
			field.Value,
			field.PosX,
			field.PosY,
			field.PosZ,
		)

		err = tmpl.Execute(&htmlContent, data)
		if err != nil {
			return nil, err
		}
	}

	return htmlContent.Bytes(), nil
}

func (gc *gameController) buildHtmlAvChars(avChars *[]model.AvChar) ([]byte, error) {
	var htmlContent bytes.Buffer

	tmpl, err := template.ParseFiles("views/game/avchar-outer.html")
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&htmlContent, nil)
	if err != nil {
		return nil, err
	}

	tmpl, err = template.ParseFiles("views/game/avchar-inner.html")
	if err != nil {
		return nil, err
	}

	for _, avChar := range *avChars {
		data := dto.HtmlAvCharData{
			ID:    int(avChar.ID),
			Value: avChar.Value,
		}

		_, err := htmlContent.Write([]byte("\n"))
		if err != nil {
			return nil, err
		}

		err = tmpl.Execute(&htmlContent, data)
		if err != nil {
			return nil, err
		}
	}

	result := append(htmlContent.Bytes(), []byte("\n</div>")...)
	return result, nil
}

// |PUBLIC| //

func (gc *gameController) GetCurrentFields(ctx *dto.WsContext) ([]byte, error) {
	fields, err := gc.fieldService.GetWithGameUUID(ctx.Game.UUID)
	if err != nil {
		return nil, err
	}

	response, err := gc.buildHtmlFields(fields)
	return response, err
}

func (gc *gameController) GetAvaibleChars(ctx *dto.WsContext) ([]byte, error) {
	avChars, err := gc.avCharService.GetWithPlayerUUID(ctx.Player.UUID)
	if err != nil {
		return nil, err
	}

	if len(*avChars) < cfg.AVAILABLE_CHARACTERS_NUMBER {
		n := cfg.AVAILABLE_CHARACTERS_NUMBER - len(*avChars)
		createdAvChars, err := gc.avCharService.CreateMany(ctx.Player, n)
		if err != nil {
			return nil, err
		}
		r := append(*avChars, *createdAvChars...)
		avChars = &r
	}

	response, err := gc.buildHtmlAvChars(avChars)
	return response, err
}

func (gc *gameController) RemoveChars(ctx *dto.WsContext, chars *[]dto.Char) error {
	charsIDs := make([]int64, 0)
	for _, char := range *chars {
		i, err := char.ParseID()
		if err != nil {
			return err
		}
		charsIDs = append(charsIDs, i)
	}

	err := gc.avCharService.DeleteMany(&charsIDs)
	if err != nil {
		return err
	}

	return nil
}

func (gc *gameController) ReceiveChars(
	ctx *dto.WsContext, playData *dto.PlayData,
) ([]byte, []byte, error) {
	err := gc.checkChars(ctx.Player.UUID, &playData.Chars)
	if err != nil {
		return nil, nil, err
	}

	word, fieldsData, err := gc.obtainWordAndFieldsData(ctx.Game.UUID, playData)
	if err != nil {
		return nil, nil, err
	}

	err = gc.wordsController.CheckWord(word)
	if err != nil {
		return nil, nil, err
	}

	newFields, err := gc.fieldService.CreateMany(
		ctx.Game.UUID,
		ctx.Player.UUID,
		ctx.Player.Appends,
		fieldsData,
	)
	if err != nil {
		return nil, nil, err
	}

	ctx.Player.Appends += 1
	err = gc.playerService.Update(ctx.Player)
	if err != nil {
		return nil, nil, err
	}

	response, err := gc.buildHtmlFields(newFields)
	if err != nil {
		return nil, nil, err
	}

	err = gc.RemoveChars(ctx, &playData.Chars)
	if err != nil {
		return nil, nil, err
	}

	avCharsResponse, err := gc.GetAvaibleChars(ctx)
	if err != nil {
		return nil, nil, err
	}

	return response, avCharsResponse, nil
}
