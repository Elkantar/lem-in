package path

import (
	"fmt"
	"lemin/data"
	"strconv"
)

var maps, connections, ants, startPoint, endPoint = data.LoadLevel()

func pathFinder(start string, end string, banned []string) []string {
	// get firsts branches
	var bannedPoints []string
	// console.log(banned)
	pathSearch := [][]string{{start}}
	if len(banned) > 1 {
		bannedPoints = banned
	} else {
		bannedPoints = []string{start}
	}

	for i := 0; i < 10; i++ {
		// Read all path from pathSearch
		// console.log(pathSearch)
		for n := 0; n < len(pathSearch); n++ {
			// Save path for spliting
			path := pathSearch[n]
			actualRoom := pathSearch[n][len(pathSearch[n])-1]
			// Get all connections from last spot of path examined
			connects, _ := GetConnectionsByRoom(actualRoom)
			// console.log(connects)
			// read connections
			first := true
			for l := 0; l < len(connects.Links); l++ {
				if actualRoom == startPoint && connects.Links[l] == endPoint {
				} else {
					// console.log(connects.links[l])
					// Execute if point wasn't visited on another road yet
					if !Includes(bannedPoints, connects.Links[l]) {
						// update the road
						if first {
							pathSearch[n] = append(pathSearch[n], connects.Links[l])
							first = false
						} else {
							// create an alternative road if connections are more than one
							newPath := append(path, connects.Links[l])
							pathSearch = append(pathSearch, newPath)
						}
						// add point examined to banned points
						bannedPoints = append(bannedPoints, connects.Links[l])
						// if point's connections is the target point
						if connects.Links[l] == endPoint {
							result := append(path, connects.Links[l])
							// console.log("Found :")
							// console.log(result)
							return result
						}
					}
				}
			}
		}
	}
	return []string{}
}

func PathAnalizer() {
	var paths [][]string
	// Get connections for the start room
	roads, _ := GetConnectionsByRoom(startPoint)
	// Repeating the loop for each roads attached to start point
	for i := -1; i < len(roads.Links); i++ {
		if i == -1 {
			connects, _ := GetConnectionsByRoom(startPoint)
			var path []string
			for i := 0; i < len(connects.Links); i++ {
				if connects.Links[i] == endPoint {
					paths = append(paths, []string{startPoint, endPoint})
					break
				}
			}
			path = pathFinder(startPoint, endPoint, []string{})
			if len(path) == 0 {
				if !checkDupplicatePath(paths, path) {
					paths = append(paths, path)
				}
			}
		} else {
			// read map
			for k := 0; k < len(maps); k++ {
				// reinitialize banned points
				bannedPoints := []string{}
				// adding every room connected to start point except one
				for n := 0; n < len(roads.Links); n++ {
					if n != i {
						bannedPoints = append(bannedPoints, roads.Links[n])
					}
				}
				// if examined point is not start or end point
				if maps[k].Name != startPoint && maps[k].Name != endPoint {
					// adding examined point to banlist
					bannedPoints = append(bannedPoints, maps[k].Name)
					// searcch path
					path := pathFinder(startPoint, endPoint, bannedPoints)
					if len(path) != 0 {
						if !checkDupplicatePath(paths, path) {
							paths = append(paths, path)
						}
					}
				}
			}
		}
	}
	startConnect, _ := GetConnectionsByRoom(startPoint)
	endConnect, _ := GetConnectionsByRoom(endPoint)
	var roadsNumb int
	if len(startConnect.Links) < len(endConnect.Links) {
		roadsNumb = len(startConnect.Links)
	} else {
		roadsNumb = len(endConnect.Links)
	}
	var optimal [][]string
	if roadsNumb > 1 {
		routes := [][][]string{}
		for i := 0; i < len(paths); i++ {
			segment := [][]string{paths[i]}
			routes = append(routes, possibilitiesAnalizer(segment, roadsNumb, paths))
		}
		optimal = getOptimalRoute(routes)
	} else {
		optimal = [][]string{paths[0]}
	}
	startRoute(optimal)
}

// Personnal Hell
func possibilitiesAnalizer(ref [][]string, step int, corpus [][]string) [][]string {
	var result [][]string
	// read all path found
	for i := 0; i < len(corpus); i++ {
		// Check if path has rooms in common with rooms established
		if !checkCommonRooms(ref, corpus[i]) {
			result = append(ref, corpus[i])
			if step > 2 {
				result = possibilitiesAnalizer(result, step-1, corpus)
			}
		}
	}
	return result
}

func startRoute(routes [][]string) {
	// Display ants on start room
	count := 0
	var ratio []int
	greaterLen := 0
	smallestLen := 0
	for i := 0; i < len(routes); i++ {
		ratio = append(ratio, 1)
		if greaterLen < len(routes[i]) {
			greaterLen = len(routes[i])
		}
		if smallestLen > len(routes[i]) {
			count = i
			smallestLen = len(routes[i])
		}
	}
	for i := 0; i < len(routes); i++ {
		ratio = append(ratio, len(routes[i])-1)
	}
	countdown := true
	if len(routes) == 1 {
		countdown = false
	}

	for k := 0; k < len(ants); k++ {
		// new version
		if !countdown {
			ants[k].Path = routes[count]
			count++
			if count == len(routes) {
				count = 0
			}
		} else {
			ants[k].Path = routes[count]
			countdown = false
			allMax := true
			for i := 0; i < len(ratio); i++ {
				if ratio[i] < greaterLen-1 {
					allMax = false
					ants[k].Path = routes[i]
					ratio[i]++
					break
				}
			}
			if allMax {
				countdown = false
			}
		}
		ants[k].Room = startPoint
	}
	MoveAnts()
}

// count the number of room inside each group of parallels routes and keep the lowest
func getOptimalRoute(routes [][][]string) [][]string {
	var prevCount int
	var bestRoute [][]string
	first := true
	for i := 0; i < len(routes); i++ {
		if len(routes[i]) > 0 {
			count := 0
			for n := 0; n < len(routes[i]); n++ {
				count = count + len(routes[i][n])
			}
			if first {
				first = false
				prevCount = count
				bestRoute = routes[i]
			} else {
				if prevCount > count {
					prevCount = count
					bestRoute = routes[i]
				}
			}
		}
	}
	return bestRoute
}

var logs []string

func MoveAnts() {
	directEnd := false
	closeAnimation := false
	logsLine := ""
	allArrived := true
	for k := 0; k < len(ants); k++ {
		moveAnts := false
		if ants[k].Room != endPoint {
			nextRoom := getNextRoomInPath(ants[k].Room, ants[k].Path)
			if !roomFilled(nextRoom) || nextRoom == endPoint && ants[k].Room != startPoint {
				moveAnts = true
			}
			if ants[k].Room == startPoint && nextRoom == endPoint && !directEnd {
				directEnd = true
				moveAnts = true
			} else if directEnd && !moveAnts {
				moveAnts = false
			}
			if moveAnts {
				allArrived = false
				if logsLine == "" {
					logsLine = "L" + strconv.Itoa(ants[k].ID+1) + "-" + nextRoom
				} else {
					logsLine = logsLine + " L" + strconv.Itoa(ants[k].ID+1) + "-" + nextRoom
				}
				ants[k].Room = nextRoom
				nRoom, _ := GetRoomByName(nextRoom)
				ants[k].X = nRoom.X
				ants[k].Y = nRoom.Y
			}
		}
	}
	if logsLine != "" {
		logs = append(logs, logsLine)
	}
	directEnd = false
	if allArrived {
		var returnText string
		for i := 0; i < len(logs); i++ {
			if i == 0 {
				returnText = logs[i]
			} else {
				returnText = returnText + "\n" + logs[i]
			}
		}
		fmt.Println(returnText)
		closeAnimation = true
	}
	if !closeAnimation {
		MoveAnts()
	}
}

func GetConnectionsByRoom(roomName string) (data.Connections, bool) {
	for k := range connections {
		if connections[k].Room == roomName {
			return connections[k], true
		}
	}
	return data.Connections{}, false
}

func Includes(corpus []string, ref string) bool {
	for _, line := range corpus {
		if line == ref {
			return true
		}
	}
	return false
}

// Check if the path already exists insite a reference array
func checkDupplicatePath(ref [][]string, path []string) bool {
	for i := 0; i < len(ref); i++ {
		if len(ref[i]) == len(path) {
			dupplicate := false
			for n := 0; n < len(path); n++ {
				if ref[i][n] == path[n] {
					if n == len(path)-1 {
						dupplicate = true
					}
				} else {
					break
				}
			}
			if dupplicate {
				return true
			}
		}
	}
	return false
}

// Check an Array of path or a path and return true if there is common room inside
func checkCommonRooms(pathA [][]string, pathB []string) bool {
	for l := 0; l < len(pathA); l++ {
		for i := 1; i < len(pathA[l])-1; i++ {
			for n := 1; n < len(pathB)-1; n++ {
				if pathA[l][i] == pathB[n] {
					return true
				}
			}
		}
	}
	return false
}

func roomFilled(room string) bool {
	for k := 0; k < len(ants); k++ {
		if ants[k].Room == room {
			return true
		}
	}
	return false
}

func getNextRoomInPath(room string, path []string) string {
	for i := 0; i < len(path); i++ {
		if path[i] == room {
			return path[i+1]
		}
	}
	return ""
}

func GetRoomByName(roomName string) (data.Map, bool) {
	for k := range maps {
		if maps[k].Name == roomName {
			return maps[k], true
		}
	}
	return data.Map{}, false
}
