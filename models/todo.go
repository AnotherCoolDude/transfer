package models

import (
	"sort"
	"time"
)

// Todo is a generic type of todos
type Todo interface {
	Timestamp() string
	Identifier() int
	ClientType() string
}

// EmptyTodo returns an empty todo
type EmptyTodo struct{ ct string }

// Timestamp satisfies the Todo interface
func (et EmptyTodo) Timestamp() string {
	return "1990-01-01T01:01:01Z07:00"
}

// Identifier satisfies the Todo interface
func (et EmptyTodo) Identifier() int {
	return 0
}

// ClientType satisfies the Todo interface
func (et EmptyTodo) ClientType() string {
	return et.ct
}

// MatchOrEmpty returns a corresponding todo from todos, or an empty todo
func MatchOrEmpty(t Todo, todos []Todo) (Todo, bool) {
	for _, sliceTodo := range todos {
		ttimestamp, _ := time.Parse(time.RFC3339, t.Timestamp())
		sttimestamp, _ := time.Parse(time.RFC3339, sliceTodo.Timestamp())
		if ttimestamp.Equal(sttimestamp) {
			return sliceTodo, true
		}
	}
	return EmptyTodo{ct: todos[0].ClientType()}, false
}

// Todos sorts todos
type Todos []Todo

// ClienttypeSorting represents a Sorting by clienttype (keys) of todos
type ClienttypeSorting map[string]*Todos

// SortByClienttype sorts the todos by clienttype (mapkey) and timestamp (mapvalues)
func (tt *Todos) SortByClienttype() *ClienttypeSorting {
	typesmap := map[string]Todos{}
	sortedmap := ClienttypeSorting{}

	// identify and store types
	for _, todo := range *tt {
		if _, ok := typesmap[todo.ClientType()]; !ok {
			typesmap[todo.ClientType()] = []Todo{}
		}
		typesmap[todo.ClientType()] = append(typesmap[todo.ClientType()], todo)
	}

	keys := []string{}
	for key, tt := range typesmap {
		keys = append(keys, key)
		sort.Slice(tt, func(i, j int) bool {
			iTS, _ := time.Parse(time.RFC3339, tt[i].Timestamp())
			jTS, _ := time.Parse(time.RFC3339, tt[j].Timestamp())
			return iTS.Before(jTS)
		})
	}
	sort.Strings(keys)
	for _, k := range keys {
		*sortedmap[k] = typesmap[k]
	}
	return &sortedmap
}

// Clienttypes returns an array of strings of all clienttypes in clienttypesorting
func (cts *ClienttypeSorting) Clienttypes() []string {
	clienttypes := []string{}
	for ct := range *cts {
		clienttypes = append(clienttypes, ct)
	}
	return clienttypes
}

// SortedByTimestamp sorts the todos by their timestamp
func (tt *Todos) SortedByTimestamp() {
	sort.Slice(*tt, func(i, j int) bool {
		iTS, _ := time.Parse(time.RFC3339, (*tt)[i].Timestamp())
		jTS, _ := time.Parse(time.RFC3339, (*tt)[j].Timestamp())
		return iTS.Before(jTS)
	})
}

// SortedByCounterpart sorts tt to match cp using the timestamp and creates empty structs if no counterpart is available
func (tt *Todos) SortedByCounterpart(cp *Todos) {
	sorted := Todos{}
	unmatched := Todos{}
	for _, t := range *tt {
		for _, cpt := range *cp {
			tTS, _ := time.Parse(time.RFC3339, t.Timestamp())
			cptTS, _ := time.Parse(time.RFC3339, cpt.Timestamp())
			if tTS.Equal(cptTS) {
				sorted = append(sorted, t)
			} else {
				sorted = append(sorted, EmptyTodo{ct: t.ClientType()})
				unmatched = append(unmatched, t)
			}
		}
	}
	sorted = append(sorted, unmatched...)
	tt = &sorted
}

// Equivalent returns the equivalent and index from slice or an empty todo
func Equivalent(t Todo, slice []Todo) (Todo, int) {
	index := -1
	if len(slice) == 0 {
		return EmptyTodo{ct: "proad"}, index
	}
	ct := slice[0].ClientType()
	for i, st := range slice {
		if st.Timestamp() == t.Timestamp() {
			return st, i
		}
	}
	return EmptyTodo{ct: ct}, index
}

// MissingTodos returns all todos, which do not have a index in indexes
func MissingTodos(slice []Todo, indexes ...int) []Todo {
	mt := []Todo{}
	existing := false
	for i, st := range slice {
		for _, idx := range indexes {
			if i == idx {
				existing = true
			}
		}
		if !existing {
			mt = append(mt, st)
		}
		existing = false
	}
	return mt
}
