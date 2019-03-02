package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Rooms struct {
	Rooms []*Room `json:"rooms"`
}

type Room struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	North   int      `json:"north"`
	South   int      `json:"south"`
	West    int      `json:"west"`
	East    int      `json:"east"`
	Objects []Object `json:"objects"`
}

type Object struct {
	Name string `json:"name"`
}

func main() {
	var (
		jRooms   = flag.String("map", "", "The map of the rooms in JSON format with 'rooms' as root item")
		start    = flag.Int("start", 0, "Starting room ID")
		sObjects = flag.String("objects", "", "The list of objects to collect separated by comma")
	)
	flag.Parse()

	fmt.Println()
	fmt.Println()
	fmt.Println(`     \            \  |                          _)                |        
    _ \          |\/ |   _` + "`" + ` | _  /   _ \         |  __ \    _` + "`" + ` |  |  |   | 
   ___ \ _____|  |   |  (   |   /    __/ _____|  |  |   |  (   |  |  |   | 
 _/    _\       _|  _| \__,_| ___| \___|        _| _|  _| \__, | _| \__, | 
                                                          |___/     ____/  `)
	fmt.Println()
	fmt.Println("--- Retro Route Puzzle ---")
	fmt.Println()
	fmt.Println()
	fmt.Println("- Listening on port :9090 -")
	fmt.Println()
	fmt.Println()

	go func() {
		mux := httprouter.New()
		mux.POST("/rooms", handle)
		log.Fatal(http.ListenAndServe(":9090", mux))
	}()

	for {
		var mapPath string
		if *jRooms == "" {
			fmt.Print("Enter path to rooms map JSON file: ")
			_, err := fmt.Scanf("%s\n", &mapPath)
			if err != nil {
				panic(err)
			}
		}
		if *start == 0 {
			fmt.Print("Enter starting room ID: ")
			_, err := fmt.Scanf("%d\n", start)
			if err != nil {
				panic(err)
			}
		}
		if *sObjects == "" {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter objects to collect separated by comma: ")
			for scanner.Scan() {
				*sObjects = scanner.Text()
				break
			}
		}

		var rooms Rooms
		objects := strings.Split(*sObjects, " ")

		if *jRooms != "" {
			err := json.Unmarshal([]byte(*jRooms), &rooms)
			if err != nil {
				panic(err)
			}
		}
		if mapPath != "" {
			raw, err := ioutil.ReadFile(mapPath)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(raw, &rooms)
			if err != nil {
				panic(err)
			}
		}

		roomsMap, steps := traverse(&rooms, *start, objects)

		fmt.Println("ID\tRoom\t\tObject collected")
		fmt.Println("_________________________________________")
		for _, s := range *steps {
			var objNames []string
			if len(roomsMap[s].Objects) > 0 {
				for _, n := range roomsMap[s].Objects {
					objNames = append(objNames, n.Name)
				}
			} else {
				objNames = append(objNames, "None")
			}
			fmt.Printf("%v\t%-10s\t%s\n", roomsMap[s].ID, roomsMap[s].Name, objNames)
		}
		fmt.Println()
		fmt.Println()
		*start = 0
		*sObjects = ""
	}
}

func traverse(rooms *Rooms, start int, objects []string) (map[int]*Room, *[]int) {
	roomsMap := make(map[int]*Room)
	visited := make(map[int]bool)

	for _, r := range rooms.Rooms {
		roomsMap[r.ID] = r
		visited[r.ID] = false
	}

	curr := *roomsMap[start]
	var steps []int
	var next *Room

	for {
		next = nil
		edges := getEdges(&curr)
		visited[curr.ID] = true
		steps = append(steps, curr.ID)

		if len(curr.Objects) > 0 {
			for i, o := range objects {
				for _, obj := range curr.Objects {
					if strings.TrimSpace(o) == obj.Name {
						objects = append(objects[:i], objects[i+1:]...)
					}
				}
			}
		}

		if len(objects) == 0 {
			break
		}

		end := true

		next = getNext(edges, visited, roomsMap)
		if next == nil {
			for _, v := range visited {
				if v == false {
					end = false
				}
			}
			for _, e := range edges {
				if e != 0 && e != steps[len(steps)-2] {
					next = roomsMap[e]
				}
			}
			if next == nil {
				next = roomsMap[steps[len(steps)-2]]
			}
		} else {
			end = false
		}

		if end == true {
			break
		}
		curr = *next
	}
	return roomsMap, &steps
}

func getEdges(r *Room) []int {
	return []int{r.North, r.South, r.East, r.West}
}

func getNext(edges []int, visited map[int]bool, rm map[int]*Room) *Room {
	for _, e := range edges {
		if e != 0 {
			vOk := visited[e]
			if !vOk {
				return rm[e]
			}
		}
	}
	return nil
}

func handle(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var r struct {
		Rooms   Rooms    `json:"rooms"`
		Start   int      `json:"start"`
		Objects []string `json:"objects"`
	}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		log.Println(err)
	}

	rm, steps := traverse(&r.Rooms, r.Start, r.Objects)

	type step struct {
		ID      int      `json:"ID"`
		Name    string   `json:"name"`
		Objects []string `json:"objects"`
	}

	var res struct {
		Steps []step `json:"steps"`
	}

	for _, s := range *steps {
		var objNames []string
		if len(rm[s].Objects) > 0 {
			for _, n := range rm[s].Objects {
				objNames = append(objNames, n.Name)
			}
		} else {
			objNames = append(objNames, "None")
		}
		res.Steps = append(res.Steps, step{rm[s].ID, rm[s].Name, objNames})
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}
