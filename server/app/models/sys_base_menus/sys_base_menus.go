package sys_base_menus

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

// SysBaseMenus 结构体 系统菜单
type SysBaseMenus struct {
	models.BaseModel
	ParentId  *uint64        `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:父菜单ID;size:20;"`       //父菜单ID
	Path      *string        `json:"path" form:"path" gorm:"column:path;comment:路由path;size:255;"`                  //路由路径
	Name      *string        `json:"name" form:"name" gorm:"column:name;comment:路由name;size:255;"`                  //路由名称
	Hidden    *int8          `json:"hidden" form:"hidden" gorm:"default:0;column:hidden;comment:是否在列表隐藏;size:1;"`   //是否隐藏
	Component *string        `json:"component" form:"component" gorm:"column:component;comment:对应前端文件路径;size:255;"` //组件路径
	Sort      *uint64        `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:排序标记;size:11;"`           //排序
	Title     *string        `json:"title" form:"title" gorm:"column:title;comment:菜单名;size:255;"`                  //菜单标题
	Icon      *string        `json:"icon" form:"icon" gorm:"column:icon;comment:菜单图标;size:255;"`                    //菜单图标
	Parent    *SysBaseMenus  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"`                     // 父节点
	Children  []SysBaseMenus `json:"children" gorm:"-"`                                                             // 子节点
	models.CommonTimestampsField
}

// TableName 系统菜单 SysBaseMenus自定义表名 sys_base_menus
func (SysBaseMenus) TableName() string {
	return "sys_base_menus"
}

// GetRootSysBaseMenus 获取所有根系统菜单
func GetRootSysBaseMenus() []SysBaseMenus {
	var sysBaseMenuss []SysBaseMenus
	database.DB.Where("parent_id = ? OR parent_id IS NULL", 0).Find(&sysBaseMenuss)
	return sysBaseMenuss
}

// GetSysBaseMenusTree 获取完整的系统菜单树
func GetSysBaseMenusTree() []SysBaseMenus {
	roots := GetRootSysBaseMenus()

	for i := range roots {
		roots[i].loadChildrenRecursive()
	}

	return roots
}

// loadChildrenRecursive 递归加载子系统菜单
func (sysBaseMenus *SysBaseMenus) loadChildrenRecursive() {
	sysBaseMenus.Children = sysBaseMenus.GetChildren()
	for i := range sysBaseMenus.Children {
		sysBaseMenus.Children[i].loadChildrenRecursive()
	}
}

// GetChildren 获取直接子系统菜单
func (sysBaseMenus *SysBaseMenus) GetChildren() []SysBaseMenus {
	var children []SysBaseMenus
	database.DB.Where("parent_id = ?", sysBaseMenus.ID).Find(&children)
	return children
}

// Create 创建系统菜单
func (sysBaseMenus *SysBaseMenus) Create() {
	database.DB.Create(&sysBaseMenus)
}

// Save 保存系统菜单
func (sysBaseMenus *SysBaseMenus) Save() (rowsAffected int64) {
	result := database.DB.Save(&sysBaseMenus)
	return result.RowsAffected
}

// Delete 删除系统菜单
func (sysBaseMenus *SysBaseMenus) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&sysBaseMenus)
	return result.RowsAffected
}

// RouteItem 前端路由项结构
type RouteItem struct {
	Path      string      `json:"path"`
	Name      string      `json:"name"`
	Component string      `json:"component,omitempty"`
	Redirect  string      `json:"redirect,omitempty"`
	Meta      RouteMeta   `json:"meta"`
	Children  []RouteItem `json:"children,omitempty"`
}

// RouteMeta 路由元信息
type RouteMeta struct {
	Title          string `json:"title"`
	Icon           string `json:"icon,omitempty"`
	Hidden         bool   `json:"hidden"`
	KeepAlive      bool   `json:"keepAlive"`
	DefaultMenu    bool   `json:"defaultMenu"`
	CloseTab       bool   `json:"closeTab"`
	ActiveName     string `json:"activeName,omitempty"`
	ActiveMenu     string `json:"activeMenu,omitempty"`
	TransitionType string `json:"transitionType,omitempty"`
	Sort           uint64 `json:"sort"`
}

// MenuListResult 菜单列表返回结构
type MenuListResult struct {
	List []RouteItem `json:"list"`
}

// GetAllMenus 获取所有菜单
func GetAllMenus() ([]SysBaseMenus, error) {
	var menus []SysBaseMenus
	err := database.DB.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// ToRouteItem 将菜单模型转换为前端路由项
func (menu *SysBaseMenus) ToRouteItem() RouteItem {
	item := RouteItem{
		Path: safeString(menu.Path),
		Name: safeString(menu.Name),
		Meta: RouteMeta{
			Title:  safeString(menu.Title),
			Icon:   safeString(menu.Icon),
			Hidden: safeBool(menu.Hidden),
			Sort:   safeInt64(menu.Sort),
		},
	}

	if menu.Component != nil && *menu.Component != "" {
		item.Component = *menu.Component
	}
	return item
}

// BuildMenuTree 构建菜单树
func BuildMenuTree() ([]RouteItem, error) {
	menus, err := GetAllMenus()
	if err != nil {
		return nil, err
	}

	if len(menus) == 0 {
		return []RouteItem{}, nil
	}

	// 使用指针类型构建树结构，方便后续处理
	type TreeNode struct {
		RouteItem
		Children []*TreeNode
	}

	// 转换为TreeNode并建立映射
	nodeMap := make(map[uint64]*TreeNode)

	for _, menu := range menus {
		node := &TreeNode{
			RouteItem: menu.ToRouteItem(),
		}
		nodeMap[menu.ID] = node
	}

	// 构建父子关系 (保持排序顺序，因为 menus 是有序的)
	var rootNodes []*TreeNode
	for _, menu := range menus {
		node := nodeMap[menu.ID]

		if menu.ParentId == nil || *menu.ParentId == 0 {
			// 根节点
			rootNodes = append(rootNodes, node)
		} else {
			// 子节点
			if parent, exists := nodeMap[*menu.ParentId]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	// 定义转换函数
	var convertNode func(*TreeNode) RouteItem
	convertNode = func(node *TreeNode) RouteItem {
		item := node.RouteItem // 复制基础属性

		if len(node.Children) > 0 {
			item.Children = make([]RouteItem, len(node.Children))
			for i, childNode := range node.Children {
				// 递归处理子节点
				childItem := convertNode(childNode)

				// ============ 核心修改逻辑开始 ============
				// 如果这是第一个子节点，建立父级重定向和高亮关联
				if i == 0 {
					// 1. 父级 Redirect 指向第一个子节点的 Path
					if item.Redirect == "" {
						item.Redirect = childItem.Path
					}

					// 2. 隐藏第一个子节点 (这样侧边栏就不会出现下拉框，而是直接点击父级)
					childItem.Meta.Hidden = true

					// 3. 让子页面激活时，高亮父级菜单 (ActiveMenu 指向父级 Path)
					if childItem.Meta.ActiveMenu == "" {
						childItem.Meta.ActiveMenu = item.Path
					}
				}
				// ============ 核心修改逻辑结束 ============

				item.Children[i] = childItem
			}
		}
		return item
	}

	// 转换所有根节点
	roots := make([]RouteItem, len(rootNodes))
	for i, rootNode := range rootNodes {
		roots[i] = convertNode(rootNode)
	}

	return roots, nil
}

// 辅助函数：安全获取字符串值
func safeString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// 辅助函数：安全获取布尔值
func safeBool(ptr *int8) bool {
	if ptr == nil {
		return false
	}
	return *ptr == 1
}

// 辅助函数：安全获取int64值
func safeInt64(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
