package data

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Connections struct {
	Room  string
	Links []string
}

type Map struct {
	Name  string
	Start bool
	End   bool
	X, Y  string
}

type Ant struct {
	ID   int
	Path []string
	Room string
	X, Y string
}

func LoadLevel() ([]Map, []Connections, []Ant, string, string) {
	var maps []Map
	var connections []Connections
	var ants []Ant

	var StartPoint, EndPoint string
	lines := strings.Split(getLevel(), "\n")
	regexCoos := regexp.MustCompile(`^\w+ \d+ \d+$`)
	regexComment := regexp.MustCompile(`^#{1}\w+`)
	regexIllegalL := regexp.MustCompile(`^L`)
	regexConnect := regexp.MustCompile(`^\w+-\w+$`)

	for i := 0; i < len(lines); i++ {
		if i == 0 {
			intLines, _ := strconv.ParseInt((lines[i]), 10, 64)
			if intLines > 0 {
				for n := 0; n < int(intLines); n++ {
					ants = append(ants, Ant{ID: n, Room: "undefined", X: "0", Y: "0"})
				}
			} else {
				log.Fatal("Error : ants not initialized")
			}
		} else {
			if !regexComment.MatchString(lines[i]) {
				start := false
				end := false
				if lines[i] == "##start" {
					i++
					start = true
				} else if lines[i] == "##end" {
					i++
					end = true
				}
				if regexCoos.MatchString(lines[i]) {
					if regexIllegalL.MatchString(lines[i]) {
						log.Fatal("Room cannot start by \"L\" or \"#\".")
					} else {
						data := strings.Split(lines[i], " ")
						if CheckDupplicate(maps, data[0]) {
							log.Fatal("Error : Room name already taken.")
						} else {
							if start {
								StartPoint = data[0]
							} else if end {
								EndPoint = data[0]
							}
							maps = append(maps, Map{Name: data[0], X: data[1], Y: data[2], Start: start, End: end})
						}
					}
				} else if regexConnect.MatchString(lines[i]) {
					data := strings.Split(lines[i], "-")
					_, boolean := GetConnectionsByRoom(connections, data[1])
					if boolean {
						// room := GetConnectionsByRoom(data[0])
						for k := range connections {
							if connections[k].Room == data[1] {
								connections[k].Links = append(connections[k].Links, data[0])
							}
						}
					} else {
						connections = append(connections, Connections{data[1], []string{data[0]}})
					}
					_, boolean = GetConnectionsByRoom(connections, data[0])
					if boolean {
						// room := GetConnectionsByRoom(data[0])
						for k := range connections {
							if connections[k].Room == data[0] {
								connections[k].Links = append(connections[k].Links, data[1])
							}
						}
					} else {
						connections = append(connections, Connections{data[0], []string{data[1]}})
					}
				} else {
					log.Fatal("Error : \"" + lines[i] + "\" invalid data format")
				}
			}
		}
	}
	return maps, connections, ants, StartPoint, EndPoint
}

func getLevel() string {
	if len(os.Args) == 2 {
		file, err := os.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return string(file)
	} else {
		log.Fatal("Invalid number of argument, try \"go run . <file.txt>\".")
	}
	return "none"
}

func CheckDupplicate(maps []Map, name string) bool {
	for k := range maps {
		if maps[k].Name == name {
			return true
		}
	}
	return false
}

func GetConnectionsByRoom(connections []Connections, roomName string) (Connections, bool) {
	for k := range connections {
		if connections[k].Room == roomName {
			return connections[k], true
		}
	}
	return Connections{}, false
}
