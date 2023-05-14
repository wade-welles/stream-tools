package show

import (
	"fmt"
)

func TestFunction() {
	fmt.Println("module show")
}

type Production struct {
}

type Episode struct {
}

//type Show struct {
//	Name string
//}

type Config struct {
	Name string
}

//func (sh *Show) PrintDebug() {
//	fmt.Printf("show: \n")
//	fmt.Printf("  object: %v\n", sh)
//	// Breaks here because no scenes, its not even set to empty
//	fmt.Printf("  scenes: %v\n", sh.Scenes)
//	fmt.Printf("  scene_count: %v\n", len(sh.Scenes))
//	for _, scene := range sh.Scenes {
//		fmt.Printf("    scene:\n")
//		fmt.Printf("      name: %v\n", scene.Name)
//		fmt.Printf("      item_count: %v\n", len(scene.Items))
//		fmt.Printf("      items:\n")
//		for _, item := range scene.Items {
//			fmt.Printf("      item:\n")
//			fmt.Printf("        name: %v\n", item.Name)
//		}
//	}
//
//}
