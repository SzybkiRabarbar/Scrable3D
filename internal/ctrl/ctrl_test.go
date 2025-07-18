package ctrl

import (
	"os"
	"path/filepath"
	"reflect"
	"scrable3/internal/dto"
	"scrable3/internal/mock"
	"scrable3/internal/model"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func setupGameControllerImplementation(t *testing.T) *gameController {
	mockController := gomock.NewController(t)
	return &gameController{
		mock.NewMockWordsController(mockController),
		mock.NewMockPlayerService(mockController),
		mock.NewMockFieldService(mockController),
		mock.NewMockAvCharService(mockController),
	}
}

func TestMakeExistingCharsInStraightAxisMap(t *testing.T) {
	gc := setupGameControllerImplementation(t)
	type mockedData struct {
		fields             []model.Field
		straightAxisNumber int
		straightAxisId     int
		sideInt            int
		idealOutput        map[int]string
	}
	datasets := []mockedData{
		{ // * 0 simple case *
			fields: []model.Field{
				{Value: "Q", PosX: 1, PosY: 5, PosZ: 6},
				{Value: "W", PosX: 2, PosY: 5, PosZ: 6},
				{Value: "E", PosX: 3, PosY: 5, PosZ: 6},
			},
			straightAxisNumber: 5,
			straightAxisId:     0,
			sideInt:            0,
			idealOutput: map[int]string{
				1: "Q",
				2: "W",
				3: "E",
			},
		},
		{ // * 1 gap between values *
			fields: []model.Field{
				{Value: "Q", PosX: 1, PosY: 5, PosZ: 6},
				{Value: "W", PosX: 2, PosY: 5, PosZ: 6},
				{Value: "E", PosX: 14, PosY: 5, PosZ: 6},
			},
			straightAxisNumber: 5,
			straightAxisId:     0,
			sideInt:            0,
			idealOutput: map[int]string{
				1:  "Q",
				2:  "W",
				14: "E",
			},
		},
		{ // * 2 one Pos not in straigth *
			fields: []model.Field{
				{Value: "Q", PosX: 1, PosY: 5, PosZ: 6},
				{Value: "W", PosX: 2, PosY: 5, PosZ: 6},
				{Value: "E", PosX: 3, PosY: 5, PosZ: 6},
				{Value: "R", PosX: 3, PosY: 6, PosZ: 6},
				{Value: "T", PosX: 5, PosY: 5, PosZ: 6},
			},
			straightAxisNumber: 5,
			straightAxisId:     0,
			sideInt:            0,
			idealOutput: map[int]string{
				1: "Q",
				2: "W",
				3: "E",
				5: "T",
			},
		},
		{ // * 3 diffrences in depth *
			fields: []model.Field{
				{Value: "Q", PosX: 1, PosY: 5, PosZ: 6},
				{Value: "W", PosX: 2, PosY: 5, PosZ: 6},
				{Value: "E", PosX: 3, PosY: 5, PosZ: 6},
				{Value: "R", PosX: 3, PosY: 5, PosZ: 5},
				{Value: "T", PosX: 5, PosY: 5, PosZ: 6},
				{Value: "Y", PosX: 5, PosY: 5, PosZ: 7},
			},
			straightAxisNumber: 5,
			straightAxisId:     0,
			sideInt:            0,
			idealOutput: map[int]string{
				1: "Q",
				2: "W",
				3: "R",
				5: "T",
			},
		},
		{ // * 4 Y not straigth and diff in depth *
			fields: []model.Field{
				{Value: "Q", PosX: 3, PosY: 4, PosZ: 3},
				{Value: "W", PosX: 3, PosY: 5, PosZ: 2},
				{Value: "E", PosX: 3, PosY: 6, PosZ: 1},
			},
			straightAxisNumber: 3,
			straightAxisId:     1,
			sideInt:            0,
			idealOutput: map[int]string{
				4: "Q",
				5: "W",
				6: "E",
			},
		},
		{ // * 5 Y not straigth and diff side *
			fields: []model.Field{
				{Value: "Q", PosX: 3, PosY: 4, PosZ: 3},
				{Value: "W", PosX: 3, PosY: 5, PosZ: 2},
				{Value: "E", PosX: 3, PosY: 6, PosZ: 1},
			},
			straightAxisNumber: 3,
			straightAxisId:     1,
			sideInt:            180,
			idealOutput: map[int]string{
				4: "Q",
				5: "W",
				6: "E",
			},
		},
		{ // * 6 check reverse depth *
			fields: []model.Field{
				{Value: "Q", PosX: 4, PosY: 6, PosZ: 3},
				{Value: "W", PosX: 4, PosY: 6, PosZ: 2},
				{Value: "E", PosX: 4, PosY: 6, PosZ: 1},
			},
			straightAxisNumber: 6,
			straightAxisId:     0,
			sideInt:            180,
			idealOutput: map[int]string{
				4: "Q",
			},
		},
		{ // * 7 empty straight axis *
			fields: []model.Field{
				{Value: "Q", PosX: 4, PosY: 6, PosZ: 3},
				{Value: "W", PosX: 4, PosY: 6, PosZ: 2},
				{Value: "E", PosX: 4, PosY: 6, PosZ: 1},
			},
			straightAxisNumber: 4,
			straightAxisId:     0,
			sideInt:            180,
			idealOutput:        map[int]string{},
		},
		{ // * 8 biggest X as depth *
			fields: []model.Field{
				{Value: "Q", PosX: 4, PosY: 12, PosZ: 3},
				{Value: "W", PosX: 4, PosY: 12, PosZ: 2},
				{Value: "E", PosX: 4, PosY: 12, PosZ: 1},
				{Value: "R", PosX: 5, PosY: 12, PosZ: 1},
				{Value: "T", PosX: 6, PosY: 12, PosZ: 1},
				{Value: "Y", PosX: 1, PosY: 12, PosZ: 5},
				{Value: "U", PosX: 2, PosY: 12, PosZ: 5},
				{Value: "I", PosX: 15, PosY: 12, PosZ: 5},
				{Value: "O", PosX: 14, PosY: 12, PosZ: 5},
			},
			straightAxisNumber: 12,
			straightAxisId:     0,
			sideInt:            90,
			idealOutput: map[int]string{
				3: "Q",
				2: "W",
				1: "T",
				5: "I",
			},
		},
		{ // * 9 smallest X as depth *
			fields: []model.Field{
				{Value: "Q", PosX: 4, PosY: 12, PosZ: 3},
				{Value: "W", PosX: 4, PosY: 12, PosZ: 2},
				{Value: "E", PosX: 4, PosY: 12, PosZ: 1},
				{Value: "R", PosX: 5, PosY: 12, PosZ: 1},
				{Value: "T", PosX: 6, PosY: 12, PosZ: 1},
				{Value: "Y", PosX: 1, PosY: 12, PosZ: 5},
				{Value: "U", PosX: 2, PosY: 12, PosZ: 5},
				{Value: "I", PosX: 15, PosY: 12, PosZ: 5},
				{Value: "O", PosX: 14, PosY: 12, PosZ: 5},
			},
			straightAxisNumber: 12,
			straightAxisId:     0,
			sideInt:            270,
			idealOutput: map[int]string{
				3: "Q",
				2: "W",
				1: "E",
				5: "Y",
			},
		},
		{ // * 10 *
			fields: []model.Field{
				{Value: "Q", PosX: 6, PosY: 1, PosZ: 3},
				{Value: "W", PosX: 3, PosY: 2, PosZ: 1},
				{Value: "E", PosX: 9, PosY: 3, PosZ: 1},
				{Value: "R", PosX: 3, PosY: 4, PosZ: 1},
				{Value: "T", PosX: 6, PosY: 4, PosZ: 1},
				{Value: "Y", PosX: 10, PosY: 7, PosZ: 5},
				{Value: "U", PosX: 9, PosY: 8, PosZ: 1},
				{Value: "I", PosX: 15, PosY: 8, PosZ: 1},
				{Value: "O", PosX: 4, PosY: 9, PosZ: 5},
			},
			straightAxisNumber: 1,
			straightAxisId:     1,
			sideInt:            90,
			idealOutput: map[int]string{
				2: "W",
				3: "E",
				4: "T",
				8: "I",
			},
		},
		{ // * 11 *
			fields: []model.Field{
				{Value: "Q", PosX: 6, PosY: 1, PosZ: 3},
				{Value: "W", PosX: 3, PosY: 2, PosZ: 1},
				{Value: "E", PosX: 9, PosY: 3, PosZ: 1},
				{Value: "R", PosX: 3, PosY: 4, PosZ: 1},
				{Value: "T", PosX: 6, PosY: 4, PosZ: 1},
				{Value: "Y", PosX: 10, PosY: 7, PosZ: 5},
				{Value: "U", PosX: 9, PosY: 8, PosZ: 1},
				{Value: "I", PosX: 15, PosY: 8, PosZ: 1},
				{Value: "O", PosX: 4, PosY: 9, PosZ: 5},
			},
			straightAxisNumber: 1,
			straightAxisId:     1,
			sideInt:            270,
			idealOutput: map[int]string{
				2: "W",
				3: "E",
				4: "R",
				8: "U",
			},
		},
		{ // * 12 *
			fields: []model.Field{
				{Value: "Q", PosX: 6, PosY: 1, PosZ: 3},
				{Value: "W", PosX: 3, PosY: 2, PosZ: 1},
				{Value: "E", PosX: 9, PosY: 3, PosZ: 1},
				{Value: "R", PosX: 3, PosY: 4, PosZ: 1},
				{Value: "T", PosX: 6, PosY: 4, PosZ: 1},
				{Value: "Y", PosX: 10, PosY: 7, PosZ: 5},
				{Value: "U", PosX: 9, PosY: 8, PosZ: 1},
				{Value: "I", PosX: 15, PosY: 8, PosZ: 1},
				{Value: "O", PosX: 4, PosY: 9, PosZ: 5},
			},
			straightAxisNumber: 11,
			straightAxisId:     0,
			sideInt:            270,
			idealOutput:        map[int]string{},
		},
	}
	for i, data := range datasets[len(datasets)-1:] {
		existingChars, existingCharsDetph := gc.makeCharsInStraightAxisFromFieldsMap(
			&data.fields,
			data.straightAxisNumber,
			data.straightAxisId,
			data.sideInt,
		)
		if !reflect.DeepEqual(existingChars, &data.idealOutput) {
			t.Errorf("%v. %v %v", i, existingChars, &data.idealOutput)
		}
		if len(*existingChars) != len(*existingCharsDetph) {
			t.Errorf(
				"%v. diffrence in len; %v %v",
				i, existingChars, existingCharsDetph,
			)
		}
		for key := range *existingChars {
			if _, exists := (*existingCharsDetph)[key]; !exists {
				t.Errorf(
					"%v. diffrence in keys; %v %v",
					i, existingChars, existingCharsDetph,
				)
			}
		}
	}
}

func TestFindDepthFromLeftMostPosition(t *testing.T) {
	gc := setupGameControllerImplementation(t)
	type mockedData struct {
		data        map[int]int
		idealOutput int
	}
	datasets := []mockedData{
		{
			data:        map[int]int{1: 1, 2: 2},
			idealOutput: 1,
		},
		{
			data:        map[int]int{2: 1, 1: 2},
			idealOutput: 2,
		},
		{
			data:        map[int]int{3: 9, 6: 4, 8: 2},
			idealOutput: 9,
		},
		{
			data:        map[int]int{3: 9, 6: 4, 8: 2, 2: 0},
			idealOutput: 0,
		},
	}
	for i, dataset := range datasets {
		depth := gc.findDepthFromLeftMostPosition(&dataset.data)
		if depth != dataset.idealOutput {
			t.Errorf("%v. wrong output; Returned: %v; IdealOutput: %v",
				i, depth, dataset.idealOutput,
			)
		}
	}
}

func TestCreateFieldData(t *testing.T) {
	gc := setupGameControllerImplementation(t)
	type mockedData struct {
		char        dto.Char
		sideInt     int
		depthLevel  int
		idealOutput dto.FieldData
	}
	datasets := []mockedData{
		{
			char:        dto.Char{Value: "Q", Position: [2]int{2, 1}},
			sideInt:     0,
			depthLevel:  3,
			idealOutput: dto.FieldData{Value: "Q", Pos: [3]int{1, 2, 3}},
		},
		{
			char:        dto.Char{Value: "Q", Position: [2]int{2, 1}},
			sideInt:     90,
			depthLevel:  3,
			idealOutput: dto.FieldData{Value: "Q", Pos: [3]int{3, 2, 1}},
		},
		{
			char:        dto.Char{Value: "Q", Position: [2]int{2, 1}},
			sideInt:     180,
			depthLevel:  3,
			idealOutput: dto.FieldData{Value: "Q", Pos: [3]int{13, 2, 3}},
		},
		{
			char:        dto.Char{Value: "Q", Position: [2]int{2, 1}},
			sideInt:     270,
			depthLevel:  3,
			idealOutput: dto.FieldData{Value: "Q", Pos: [3]int{3, 2, 13}},
		},
	}
	for i, dataset := range datasets {
		fieldData := gc.createFieldData(
			&dataset.char,
			dataset.sideInt,
			dataset.depthLevel,
		)
		if !reflect.DeepEqual(fieldData, &dataset.idealOutput) {
			t.Errorf(
				"%v. wrong output; Returns: %v; IdealOutput: %v",
				i, fieldData, &dataset.idealOutput,
			)
		}
	}
}

func TestValidCheckChars(t *testing.T) {
	gc := setupGameControllerImplementation(t)
	mockAvChars := &[]model.AvChar{
		{ID: 1, Value: "Q"}, {ID: 2, Value: "W"}, {ID: 3, Value: "E"},
		{ID: 4, Value: "R"}, {ID: 5, Value: "T"}, {ID: 6, Value: "Y"},
		{ID: 7, Value: "U"}, {ID: 8, Value: "I"}, {ID: 9, Value: "O"},
		{ID: 10, Value: "P"}, {ID: 11, Value: "A"}, {ID: 12, Value: "S"},
		{ID: 13, Value: "D"}, {ID: 14, Value: "F"}, {ID: 15, Value: "G"},
		{ID: 16, Value: "H"}, {ID: 17, Value: "J"},
	}
	gc.avCharService.(*mock.MockAvCharService).
		EXPECT().
		GetWithPlayerUUID(uuid.UUID{}).
		Return(mockAvChars, nil).
		AnyTimes()

	validChars := [][]dto.Char{
		{
			{Value: "Q", Position: [2]int{3, 6}, HtmlIdentifier: "char-Q1"},
			{Value: "W", Position: [2]int{3, 7}, HtmlIdentifier: "char-W2"},
			{Value: "E", Position: [2]int{3, 9}, HtmlIdentifier: "char-E3"},
		},
		{
			{Value: "R", Position: [2]int{8, 1}, HtmlIdentifier: "char-R4"},
			{Value: "T", Position: [2]int{8, 2}, HtmlIdentifier: "char-T5"},
			{Value: "Y", Position: [2]int{8, 4}, HtmlIdentifier: "char-Y6"},
			{Value: "U", Position: [2]int{8, 5}, HtmlIdentifier: "char-U7"},
		},
		{
			{Value: "I", Position: [2]int{6, 5}, HtmlIdentifier: "char-I8"},
			{Value: "O", Position: [2]int{7, 5}, HtmlIdentifier: "char-O9"},
			{Value: "P", Position: [2]int{9, 5}, HtmlIdentifier: "char-P10"},
			{Value: "A", Position: [2]int{10, 5}, HtmlIdentifier: "char-A11"},
		},
		{
			{Value: "S", Position: [2]int{2, 14}, HtmlIdentifier: "char-S12"},
			{Value: "D", Position: [2]int{3, 14}, HtmlIdentifier: "char-D13"},
			{Value: "F", Position: [2]int{5, 14}, HtmlIdentifier: "char-F14"},
			{Value: "G", Position: [2]int{7, 14}, HtmlIdentifier: "char-G15"},
			{Value: "H", Position: [2]int{8, 14}, HtmlIdentifier: "char-H16"},
			{Value: "J", Position: [2]int{9, 14}, HtmlIdentifier: "char-J17"},
		},
	}
	for i, chars := range validChars {
		err := gc.checkChars(uuid.UUID{}, &chars)
		if err != nil {
			t.Errorf("Check() number %v failed; %v", i+1, err)
		}
	}
}

func TestInvalidCheckChars(t *testing.T) {
	gc := setupGameControllerImplementation(t)

	mockAvChars := &[]model.AvChar{
		{ID: 1, Value: "Q"}, {ID: 2, Value: "W"}, {ID: 3, Value: "E"},
		{ID: 4, Value: "R"}, {ID: 5, Value: "T"}, {ID: 6, Value: "Y"},
		{ID: 7, Value: "U"}, {ID: 8, Value: "I"}, {ID: 9, Value: "O"},
		{ID: 10, Value: "P"}, {ID: 11, Value: "A"},
	}

	gc.avCharService.(*mock.MockAvCharService).
		EXPECT().
		GetWithPlayerUUID(uuid.UUID{}).
		Return(mockAvChars, nil).
		AnyTimes()

	testCases := []struct {
		name        string
		fields      []dto.Char
		expectedErr string
	}{
		{
			name:        "empty field slice",
			fields:      []dto.Char{},
			expectedErr: "characters cannot be empty",
		},
		{
			name: "value is wrong length",
			fields: []dto.Char{
				{Value: "Q", Position: [2]int{3, 6}, HtmlIdentifier: "char-Q1"},
				{Value: "WZ", Position: [2]int{3, 7}, HtmlIdentifier: "char-WZ2"},
				{Value: "E", Position: [2]int{3, 9}, HtmlIdentifier: "char-E3"},
			},
			expectedErr: "field.value length should be 1, not 2",
		},
		{
			name: "value is not uppercase",
			fields: []dto.Char{
				{Value: "R", Position: [2]int{8, 1}, HtmlIdentifier: "char-R4"},
				{Value: "t", Position: [2]int{8, 2}, HtmlIdentifier: "char-t5"},
				{Value: "Y", Position: [2]int{8, 4}, HtmlIdentifier: "char-Y6"},
				{Value: "U", Position: [2]int{8, 5}, HtmlIdentifier: "char-U7"},
			},
			expectedErr: "field.value ('t') should be allowed character; allowed characters: ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			name: "fields are not aligned",
			fields: []dto.Char{
				{Value: "I", Position: [2]int{6, 5}, HtmlIdentifier: "char-I8"},
				{Value: "O", Position: [2]int{7, 6}, HtmlIdentifier: "char-O9"},
				{Value: "P", Position: [2]int{9, 5}, HtmlIdentifier: "char-P10"},
				{Value: "A", Position: [2]int{10, 5}, HtmlIdentifier: "char-A11"},
			},
			expectedErr: "fields overlap or are not aligned; Y: 4 X: 2",
		},
		{
			name: "fields overlap",
			fields: []dto.Char{
				{Value: "S", Position: [2]int{2, 14}, HtmlIdentifier: "char-S12"},
				{Value: "D", Position: [2]int{3, 14}, HtmlIdentifier: "char-D13"},
				{Value: "F", Position: [2]int{5, 14}, HtmlIdentifier: "char-F14"},
				{Value: "G", Position: [2]int{7, 14}, HtmlIdentifier: "char-G15"},
				{Value: "H", Position: [2]int{8, 14}, HtmlIdentifier: "char-H16"},
				{Value: "J", Position: [2]int{8, 14}, HtmlIdentifier: "char-J17"},
			},
			expectedErr: "fields overlap or are not aligned; Y: 5 X: 1",
		},
		{
			name: "character not in available characters",
			fields: []dto.Char{
				{Value: "X", Position: [2]int{1, 1}, HtmlIdentifier: "char-X99"},
			},
			expectedErr: "char char-X99 is not in available characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := gc.checkChars(uuid.UUID{}, &tc.fields)
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if err.Error() != tc.expectedErr {
				t.Errorf("expected error '%s', got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestLoadWords(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("geting current working directory failed; %v", err)
	}
	filePath := filepath.Join(
		cwd, "..", "..", "words", "words_alpha.txt",
	)

	// Init wordsController
	wordsController, err := NewWordsController(filePath)
	if err != nil {
		t.Errorf("loading words failed; %v", err)
	}

	// Test if not empty
	if wordsController.WordsNumber() == 0 {
		t.Error("Empty words.Map")
	}
}

func TestCheckWord(t *testing.T) {
	realWords := []string{
		"breeze",
		"galaxy",
		"whisper",
		"echo",
		"mystery",
		"harmony",
		"journey",
		"twilight",
		"serenity",
		"radiance",
	}
	notRealWords := []string{
		"briize",
		"balaxy",
		"abdc",
		"sssss",
		"houseapple",
	}

	// Init mocked words
	wMap := make(map[string]bool)
	for _, word := range realWords {
		wMap[word] = true
	}
	wordsController := wordsController{
		words: &wMap,
	}

	// Check Words.CheckWord() with real words
	for _, word := range realWords {
		err := wordsController.CheckWord(word)
		if err != nil {
			t.Errorf("`%v` not checked as word; %v", word, err.Error())
		}
	}

	// Check Words.CheckWord() with not real words
	for _, notWord := range notRealWords {
		err := wordsController.CheckWord(notWord)
		if err == nil {
			t.Errorf("`%v` checked as word", notWord)
		}
	}
}
