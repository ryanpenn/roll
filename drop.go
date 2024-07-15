package main

import (
	"math/rand"
	"time"
)

type (
	Drop struct {
		DropId int64 `json:"DropId"` // 掉落组
		Weight int   `json:"Weight"` // 组权重(万分比)
		Result int64 `json:"Result"` // 掉落结果
		IsEnd  int   `json:"IsEnd"`  // 是否为结束节点
	}

	DropGroup struct {
		DropId      int64   // 掉落组ID
		WeightCount int     // 组权重
		Drops       []*Drop // 掉落项
	}
)

const (
	FiveStarDropTimesLimit = 73  // 抽73次后增加概率
	FiveStarDropAddValue   = 600 // 超过后每次增加 6% 的概率

	FourStarDropTimesLimit = 8    // 抽8次后增加概率
	FourStarDropAddValue   = 5100 // 超过后每次增加 51% 的概率
)

var (
	dropList     []*Drop
	dropGroupMap map[int64]*DropGroup
	random       *rand.Rand
)

func init() {
	var err error
	dropList, err = LoadFile[*Drop]("data/Drop.csv", "json")
	if err != nil {
		panic(err)
	}

	makeDropGroupMap()
	random = rand.New(rand.NewSource(time.Now().Unix()))
}

func makeDropGroupMap() {
	dropGroupMap = make(map[int64]*DropGroup)
	for _, v := range dropList {
		drop, ok := dropGroupMap[v.DropId]
		if !ok {
			drop = &DropGroup{
				DropId: v.DropId,
				Drops:  []*Drop{},
			}
			dropGroupMap[v.DropId] = drop
		}

		drop.WeightCount += v.Weight
		drop.Drops = append(drop.Drops, v)
	}
}

func GetDropList() []*Drop {
	return dropList
}

func GetDropGroup(dropId int64) *DropGroup {
	return dropGroupMap[dropId]
}
