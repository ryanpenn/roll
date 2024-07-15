package main

type (
	Item struct {
		ItemId   int64  `json:"ItemId"`   // 物品ID
		SortType int    `json:"SortType"` // 物品类别
		ItemName string `json:"ItemName"` // 物品名称
	}

	Role struct {
		RoleId int64 // 物品ID
		Star   int   // 星级
	}

	Weapon struct {
		WeaponId int64 // 物品ID
		Type     int   // 武器类别，如：双手剑
		Star     int   // 星级
	}
)

var (
	itemMap   map[int64]*Item
	roleMap   map[int64]*Role
	weaponMap map[int64]*Weapon
)

func init() {
	itemList, err := LoadFile[*Item]("data/Item.csv", "json")
	if err != nil {
		panic(err)
	}

	initItems(itemList)
	initRoles(itemList)
	initWeapons(itemList)
}

func initItems(list []*Item) {
	itemMap = make(map[int64]*Item)
	for _, v := range list {
		itemMap[v.ItemId] = v
	}
}

func initRoles(list []*Item) {
	// mock role data
	roleMap = make(map[int64]*Role)
	start := 3
	for _, v := range list {
		if v.ItemId > 2000000 && v.ItemId < 3000000 {
			if v.ItemId > 2000000 && v.ItemId <= 2000006 {
				start = 5
			} else if v.ItemId > 2000006 && v.ItemId <= 2000030 {
				start = 4
			} else {
				start = 3
			}

			roleMap[v.ItemId] = &Role{
				RoleId: v.ItemId,
				Star:   start,
			}
		}
	}
}

func initWeapons(list []*Item) {
	// mock weapon data
	weaponMap = make(map[int64]*Weapon)
	start := 3
	for _, v := range list {
		if v.ItemId > 6000000 && v.ItemId < 7000000 {
			if v.ItemId == 6000002 {
				start = 4
			} else if v.ItemId == 6000003 {
				start = 5
			} else {
				start = 3
			}

			weaponMap[v.ItemId] = &Weapon{
				WeaponId: v.ItemId,
				Type:     1,
				Star:     start,
			}
		}
	}
}

func GetItem(itemId int64) *Item {
	return itemMap[itemId]
}

func GetRole(roleId int64) *Role {
	return roleMap[roleId]
}

func GetWeapon(weaponId int64) *Weapon {
	return weaponMap[weaponId]
}
