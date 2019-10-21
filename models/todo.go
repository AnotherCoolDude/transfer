package models

// Todo is a generic type of todos
type Todo interface {
	Timestamp() string
	Identifier() int
	ClientType() string
}

// Todosorter sorts todos
type Todosorter []Todo

func (ts Todosorter) pairs() [][]Todo {
	types := []string{}
	typesmap := map[string][]Todo{}
	amountTypes := 0
	highestIndex := 0
	pairs := [][]Todo{}

	// identify and store types
	for _, todo := range ts {
		if _, ok := typesmap[todo.ClientType()]; !ok {
			amountTypes++
			typesmap[todo.ClientType()] = []Todo{}
		}
		typesmap[todo.ClientType()] = append(typesmap[todo.ClientType()], todo)
	}

	for clienttype := range typesmap {
		if len(typesmap[clienttype]) > highestIndex {
			highestIndex = len(typesmap[clienttype])
		}
	}

	for i := 0; i < highestIndex; i++ {

		for clienttype, tt := range typesmap {
			timestamp := ""
			if t, ok := typesmap[clienttype][i]; ok {
				if timestamp == "" {
					pairs[len(pairs)-1] = append(pairs[len(pairs)-1], t)
					timestamp = t.Timestamp()
				}

				if t.Timestamp() == timestamp {
					pairs[len(pairs)-1] = append(pairs[len(pairs)-1], t)
				}
			}
		}
	}
	return pairs
}
