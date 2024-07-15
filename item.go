package main

type (
	Item struct {
		ItemId   int64  `json:"ItemId"`   // 物品ID
		SortType int    `json:"SortType"` // 物品类别
		ItemName string `json:"ItemName"` // 物品名称
	}

	Role struct {
		RoleId          int64 `json:"RoleId"` // 物品ID
		Star            int   `json:"Star"`   // 星级
		Stuff           int   `json:"Stuff"`
		StuffNum        int64 `json:"StuffNum"`
		StuffItem       int   `json:"StuffItem"`
		StuffItemNum    int64 `json:"StuffItemNum"`
		MaxStuffItem    int   `json:"MaxStuffItem"`
		MaxStuffItemNum int64 `json:"MaxStuffItemNum"`
		Type            int   `json:"Type"`
	}

	Weapon struct {
		WeaponId int64 `json:"WeaponId"` // 物品ID
		Type     int   `json:"Type"`     // 武器类别，如：双手剑
		Star     int   `json:"Star"`     // 星级
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
	initRoles()
	initWeapons()
}

func initItems(list []*Item) {
	itemMap = make(map[int64]*Item)
	for _, v := range list {
		itemMap[v.ItemId] = v
	}
}

func initRoles() {
	// role data
	roleList, err := LoadFile[*Role]("data/Role.csv", "json")
	if err != nil {
		panic(err)
	}

	roleMap = make(map[int64]*Role)
	for _, v := range roleList {
		roleMap[v.RoleId] = v
	}
}

func initWeapons() {
	// weapon data
	weaponList, err := LoadFile[*Weapon]("data/Weapon.csv", "json")
	if err != nil {
		panic(err)
	}

	weaponMap = make(map[int64]*Weapon)
	for _, v := range weaponList {
		weaponMap[v.WeaponId] = v
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
