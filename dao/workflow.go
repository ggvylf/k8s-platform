package dao

import (
	"errors"
	"k8s-platform/db"
	"k8s-platform/model"

	"github.com/wonderivan/logger"
)

// dao层自己的Workflow
type workflow struct{}

var Workflow workflow

type WorkflowResp struct {
	// model层的Workflow
	Items []*model.Workflow `json:"items`
	Total int               `json:"total"`
}

// 分页查询
func (w *workflow) GetList(name, namespace string, page, limit int) (data *WorkflowResp, err error) {

	start := (page - 1) * limit

	var workflowList []*model.Workflow

	// 查询数据
	tx := db.DB.Where("namespace is ?", namespace).
		Where("name like ? ", namespace, "%"+name+"%").
		Limit(limit).
		Offset(start).
		Order("id desc").
		Find(&workflowList)

	// 空数据处理
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow数据失败," + tx.Error.Error())
		return nil, errors.New("获取workflow数据失败," + tx.Error.Error())
	}

	return &WorkflowResp{
		Items: workflowList,
		Total: len(workflowList),
	}, nil
}

// 获取单条数据
func (w *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}

	tx := db.DB.Where("id=?", id).
		First(&workflow)

	// 空数据处理
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow数据失败," + tx.Error.Error())
		return nil, errors.New("获取workflow数据失败," + tx.Error.Error())
	}

	return workflow, nil

}

// 新增workflow
// dao层不负责组装workflow数据，组装交给service层做
func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.DB.Create(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("新增workflow失败，" + tx.Error.Error())
		return errors.New("新增workflow失败," + tx.Error.Error())
	}
	return nil
}

// 删除workflow
func (w *workflow) DeleteByID(id int) (err error) {
	tx := db.DB.Where("id=?", id).
		Delete(&model.Workflow{})

	if tx.Error != nil {
		logger.Error("删除workflow失败," + tx.Error.Error())
		return errors.New("删除workflow失败," + tx.Error.Error())
	}
	return nil
}
