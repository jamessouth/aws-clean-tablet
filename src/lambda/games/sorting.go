package main

import "sort"

func (players listPlayerList) sortByName() {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})
}

func (players livePlayerList) sortByName() {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})
}

func (players livePlayerList) sortByAnswerThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Answer.Answer != players[j].Answer.Answer:
			return players[i].Answer.Answer < players[j].Answer.Answer
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) sortByScoreThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Score != players[j].Score:
			return players[i].Score > players[j].Score
		default:
			return players[i].Name < players[j].Name
		}
	})
}

// type lessFuncList func(p1, p2 *listPlayer) int

// type multiSorterList struct {
// 	players listPlayerList
// 	less    []lessFuncList
// }

// func (ms *multiSorterList) Sort(players listPlayerList) {
// 	ms.players = players
// 	sort.Sort(ms)
// }

// func OrderedByList(less ...lessFuncList) *multiSorterList {
// 	return &multiSorterList{
// 		less: less,
// 	}
// }

// func (ms *multiSorterList) Len() int {
// 	return len(ms.players)
// }

// func (ms *multiSorterList) Swap(i, j int) {
// 	ms.players[i], ms.players[j] = ms.players[j], ms.players[i]
// }

// func (ms *multiSorterList) Less(i, j int) bool {
// 	for _, k := range ms.less {
// 		switch k(&ms.players[i], &ms.players[j]) {
// 		case 1:
// 			return true
// 		case -1:
// 			return false
// 		}
// 	}

// 	return true
// }

// func (p listPlayerList) sort(fs ...lessFuncList) listPlayerList {
// 	OrderedByList(fs...).Sort(p)

// 	return p
// }

// var namesList = func(a, b *listPlayer) int {
// 	if a.Name > b.Name {
// 		return -1
// 	}

// 	return 1
// }

// type lessFuncLive func(p1, p2 *livePlayer) int

// type multiSorterLive struct {
// 	players livePlayerList
// 	less    []lessFuncLive
// }

// func (ms *multiSorterLive) Sort(players livePlayerList) {
// 	ms.players = players
// 	sort.Sort(ms)
// }

// func OrderedByLive(less ...lessFuncLive) *multiSorterLive {
// 	return &multiSorterLive{
// 		less: less,
// 	}
// }

// func (ms *multiSorterLive) Len() int {
// 	return len(ms.players)
// }

// func (ms *multiSorterLive) Swap(i, j int) {
// 	ms.players[i], ms.players[j] = ms.players[j], ms.players[i]
// }

// func (ms *multiSorterLive) Less(i, j int) bool {
// 	for _, k := range ms.less {
// 		switch k(&ms.players[i], &ms.players[j]) {
// 		case 1:
// 			return true
// 		case -1:
// 			return false
// 		}
// 	}

// 	return true
// }

// func (p livePlayerList) sort(fs ...lessFuncLive) livePlayerList {
// 	OrderedByLive(fs...).Sort(p)

// 	return p
// }

// var namesLive = func(a, b *livePlayer) int {
// 	if a.Name > b.Name {
// 		return -1
// 	}

// 	return 1
// }

// var scores = func(a, b *livePlayer) int {
// 	if a.Score < b.Score {
// 		return -1
// 	}
// 	if a.Score > b.Score {
// 		return 1
// 	}

// 	return 0
// }

// var answers = func(a, b *livePlayer) int {
// 	if a.Answer.Answer > b.Answer.Answer {
// 		return -1
// 	}

// 	return 1
// }
