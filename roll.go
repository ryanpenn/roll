package main

import (
	"fmt"
)

type RollInfo struct {
	PoolId        int64 // 奖池
	FiveStarTimes int   // 5星次数
	FourStarTimes int   // 4星次数
	FiveMustHit   bool  // 5星大保底
}

var roleRollInfo *RollInfo

func init() {
	roleRollInfo = &RollInfo{
		PoolId:      1000,
		FiveMustHit: false,
	}
}

func rollInGroup(dropGroup *DropGroup) *Drop {
	if dropGroup != nil {
		num := random.Intn(dropGroup.WeightCount)
		weight := 0
		for _, v := range dropGroup.Drops {
			weight += v.Weight
			if weight > num {
				if v.IsEnd == 1 {
					return v
				}

				return rollInGroup(GetDropGroup(v.Result))
			}
		}
	}

	return nil
}

func Roll(times int) {
	// 统计结果
	result := make(map[int64]int)

	for i := 0; i < times; i++ {
		// 增加次数
		roleRollInfo.FiveStarTimes++
		roleRollInfo.FourStarTimes++

		group := GetDropGroup(roleRollInfo.PoolId)
		if group == nil {
			fmt.Printf("配置错误: 奖池 %d 不存在", roleRollInfo.PoolId)
			return
		}

		// 4星、5星保底算法
		if roleRollInfo.FiveStarTimes > FiveStarDropTimesLimit ||
			roleRollInfo.FourStarTimes > FourStarDropTimesLimit {
			newGroup := &DropGroup{
				DropId:      group.DropId,
				WeightCount: group.WeightCount,
			}

			fiveStarAddVal := (roleRollInfo.FiveStarTimes - FiveStarDropTimesLimit) * FiveStarDropAddValue
			if fiveStarAddVal < 0 {
				fiveStarAddVal = 0
			}
			fourStarAddVal := (roleRollInfo.FourStarTimes - FourStarDropTimesLimit) * FourStarDropAddValue
			if fourStarAddVal < 0 {
				fourStarAddVal = 0
			}

			for _, v := range group.Drops {
				newDrop := &Drop{
					DropId: v.DropId,
					Result: v.Result,
					IsEnd:  v.IsEnd,
				}

				// 调权重
				switch v.Result {
				case 10001: // 5
					newDrop.Weight = v.Weight + fiveStarAddVal
				case 10002: // 4
					newDrop.Weight = v.Weight + fourStarAddVal
				case 10003: // 3
					newDrop.Weight = v.Weight - fiveStarAddVal - fourStarAddVal
				}

				newGroup.Drops = append(newGroup.Drops, newDrop)
			}

			group = newGroup
		}

		ret := rollInGroup(group)
		if ret != nil {
			// 检查是否已获得5星英雄
			role := GetRole(ret.Result)
			if role != nil {
				switch role.Star {
				case 5:
					roleRollInfo.FiveStarTimes = 0 // 恢复5星概率

					if roleRollInfo.FiveMustHit {
						mustGroup := GetDropGroup(100012) // 5星保底
						if mustGroup != nil {
							ret = rollInGroup(mustGroup)
							if ret == nil {
								fmt.Println("5星保底数据配置错误")
								return
							}
						}
					}

					if ret.DropId == 100012 { // 是否抽到当期的保底英雄
						roleRollInfo.FiveMustHit = false
					} else {
						roleRollInfo.FiveMustHit = true
					}

				case 4:
					roleRollInfo.FourStarTimes = 0 // 恢复4星概率
					// todo 4星保底?
				default:
					// do nothing
				}
			}

			// 检查获得的武器
			if weapon := GetWeapon(ret.Result); weapon != nil {
				switch weapon.Star {
				case 5:
					roleRollInfo.FiveStarTimes = 0 // 恢复5星概率
				case 4:
					roleRollInfo.FourStarTimes = 0 // 恢复4星概率
				}
			}

			result[ret.Result]++
		}
	}

	//-----------------------------------
	// 概率统计

	fmt.Println()
	fmt.Printf("%d 次英雄抽取结果统计:\n", times)

	var sum int
	wpMap := make(map[int64]int) // 武器概率
	starMap := make(map[int]int) // 英雄星级概率
	for k, v := range result {
		if role := GetRole(k); role != nil {
			// 抽到英雄
			if _, ok := starMap[role.Star]; !ok {
				starMap[role.Star] = 0
			}

			sum += v
			starMap[role.Star] += v
			fmt.Printf("英雄:%s  \t%d星  \t数量:%d  \t概率:%.2f%%\n", GetItem(k).ItemName, role.Star, v, float32(v)*100.0/float32(times))
		} else {
			// 抽到武器
			if wp := GetWeapon(k); wp != nil {
				starMap[wp.Star] += v
			}
			if _, ok := wpMap[k]; !ok {
				wpMap[k] = 0
			}
			wpMap[k] += v
			sum += v
		}
	}

	fmt.Println()
	fmt.Printf("%d 次英雄星级概率:\n", sum)
	expect := map[int][2]string{
		5: {"0.6%", "1.6052%"},
		4: {"5.1%", "13.057%"},
	}
	for k, v := range starMap {
		fmt.Printf("%d星物品  \t%d次  \t概率: %.4f%%   \t%s \t%s\n", k, v, float32(v)*100.0/float32(times), expect[k][0], expect[k][1])
	}

	fmt.Println()
	for k, v := range wpMap {
		fmt.Printf("武器:%s  \t%d星  \t数量:%d  \t概率:%.2f%%\n", GetItem(k).ItemName, GetWeapon(k).Star, v, float32(v)*100.0/float32(times))
	}

	fmt.Println()
}
