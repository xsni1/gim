package editor

type KeyBinder interface {
	Get(key string) KeyEvent
}

func NewMapKeyBinder(kb ActionsMap) KeyBinder {
	binder := &MapKeyBinder{
		actionsMap:  kb,
		keybindings: defaultKeyBindings,
	}

	return binder
}

type MapKeyBinder struct {
	actionsMap  ActionsMap
	keybindings KeyBindings
}

func (b *MapKeyBinder) Get(key string) KeyEvent {
	actionName, ok := b.keybindings[key]
	if !ok {
		return func() {}
	}

	return b.actionsMap[actionName]
}

type ActionsMap map[string]KeyEvent

type KeyBindings map[string]string

type KeyEvent func()

var defaultKeyBindings = KeyBindings{
	"k":      "CursorUp",
	"j":      "CursorDown",
	"h":      "CursorLeft",
	"l":      "CursorRight",
	"i":      "InsertMode",
	"Ctrl+C": "NormalMode",
	"Esc":    "Quit",
	"Ctrl+S": "Save",
}
