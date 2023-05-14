package scene

type Item struct {
	Id   string
	Name string

	sceneId string
	showId  string
}

type Items []*Item

//func (i Item) TestFunction() *Item {
//	fmt.Printf("streamkit/show/scenes\n")
//	return i
//}
