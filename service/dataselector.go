package service

import (
	"sort"
	"strings"
	"time"
)

// 可以支持排序 过滤 分页的的数据类型
type dataSelector struct {
	// 数据列表
	GenericDataList []DataCell

	// 数据的参数
	dataSelectQuery *DataSelectQuery
}

// 型转换接口。实现接口以后，都可以用dataSelector的方法
type DataCell interface {
	// 创建时间
	GetCreation() time.Time

	// 资源的名称
	GetName() string
}

// 定义过滤和排序的字段
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

// 过滤的关键字
type FilterQuery struct {
	Name string
}

// 分页的单页数据大小和页数
type PaginateQuery struct {
	Limit int
	Page  int
}

// 排序
// 使用sort.Sort()方法
// 需要实现sort的Len Swap Less方法
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]

}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()

	// b的时间是否在a之前
	return b.Before(a)
}

// 降序排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 过滤
// 根据Name字段判断
func (d *dataSelector) Filter() *dataSelector {
	// name为空，直接返回d，不做处理
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}

	// name不为空，遍历数据并append到list中
	// 创建空list
	dcList := []DataCell{}
	for _, dc := range d.GenericDataList {

		objName := dc.GetName()
		// 字符串判断是否包含
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			continue
		}

		dcList = append(dcList, dc)

	}

	// 返回过滤后的结果
	d.GenericDataList = dcList
	return d
}

// 分页
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page

	// 参数合规
	if limit <= 0 || page <= 0 {
		return d
	}

	// 单页上元素索引
	// 假设元素个数=25 limit=10 page=3
	// 索引从0开始计算，第3页中第一个元素索引20，最后一个元素索引为29，实际上最后一个元素索引是24
	startIndex := limit * (page - 1)
	endIndex := (limit * page) - 1

	// endIndex的处理
	if endIndex > len(d.GenericDataList) {
		endIndex = len(d.GenericDataList) - 1
	}

	// 返回分页的数据
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d

}
